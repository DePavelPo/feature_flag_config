package service

import (
	"context"

	pb "feature_flag_config/pkg/pb/feature_flag_config"
)

func (s *Service) CheckFeatureAvailability(
	ctx context.Context,
	request *pb.FeatureAvailabilityRequest,
) (*pb.FeatureAvailabilityResponse, error) {
	enabled, reason := s.checkFeature(ctx, request.FeatureName, request.ItemId)

	return &pb.FeatureAvailabilityResponse{
		Data: &pb.FeatureAvailabilityResponse_FeatureAvailabilityResponseData{
			Enabled: enabled,
			Reason:  reason,
		},
		Error: nil,
	}, nil
}
