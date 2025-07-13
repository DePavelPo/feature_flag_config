package service

import (
	"context"

	pb "feature_flag_config/pkg/pb/feature_flag_config"
)

func (s *Service) GetFeaturesByOptions(
	ctx context.Context,
	request *pb.GetFeaturesByOptionsRequest,
) (*pb.GetFeaturesByOptionsResponse, error) {

	return &pb.GetFeaturesByOptionsResponse{}, nil
}
