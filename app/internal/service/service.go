package service

import (
	"errors"
	"fmt"
	"imager/internal/entities/picture"
	"imager/internal/storage"
	"os"

	"github.com/segmentio/kafka-go"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrNotAffected  = errors.New("not affected")
	ErrNotValidData = errors.New("not valid data")
)

type str interface {
	NewPicture(id uint64) error
	Commit(msg kafka.Message) error

	CreatePicture(pc picture.Picture) (uint64, error)
	Picture(id uint64) (picture.Picture, error)
	UpdatePictureStatus(id uint64, status string) error
	DeletePicture(id uint64) error
}

type Service struct {
	storage str
}

func New(str str) *Service {
	return &Service{
		storage: str,
	}
}

func (s *Service) Shutdown() {
	// smth
}

func (s *Service) CreatePicture(pc picture.Picture) (uint64, error) {
	const op = "internal.service.CreatePicture"

	id, err := s.storage.CreatePicture(pc)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Service) NewPicture(id uint64) error {
	const op = "internal.service.NewPicture"

	err := s.storage.NewPicture(id)
	if err != nil {
		_ = s.DeletePicture(id)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Picture(id uint64) (picture.Picture, error) {
	const op = "internal.service.Picture"

	if id <= 0 {
		return picture.Picture{}, fmt.Errorf(
			"%w: %s", ErrNotValidData, "wrong id, should be > 0",
		)
	}

	pc, err := s.storage.Picture(id)
	if errors.Is(err, storage.ErrNotFound) {
		return pc, ErrNotFound
	} else if err != nil {
		return pc, fmt.Errorf("%s: %w", op, err)
	}

	return pc, nil
}

func (s *Service) UpdatePictureStatus(id uint64, haveError bool) error {
	const op = "internal.service.UpdatePictureStatus"

	if id <= 0 {
		return fmt.Errorf(
			"%w: %s", ErrNotValidData, "wrong id, should be > 0",
		)
	}

	var status string
	if haveError {
		status = picture.StatusError
	} else {
		status = picture.StatusComplete
	}

	err := s.storage.UpdatePictureStatus(id, status)
	if errors.Is(err, storage.ErrNotAffected) {
		return ErrNotAffected
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeletePicture(id uint64) error {
	const op = "internal.service.DeletePicture"

	if id <= 0 {
		return fmt.Errorf(
			"%w: %s", ErrNotValidData, "wrong id, should be > 0",
		)
	}

	pc, err := s.storage.Picture(id)
	if errors.Is(err, storage.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fileName := fmt.Sprintf(
		"%s/%d.%s", picture.PicturesFolder, pc.ID, pc.Extension,
	)
	_ = os.Remove(fileName)
	if pc.Status == picture.StatusComplete {
		fileName := fmt.Sprintf(
			"%s/%d.%s", picture.ProcessedPicturesFolder, pc.ID, pc.Extension,
		)
		_ = os.Remove(fileName)
		fileNameMin := fmt.Sprintf(
			"%s/%d_%s.%s", picture.ProcessedPicturesFolder,
			pc.ID, picture.MiniaturePostfix, pc.Extension,
		)
		_ = os.Remove(fileNameMin)
	}

	err = s.storage.DeletePicture(id)
	if errors.Is(err, storage.ErrNotAffected) {
		return ErrNotAffected
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
