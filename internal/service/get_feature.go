package service

import (
	"context"

	"feature_flag_config/db/cache"
	pb "feature_flag_config/pkg/pb/feature_flag_config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func (s *Service) GetFeaturesByOptions(
	ctx context.Context,
	request *pb.GetFeaturesByOptionsRequest,
) (*pb.GetFeaturesByOptionsResponse, error) {
	// if we didn't get features, we need to return all feature
	// so we need to get all keys first
	if len(request.FeatureNames) == 0 {
		var err error
		if request.FeatureNames, err = getAllKeys(ctx, s.redisDB); err != nil {
			logrus.Errorf("getAllKeys error: %v", err)
			return &pb.GetFeaturesByOptionsResponse{
				Error: &pb.Error{
					Code:    500,
					Message: err.Error(),
				},
			}, nil
		}
	}

	// go to get features by feature names
	features, err := getByNames(ctx, request.FeatureNames, s.redisDB)
	if err != nil {
		logrus.Errorf("getByNames error: %v", err)
		return &pb.GetFeaturesByOptionsResponse{
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	// filter output data by IsActive flag if it's provided
	if request.IsActive != nil {
		needToSaveIdx := make(map[int]struct{})

		for i := 0; i < len(features); i++ {
			if features[i].IsActive == *request.IsActive {
				needToSaveIdx[i] = struct{}{}
			}
		}

		features = filterByIsActive(features, needToSaveIdx)
	}

	return getFeaturesByOptionsFeatureToResponse(features), nil
}

func getAllKeys(ctx context.Context, db *redis.Client) ([]string, error) {
	keys := []string{}
	var cursor uint64
	for {
		var keysBatch []string
		var err error
		keysBatch, cursor, err = db.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, keysBatch...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func getByNames(ctx context.Context, featureNames []string, db *redis.Client) (features []*cache.Feature, err error) {
	for i := 0; i < len(featureNames); i++ {
		feature := &cache.Feature{}
		err = feature.GetByName(ctx, featureNames[i], db)
		if err != nil {
			return
		}

		features = append(features, feature)
	}

	return
}

func filterByIsActive[T any](features []T, needToSaveIdx map[int]struct{}) []T {
	write := 0
	for read := 0; read < len(features); read++ {
		if _, found := needToSaveIdx[read]; found {
			features[write] = features[read]
			write++
		}
	}
	return features[:write]
}

func getFeaturesByOptionsFeatureToResponse(features []*cache.Feature) *pb.GetFeaturesByOptionsResponse {
	data := make([]*pb.GetFeaturesByOptionsResponse_GetFeaturesByOptionsResponseData, 0, len(features))
	for _, feature := range features {
		data = append(data, &pb.GetFeaturesByOptionsResponse_GetFeaturesByOptionsResponseData{
			Name:          feature.Name,
			IsActive:      feature.IsActive,
			BucketsOpened: int32(feature.BucketsOpened),
			WhiteList:     feature.Whitelist,
			BlackList:     feature.Blacklist,
		})
	}

	return &pb.GetFeaturesByOptionsResponse{
		Data: data,
	}
}
