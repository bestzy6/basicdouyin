package mq

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

var (
	FollowConsumerMsg chan string
	FollowProducerMsg chan string
	FollowNotifyMsg   chan struct{}
	//
	FavoriteConsumerMsg chan string
	FavoriteProducerMsg chan string
	FavoriteNotifyMsg   chan struct{}
	MaxLength           = 20
)

func init() {
	FollowConsumerMsg = make(chan string)
	FollowProducerMsg = make(chan string, MaxLength)
	FollowNotifyMsg = make(chan struct{})
	//
	FavoriteConsumerMsg = make(chan string)
	FavoriteProducerMsg = make(chan string, MaxLength)
	FavoriteNotifyMsg = make(chan struct{})
}

//follow and favorite，总共需要两个topic，partition都只能是一个
func InitKafka() {
	go producer(os.Getenv("FOLLOW_TOPIC"), 0, FollowProducerMsg)
	go producer(os.Getenv("FAVORITE_TOPIC"), 0, FavoriteProducerMsg)
	//
	go consumer(os.Getenv("FOLLOW_TOPIC"), 0, FollowConsumerMsg, FollowNotifyMsg)
	go consumer(os.Getenv("FAVORITE_TOPIC"), 0, FavoriteConsumerMsg, FavoriteNotifyMsg)
}

func producer(topic string, part int, ch <-chan string) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("KAFKA_HOST")),
		Topic:        topic,
		RequiredAcks: 1,
		Balancer:     &kafka.LeastBytes{},
		Async:        false, //使用同步的方式
	}
	defer w.Close()
	for {
		msg := <-ch
		err := w.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(msg),
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func consumer(topic string, part int, ch chan<- string, notify chan struct{}) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_HOST")},
		GroupID:   "consumer-group-id",
		Topic:     topic,
		Partition: part,
	})
	defer reader.Close()
	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		//写入通道
		ch <- string(msg.Value)
		//检测是否成功处理消息，成功则向kafka提交offset
		<-notify
		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Fatalln("failed to commit messages:", err)
		}
	}
}
