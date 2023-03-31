package services

import (
	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type NATSService struct {
	log               *zerolog.Logger
	Jetstream         nats.JetStreamContext
	JetStreamName     string
	JetStreamSubject  string
	AckTimeoutMinutes int
	DurableConsumer   string
}

func NewNATSService(settings *config.Settings, log *zerolog.Logger) (*NATSService, error) {
	n, err := nats.Connect(settings.NATSURL)
	if err != nil {
		return nil, err
	}

	js, err := n.JetStream()
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      settings.NATSStreamName,
		Retention: nats.WorkQueuePolicy,
		Subjects:  []string{settings.NATSDataDownloadSubject},
	})
	if err != nil {
		return nil, err
	}

	return &NATSService{
		log:               log,
		Jetstream:         js,
		JetStreamName:     settings.NATSStreamName,
		JetStreamSubject:  settings.NATSDataDownloadSubject,
		AckTimeoutMinutes: settings.NATSAckTimeoutMinutes,
		DurableConsumer:   settings.NATSDurableConsumer}, nil
}
