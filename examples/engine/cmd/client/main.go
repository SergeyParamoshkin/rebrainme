package main

import (
	"context"
	"flag"
	"log"
	"path"
	"time"

	"github.com/SergeyParamoshkin/rebrainme/examples/engine/pb/v1/calculation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:8080", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption

	if *tls {
		if *caFile == "" {
			*caFile = path.Base("x509/ca_cert.pem")
		}

		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := calculation.NewCalculationServiceClient(conn)

	printFeature(client, &calculation.Calculation{
		Id:      0,
		Comment: "hello world",
	})

}

func printFeature(client calculation.CalculationServiceClient, c *calculation.Calculation) {
	// log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feature, err := client.Get(ctx, c)
	if err != nil {

	}

	log.Println(feature)
}
