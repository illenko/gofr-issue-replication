package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type transaction struct {
	ID       uuid.UUID `json:"id"`
	Datetime time.Time `json:"datetime"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Country  string    `json:"country"`
}

func main() {

	app := gofr.New()

	//app.POST("/transactions", func(ctx *gofr.Context) (interface{}, error) {
	//
	//	var data []transaction
	//
	//	err := ctx.Bind(&data)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	for _, t := range data {
	//		msg, err := json.Marshal(t)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		err = ctx.GetPublisher().Publish(ctx, "transactions", msg)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	return "Published", nil
	//})

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})

	if err != nil {
		panic(err)
	}

	app.POST("/transactions", func(ctx *gofr.Context) (interface{}, error) {

		var data []transaction

		err := ctx.Bind(&data)
		if err != nil {
			return nil, err
		}

		for _, t := range data {
			msg, err := json.Marshal(t)
			if err != nil {
				return nil, err
			}

			topic := "transactions"

			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          msg},
				nil,
			)
		}

		return "Published", nil
	})

	app.Run()
}
