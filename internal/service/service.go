package service

import (
	"github.com/redis/go-redis/v9"

	pb "feature_flag_config/pkg/pb/feature_flag_config"
)

type Service struct {
	redisDB *redis.Client
	pb.UnimplementedFeatureFlagConfigServiceServer
}

func NewService(db *redis.Client) *Service {
	return &Service{
		redisDB: db,
	}
}
