package outbound

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	model "github.com/evri/CashlessPayments/SBController/model"
	pb "github.com/evri/CashlessPayments/SBController/proto"
)

// Load SB
func (sb SBClient) Load() (string, error) {

	conn, err := grpc.Dial(configuration.SBHOST, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("error in connection: %v", err)
	}

	defer conn.Close()

	c := pb.NewFundClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var resp string

	r, err := c.Load(ctx, &pb.TransferFunds{
		CashableMoneyInCents:      sb.request.Cents,
		RestrictedMoneyInCents:    0,
		NonRestrictedMoneyInCents: 0,
	})
	if err != nil {
		log.Fatalf("load error: %v", err)
	}

	resp = r.GetResponse()

	if err != nil {
		log.Fatalf("could not disable: %v", err)
		return "", err
	}

	log.Printf("Greeting: %s", resp)
	return resp, nil
}

func (sb SBClient) UnLoad() (string, error) {

	conn, err := grpc.Dial(configuration.SBHOST, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("error in connection: %v", err)
	}

	defer conn.Close()

	c := pb.NewFundClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var resp string
	r, err := c.UnLoad(ctx, &pb.Empty{})

	if err != nil {
		log.Fatalf("unload error: %v", err)
	}

	resp = r.GetResponse()

	if err != nil {
		log.Fatalf("could not disable: %v", err)
		return "", err
	}

	log.Printf("Greeting: %s", resp)
	return resp, nil
}

func (sb SBClient) UpdateJx(jxRequest model.JXUpdateRequest) model.CommonResponse {

	conn, err := grpc.Dial(configuration.SBHOST, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("error in connection: %v", err)
	}

	defer conn.Close()

	c := pb.NewFundClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UpdateJx(ctx, &pb.JXRequest{
		DispenseType: jxRequest.DispenseType,
		Amount:       jxRequest.Amount,
	})

	if err != nil {
		log.Fatalf("could not disable: %v", err)
	}

	log.Printf("Greeting: %s", r.GetResponse())

	return model.CommonResponse{Status: r.GetResponse()}
}
