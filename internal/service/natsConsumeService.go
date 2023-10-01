package service

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"nats-streaming-consumer/internal/entity"
)

type ModelSaver interface {
	AddNewModel(model entity.Model) (int, error)
}

type NatsConsumeService struct {
	connection *nats.Conn
	chanelName string
	modelSaver ModelSaver
}

func NewNatsConsumeService(chanelName string, saver ModelSaver) (*NatsConsumeService, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	consumer := &NatsConsumeService{
		connection: nc,
		chanelName: chanelName,
		modelSaver: saver,
	}

	_, err = nc.Subscribe(chanelName, consumer.ConsumeMessage)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func (s *NatsConsumeService) Close() {
	s.connection.Close()
}

func (s *NatsConsumeService) ConsumeMessage(m *nats.Msg) {
	var unmarshalledModel entity.Model
	err := json.Unmarshal(m.Data, &unmarshalledModel)
	if err != nil {
		logrus.Errorf("Failed to unmarshall model: %s", err.Error())
		return
	}

	newModelId, err := s.modelSaver.AddNewModel(unmarshalledModel)
	if err != nil {
		logrus.Errorf("Failed to save model: %s", err.Error())
		return
	}
	logrus.Infof("Model successfully saved with id=%d", newModelId)
}
