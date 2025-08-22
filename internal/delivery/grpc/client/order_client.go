package client

import (
	"context"
	"fmt"
	"time"

	pb "golectro-payment/internal/delivery/grpc/proto/order"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type OrderClient struct {
	client pb.OrderServiceClient
	log    *logrus.Logger
	viper  *viper.Viper
}

func NewOrderClient(log *logrus.Logger, viper *viper.Viper) *OrderClient {
	conn, err := grpc.Dial(viper.GetString("GRPC_ORDER_SERVICE"), grpc.WithInsecure())
	if err != nil {
		return nil
	}

	client := pb.NewOrderServiceClient(conn)

	return &OrderClient{
		client: client,
		log:    log,
		viper:  viper,
	}
}

func (p *OrderClient) GetOrderByID(ctx context.Context, orderID string) (*pb.GetOrderByIdResponse, error) {
	req := &pb.GetOrderByIdRequest{
		Id: orderID,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := p.client.GetOrderByID(ctx, req)
	if err != nil {
		p.log.WithError(err).Error("Failed to get order by ID")
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}
	if resp == nil {
		p.log.Warn("Order not found for ID: ", orderID)
		return nil, fmt.Errorf("order not found for ID: %s", orderID)
	}
	return resp, nil
}
