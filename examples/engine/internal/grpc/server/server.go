package server

import (
	"context"

	"github.com/SergeyParamoshkin/rebrainme/examples/engine/pb/v1/calculation"
	"google.golang.org/grpc"
)

type CalculationServiceServer struct {
	calculation.UnimplementedCalculationServiceServer

	Response *calculation.CalculationResponse
}

func (s *CalculationServiceServer) Get(
	context.Context, *calculation.Calculation,
) (*calculation.CalculationResponse, error) {
	return &calculation.CalculationResponse{
		Result: 0.9999,
	}, nil
}

func NewServer() *grpc.Server {
	gsrv := grpc.NewServer()
	srv := &CalculationServiceServer{
		UnimplementedCalculationServiceServer: calculation.UnimplementedCalculationServiceServer{},
	}
	calculation.RegisterCalculationServiceServer(gsrv, srv)

	return gsrv
}
