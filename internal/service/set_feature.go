package service

import (
	"context"

	pb "feature_flag_config/pkg/pb/feature_flag_config"
)

func (s *Service) SetFeature(ctx context.Context, request *pb.SetFeatureRequest) (*pb.SetFeatureResponse, error) {

	return &pb.SetFeatureResponse{}, nil
}
