package server_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/SergeyParamoshkin/rebrainme/examples/engine/internal/grpc/server"
	"github.com/SergeyParamoshkin/rebrainme/examples/engine/pb/v1/calculation"
)

func TestCalculationServiceServer_Get(t *testing.T) {
	type fields struct {
		UnimplementedCalculationServiceServer calculation.UnimplementedCalculationServiceServer
		Response                              *calculation.CalculationResponse
	}
	type args struct {
		in0 context.Context
		in1 *calculation.Calculation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *calculation.CalculationResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server.CalculationServiceServer{
				UnimplementedCalculationServiceServer: tt.fields.UnimplementedCalculationServiceServer,
				Response:                              tt.fields.Response,
			}
			got, err := s.Get(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculationServiceServer.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculationServiceServer.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
