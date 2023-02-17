package service

import (
	"errors"
	"fmt"
	"github.com/COTBU/notifier/service/sender"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/COTBU/notifier/config"
	"github.com/Shopify/sarama"
)

type Service struct {
	appConfig *config.Config

	saramaClient sarama.Client
	isRunning    bool

	notifySender sender.Sender
}

func (s *Service) CloseClient() error {
	return s.saramaClient.Close()
}

func New(appConfig *config.Config) *Service {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = "notifier"
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Producer.Return.Successes = true

	client, err := sarama.NewClient([]string{
		fmt.Sprintf(
			"%s:%d",
			appConfig.Broker.Host,
			appConfig.Broker.Port,
		),
	}, saramaConfig)
	if err != nil {
		log.Fatalln(err)
	}

	return &Service{
		saramaClient: client,
	}
}

func (s *Service) RunConsumer() error {
	if s.isRunning {
		return errors.New("consumer already running")
	}

	master, err := sarama.NewConsumerFromClient(s.saramaClient)
	if err != nil {
		return err
	}

	s.isRunning = true

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topics, err := master.Topics()
	if err != nil {
		panic(err)
	}

	fmt.Println("got topics", topics)

	messages, errChan := consume(topics, master)

	// Get signal for finish
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	doneCh := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case msg := <-messages:
				fmt.Println(
					"Received messages",
					"key:", string(msg.Key),
					"value:", string(msg.Value),
				)
				// todo send email
			case consumerError := <-errChan:
				fmt.Println(
					"Received messages error ",
					consumerError.Topic,
					string(consumerError.Partition),
					consumerError.Err,
				)
				doneCh <- struct{}{}
				return
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
				return
			}
		}
	}()

	<-doneCh
	s.isRunning = false

	return nil
}

func consume(topics []string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumerMessages := make(chan *sarama.ConsumerMessage)
	consumerErrors := make(chan *sarama.ConsumerError)

	for _, topic := range topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}

		partitions, _ := master.Partitions(topic)
		consumer, err := master.ConsumePartition(topic, partitions[0], sarama.OffsetOldest)
		if nil != err {
			fmt.Printf("Topic %v Partitions: %v", topic, partitions)
			panic(err)
		}

		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					consumerErrors <- consumerError
					fmt.Println("consumerError: ", consumerError.Err)

				case msg := <-consumer.Messages():
					consumerMessages <- msg
				}
			}
		}(topic, consumer)
	}

	return consumerMessages, consumerErrors
}
