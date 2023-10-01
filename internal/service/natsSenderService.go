package service

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"nats-streaming-consumer/internal/entity"
)

type NatsSenderService struct {
	connection *nats.Conn
	chanelName string
}

func NewNatsSenderService(chanelName string) (*NatsSenderService, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	return &NatsSenderService{
		connection: nc,
		chanelName: chanelName,
	}, nil
}

func (s *NatsSenderService) PublishModel(model entity.Model) error {
	marshalledModel, err := json.Marshal(model)
	if err != nil {
		return err
	}
	err = s.connection.Publish(s.chanelName, marshalledModel)
	if err != nil {
		return err
	}
	return nil
}

func (s *NatsSenderService) Close() {
	s.connection.Close()
}
