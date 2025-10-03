package storage

import (
	"context"
	"errors"
	"imager/internal/entities/picture"

	"github.com/segmentio/kafka-go"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrNotAffected = errors.New("not affected")
)

type db interface {
	CreatePicture(pc picture.Picture) (uint64, error)
	Picture(id uint64) (picture.Picture, error)
	UpdatePictureStatus(id uint64, status string) error
	DeletePicture(id uint64) error

	Shutdown()
}

type broker interface {
	NewPicture(id uint64) error
	StartConsuming(ctx context.Context, out chan<- kafka.Message)
	Commit(msg kafka.Message) error

	Shutdown()
}

type Storage struct {
	db
	broker
}

func New(db db, broker broker) *Storage {
	return &Storage{
		db:     db,
		broker: broker,
	}
}

func (s *Storage) Shutdown() {
	s.db.Shutdown()
	s.broker.Shutdown()
}
