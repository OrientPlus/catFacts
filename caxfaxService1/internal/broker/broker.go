package broker

import (
	"caxfaxService1/internal/entity"
	"caxfaxService1/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type Broker struct {
	consumer  sarama.PartitionConsumer
	producer  sarama.SyncProducer
	logger    logger.Logger
	RequestCh chan entity.Request
	closeCh   chan bool
}

func NewBroker(config entity.Config) (*Broker, error) {
	broker := &Broker{}
	broker.logger = config.Logger
	broker.RequestCh = make(chan entity.Request, 10)
	broker.closeCh = make(chan bool)

	addrs := fmt.Sprintf("%s:%s", config.BrokerConsumerHost, config.BrokerConsumerPort)
	consumer, err := sarama.NewConsumer([]string{addrs}, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания cosumer'а: %s", err.Error())
	}

	broker.consumer, err = consumer.ConsumePartition("service1.requests", 0, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания partition cosumer'а: %s", err.Error())
	}

	broker.producer, err = sarama.NewSyncProducer([]string{addrs}, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания producer'а: %s", err)
	}

	return broker, nil
}

func (b *Broker) SendResponse(request entity.Request, fact entity.Fact) error {
	responseData, _ := json.Marshal(fact)

	message := &sarama.ProducerMessage{
		Topic: request.ReplyTopic,
		Value: sarama.ByteEncoder(responseData),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("Correlation-ID"),
				Value: []byte(request.CorrelationID),
			},
		},
	}

	_, _, err := b.producer.SendMessage(message)
	if err != nil {
		b.logger.Errorf("ошибка отправки сообщения %s", err.Error())
	}
	return err
}

func (b *Broker) Run() {
	go b.handleRequests()
	return
}

func (b *Broker) handleRequests() {
	for {
		select {
		case msg := <-b.consumer.Messages():
			b.logger.Debugf("Получено сообщение; Topic: %s, Partition: %d", msg.Topic, msg.Partition)
			req := entity.Request{}
			if err := json.Unmarshal(msg.Value, &req); err != nil {
				b.logger.Errorf("ошибка анмаршаллинга json объекта: %s", err.Error())
				continue
			}

			b.RequestCh <- req
			break

		case <-b.closeCh:
			return
		}
	}
}

func (b *Broker) Stop() {
	b.closeCh <- true
	b.consumer.Close()
	b.producer.Close()
}
