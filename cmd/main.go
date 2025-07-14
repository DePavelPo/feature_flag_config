package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/lib/pq"

	"github.com/redis/go-redis/v9"

	"feature_flag_config/internal/service"
	pb "feature_flag_config/pkg/pb/feature_flag_config"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	redisDB := initRedis()

	internalService := service.NewService(redisDB)

	grpcServer := grpc.NewServer()
	pb.RegisterFeatureFlagConfigServiceServer(grpcServer, internalService)

	listener, err := net.Listen("tcp", os.Getenv("API_ADDR_GRPC"))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		logrus.Infof("starting gRPC server at %s", os.Getenv("API_ADDR_GRPC"))
		if err := grpcServer.Serve(listener); err != nil {
			logrus.Fatalf("gRpc server error: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	pb.RegisterFeatureFlagConfigServiceHandlerFromEndpoint(ctx, mux, os.Getenv("API_ADDR_GRPC"), opts)

	wg.Add(1)
	go func() {
		defer wg.Done()

		logrus.Infof("starting HTTP server at %s", os.Getenv("API_ADDR_HTTP"))
		srv := &http.Server{
			Handler:      mux,
			Addr:         os.Getenv("API_ADDR_HTTP"),
			WriteTimeout: 1 * time.Second,
			ReadTimeout:  1 * time.Second,
		}
		logrus.Fatalf("http server error: %v", srv.ListenAndServe())
	}()

	wg.Wait()
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       1,
	})
}
