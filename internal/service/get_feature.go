package service

import (
	"context"

	"feature_flag_config/db/cache"
	pb "feature_flag_config/pkg/pb/feature_flag_config"

	"github.com/sirupsen/logrus"
)

func (s *Service) GetFeaturesByOptions(
	ctx context.Context,
	request *pb.GetFeaturesByOptionsRequest,
) (*pb.GetFeaturesByOptionsResponse, error) {
	if len(request.FeatureNames) == 0 {
		return &pb.GetFeaturesByOptionsResponse{}, nil
	}

	feature := cache.Feature{}
	if err := feature.GetByName(ctx, request.FeatureNames[0], s.redisDB); err != nil {
		logrus.Errorf("GetByName error: %v", err)
		return &pb.GetFeaturesByOptionsResponse{
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	return getFeaturesByOptionsFeatureToResponse(feature), nil
}

func getFeaturesByOptionsFeatureToResponse(feature cache.Feature) *pb.GetFeaturesByOptionsResponse {
	data := make([]*pb.GetFeaturesByOptionsResponse_GetFeaturesByOptionsResponseData, 1)
	data[0] = &pb.GetFeaturesByOptionsResponse_GetFeaturesByOptionsResponseData{
		Name:          feature.Name,
		IsActive:      feature.IsActive,
		BucketsOpened: int32(feature.BucketsOpened),
		WhiteList:     feature.Whitelist,
		BlackList:     feature.Blacklist,
	}

	return &pb.GetFeaturesByOptionsResponse{
		Data: data,
	}
}
