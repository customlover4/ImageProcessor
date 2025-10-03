package kafka

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	kafkaLib "github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

type Kafka struct {
	cons *kafka.Consumer
	prod *kafka.Producer
}

var (
	strategy = retry.Strategy{
		Attempts: 5,
		Delay:    1,
		Backoff:  1.25,
	}

	buffer = new(bytes.Buffer)
)

func New(brokers []string, topic string, groupID string) *Kafka {
	// for _, v := range brokers {
	// 	conn, err := kafkaLib.Dial("tcp", v)
	// 	if err != nil {
	// 		panic(fmt.Sprintf("❌ Не удалось подключиться: %v", err))
	// 	}
	// 	defer conn.Close()

	// 	conn.SetDeadline(time.Now().Add(10 * time.Second))
	// 	_, err = conn.ReadPartitions() // запрос метаданных
	// 	if err != nil {
	// 		panic(fmt.Sprintf("❌ Ошибка чтения топиков: %v", err))
	// 	}

	// }

	cons := kafka.NewConsumer(brokers, topic, groupID)
	prod := kafka.NewProducer(brokers, topic)

	return &Kafka{
		cons: cons,
		prod: prod,
	}
}

func (k *Kafka) Shutdown() {
	_ = k.cons.Close()
	_ = k.prod.Close()
}

func (k *Kafka) NewPicture(id uint64) error {
	const op = "internal.storage.kafka.NewPicture"

	binary.Write(buffer, binary.LittleEndian, id)
	bts := buffer.Bytes()
	err := k.prod.SendWithRetry(context.Background(), strategy, bts, bts)
	buffer.Reset()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (k *Kafka) Commit(msg kafkaLib.Message) error {
	const op = "internal.storage.kafka.Commit"

	err := k.cons.Commit(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (k *Kafka) StartConsuming(ctx context.Context, out chan<- kafkaLib.Message) {
	k.cons.StartConsuming(ctx, out, strategy)
}
