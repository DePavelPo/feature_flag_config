package service

import (
	"context"

	"feature_flag_config/db/cache"
	pb "feature_flag_config/pkg/pb/feature_flag_config"

	"github.com/sirupsen/logrus"
)

func (s *Service) SetFeature(ctx context.Context, request *pb.SetFeatureRequest) (*pb.SetFeatureResponse, error) {
	feature := setFeatureRequestToFeature(request)

	if err := feature.SetInRedis(ctx, s.redisDB); err != nil {
		logrus.Errorf("SetInRedis error: %v", err)
		return &pb.SetFeatureResponse{}, err
	}

	return &pb.SetFeatureResponse{}, nil
}

func setFeatureRequestToFeature(request *pb.SetFeatureRequest) cache.Feature {
	return cache.Feature{
		Name:          request.Name,
		IsActive:      request.IsActive,
		BucketsOpened: int(request.BucketsOpened),
		Whitelist:     request.GetWhiteList(),
		Blacklist:     request.GetBlackList(),
	}
}
