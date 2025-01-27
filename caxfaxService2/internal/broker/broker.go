package broker

import (
	"caxfaxService2/internal/entity"
	"caxfaxService2/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"math/rand"
	"time"
)

type Broker struct {
	consumer      sarama.PartitionConsumer
	producer      sarama.SyncProducer
	responseTopic string
	requestTopic  string
	logger        logger.Logger
	closeCh       chan bool
}

type Fact struct {
	Message string `json:"fact"`
	Length  int32  `json:"length"`
}

type RequestMessage struct {
	CorrelationID string `json:"correlation_id"`
	ReplyTopic    string `json:"reply_topic"`
}

type ResponseMessage struct {
	CorrelationID string `json:"correlation_id"`
	Fact          Fact   `json:"fact"`
}

func (b *Broker) RequestFact(ctx context.Context) (Fact, error) {
	correlationID := generateCorrelationID()

	// Формируем сообщение-запрос
	reqMsg := RequestMessage{
		CorrelationID: correlationID,
		ReplyTopic:    b.responseTopic,
	}

	reqBytes, err := json.Marshal(reqMsg)
	if err != nil {
		b.logger.Errorf("Не удалось смаршалить запрос: %s", err.Error())
		return Fact{}, err
	}

	// Отправляем запрос
	_, _, err = b.producer.SendMessage(&sarama.ProducerMessage{
		Topic: b.requestTopic,
		Value: sarama.ByteEncoder(reqBytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("Correlation-ID"),
				Value: []byte(correlationID),
			},
		},
	})
	if err != nil {
		b.logger.Errorf("Не удалось отправить запрос: %s", err.Error())
		return Fact{}, err
	}
	b.logger.Debugf("Сообщение отправлено в топик %s", b.requestTopic)

	select {
	case msg := <-b.consumer.Messages():
		var fact Fact
		if err = json.Unmarshal(msg.Value, &fact); err != nil {
			b.logger.Errorf("Не удалось анмаршалить ответ: %s", err.Error())
			return Fact{}, err
		}

		return fact, nil

	case <-ctx.Done():
		return Fact{}, errors.New("request timed out")
	}
}

func generateCorrelationID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func (b *Broker) Stop() {
	b.closeCh <- true
	b.consumer.Close()
	b.producer.Close()
}

func NewBroker(config entity.Config) (*Broker, error) {
	broker := &Broker{}
	broker.logger = config.Logger
	broker.closeCh = make(chan bool, 1)
	broker.responseTopic = config.ResponseTopic
	broker.requestTopic = config.RequestTopic

	addrs := fmt.Sprintf("%s:%s", config.BrokerConsumerHost, config.BrokerConsumerPort)
	var err error
	broker.producer, err = sarama.NewSyncProducer([]string{addrs}, nil)
	if err != nil {
		broker.logger.Errorf("Не удалось создать producer'а: %s", err.Error())
		return nil, err
	}

	consumer, err := sarama.NewConsumer([]string{addrs}, nil)
	if err != nil {
		broker.logger.Errorf("Не удалось создать consumer'а: %s", err.Error())
		return nil, err
	}
	broker.consumer, err = consumer.ConsumePartition(
		broker.responseTopic,
		0,
		sarama.OffsetNewest,
	)
	if err != nil {
		broker.logger.Errorf("Не удалось создать partition consumer'а: %s", err.Error())
		return nil, err
	}

	return broker, nil
}
