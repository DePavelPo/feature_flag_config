package service

import (
	"context"

	pb "feature_flag_config/pkg/pb/feature_flag_config"

	"github.com/sirupsen/logrus"
)

func (s *Service) CheckFeatureAvailability(
	ctx context.Context,
	request *pb.FeatureAvailabilityRequest,
) (*pb.FeatureAvailabilityResponse, error) {
	logrus.Infof("request %v", request)

	enabled, reason := s.checkFeature(ctx, request.FeatureName, request.ItemId)

	return &pb.FeatureAvailabilityResponse{
		Data: &pb.FeatureAvailabilityResponse_FeatureAvailabilityResponseData{
			Enabled: enabled,
			Reason:  reason,
		},
		Error: nil,
	}, nil
}
