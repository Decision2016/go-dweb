/**
  @author: decision
  @date: 2024/7/24
  @note:
**/

package main

import (
	"context"
	"encoding/hex"
	"github.com/gogo/protobuf/proto"
	"github.com/gookit/config/v2"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/plugins/norn/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NornChain struct {
	rpc        string
	privateKey string

	client  *pb.BlockchainClient
	address string
}

func (n *NornChain) Identity() (string, error) {
	req := &pb.ReadContractAddressReq{
		Address: proto.String(n.address),
		Key:     proto.String("identity"),
	}

	rsp, err := (*n.client).ReadContractAddress(context.TODO(), req)
	if err != nil {
		return "", err
	}

	identity, err := hex.DecodeString(*rsp.Hex)
	if err != nil {
		return "", err
	}

	return string(identity), nil
}

func (n *NornChain) Bootstrap() (string, error) {
	req := &pb.ReadContractAddressReq{
		Address: proto.String(n.address),
		Key:     proto.String("bootstrap"),
	}

	rsp, err := (*n.client).ReadContractAddress(context.TODO(), req)
	if err != nil {
		return "", err
	}

	identity, err := hex.DecodeString(*rsp.Hex)
	if err != nil {
		return "", err
	}

	return string(identity), nil
}

func (n *NornChain) Initial(ident string, url string) error {
	return n.SetIdentity(ident)
}

func (n *NornChain) SetIdentity(ident string) error {
	req := &pb.SendTransactionWithDataReq{
		Type:     proto.String("set"),
		Receiver: proto.String(n.address),
		Key:      proto.String("identity"),
		Value:    proto.String(ident),
	}

	_, err := (*n.client).SendTransactionWithData(context.TODO(), req)
	if err != nil {
		return err
	}

	return nil
}

func (n *NornChain) Join(url string) error {
	req := &pb.SendTransactionWithDataReq{
		Type:     proto.String("set"),
		Receiver: proto.String(n.address),
		Key:      proto.String("boostrap"),
		Value:    proto.String(url),
	}

	_, err := (*n.client).SendTransactionWithData(context.TODO(), req)
	if err != nil {
		return err
	}

	return nil
}

func (n *NornChain) Setup(address string) error {
	n.rpc = config.String("chain.url")
	n.privateKey = config.String("chain.private")

	conn, err := grpc.Dial(n.rpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithError(err).Errorln("error occurred when dial rpc server")
		return err
	}

	c := pb.NewBlockchainClient(conn)
	n.client = &c

	return nil
}
