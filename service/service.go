package service

import (
	"SOTBI/notifier/service/sender"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"SOTBI/notifier/config"
	"github.com/Shopify/sarama"
)

type Service struct {
	appConfig config.Config

	saramaClient sarama.Client
	isRunning    bool

	notifySender sender.Sender
}

func (s *Service) New(appConfig config.Config) *Service {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = "notifier"
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Producer.Return.Successes = true

	client, err := sarama.NewClient([]string{appConfig.RedPanda.Host}, saramaConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	return &Service{
		saramaClient: nil,
		isRunning:    false,
	}
}

func (s *Service) Run(appConfig config.Config) error {
	if s.isRunning {
		return errors.New("consumer already running")
	}

	master, err := sarama.NewConsumerFromClient(s.saramaClient)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topics, err := master.Topics()
	if err != nil {
		panic(err)
	}

	messages, errChan := consume(topics, master)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Get signal for finish
	doneCh := make(chan struct{})

	s.isRunning = true
	go func() {
		for {
			select {
			case msg := <-messages:
				fmt.Println(
					"Received messages",
					string(msg.Key),
					string(msg.Value),
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
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
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
		// this only consumes partition no 1, you would probably want to consume all partitions
		consumer, err := master.ConsumePartition(topic, partitions[0], sarama.OffsetOldest)
		if nil != err {
			fmt.Printf("Topic %v Partitions: %v", topic, partitions)
			panic(err)
		}
		fmt.Println(" Start consuming topic ", topic)
		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					consumerErrors <- consumerError
					fmt.Println("consumerError: ", consumerError.Err)

				case msg := <-consumer.Messages():
					consumerMessages <- msg
					fmt.Println("Got message on topic ", topic, msg.Value)
				}
			}
		}(topic, consumer)
	}

	return consumerMessages, consumerErrors
}
