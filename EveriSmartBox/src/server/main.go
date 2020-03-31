package main

import (
	"context"
	"log"
	"net"
	"runtime"

	"google.golang.org/grpc"

	"github.com/evri/CashlessPayments/EveriSmartBox/src/config"
	core "github.com/evri/CashlessPayments/EveriSmartBox/src/core"
	"github.com/evri/CashlessPayments/EveriSmartBox/src/lib"
	pb "github.com/evri/CashlessPayments/EveriSmartBox/src/proto"
)

type server struct{}

const (
	port = ":50051"
)

var configuration = config.New()

func (s *server) DisableEGM(ctx context.Context, req *pb.Empty) (*pb.CommonResponse, error) {
	core.DisableEGM()
	return &pb.CommonResponse{
		Response: "Success",
	}, nil
}
func (s *server) EnableEGM(ctx context.Context, req *pb.Empty) (*pb.CommonResponse, error) {
	core.EnableEGM()
	return &pb.CommonResponse{
		Response: "Success",
	}, nil
}
func (s *server) Load(ctx context.Context, req *pb.TransferFunds) (*pb.CommonResponse, error) {
	go core.LoadEGMwithFunds(req.GetCashableMoneyInCents(), req.GetRestrictedMoneyInCents(), req.GetNonRestrictedMoneyInCents())
	return &pb.CommonResponse{
		Response: "Success",
	}, nil
}

func (s *server) UnLoad(ctx context.Context, req *pb.Empty) (*pb.CommonResponse, error) {
	go core.UnLoadFunds()
	return &pb.CommonResponse{
		Response: "Success",
	}, nil
}

func (s *server) UpdateJx(ctx context.Context, req *pb.JXRequest) (*pb.CommonResponse, error) {
	go core.HandsetRepay(req.GetDispenseType())
	return &pb.CommonResponse{
		Response: "Success",
	}, nil
}

func main() {

	lib.Log("Hello %s/%s\n", runtime.GOOS, runtime.GOARCH)

	go core.InitializePort(configuration.ActivePort, configuration.PassivePort)
	go core.DequeueIncomingMessage()

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterFundServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
