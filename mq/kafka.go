package mq

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

var (
	FollowConsumerMsg   chan string
	FollowProducerMsg   chan string
	FavoriteConsumerMsg chan string
	FavoriteProducerMsg chan string
)

func init() {
	FollowConsumerMsg = make(chan string)
	FollowProducerMsg = make(chan string)
	FavoriteConsumerMsg = make(chan string)
	FavoriteProducerMsg = make(chan string)
}

//follow and favorite，总共需要两个topic，partition都只能是一个
func InitKafka() {
	go producer(os.Getenv("FOLLOW_TOPIC"), 0, FollowProducerMsg)
	go producer(os.Getenv("FAVORITE_TOPIC"), 0, FavoriteProducerMsg)
	go consumer(os.Getenv("FOLLOW_TOPIC"), 0, FollowConsumerMsg)
	go consumer(os.Getenv("FAVORITE_TOPIC"), 0, FavoriteConsumerMsg)
}

func getKafkaConn(topic string, part int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", os.Getenv("KAFKA_HOST"), topic, part)
	if err != nil {
		return nil, err
	}
	err = conn.SetWriteDeadline(time.Now())
	if err != nil {
		return nil, err
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	}
	return conn, err
}

func producer(topic string, part int, ch <-chan string) {
	conn, err := getKafkaConn(topic, part)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	for {
		msg := <-ch
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func consumer(topic string, part int, ch chan<- string) {
	conn, err := getKafkaConn(topic, part)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	bytes := make([]byte, 1024)
	for {
		n, err := conn.Read(bytes)
		if err != nil {
			log.Fatalln(err)
		}
		ch <- string(bytes[:n])
	}
}
