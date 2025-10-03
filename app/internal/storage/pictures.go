package storage

import (
	"errors"
	"fmt"
	"imager/internal/entities/picture"
	"imager/pkg/errs"
)

func (s *Storage) CreatePicture(pc picture.Picture) (uint64, error) {
	const op = "internal.storage.CreatePicture"

	id, err := s.db.CreatePicture(pc)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Picture(id uint64) (picture.Picture, error) {
	const op = "internal.storage.Picture"

	pc, err := s.db.Picture(id)
	if errors.Is(err, errs.ErrDBNotFound) {
		return pc, ErrNotFound
	} else if err != nil {
		return pc, fmt.Errorf("%s: %w", op, err)
	}

	return pc, nil
}

func (s *Storage) UpdatePictureStatus(id uint64, status string) error {
	const op = "internal.storage.UpdatePictureStatus"

	err := s.db.UpdatePictureStatus(id, status)
	if errors.Is(err, errs.ErrDBNotAffected) {
		return ErrNotAffected
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeletePicture(id uint64) error {
	const op = "internal.storage.DeletePicture"

	err := s.db.DeletePicture(id)
	if errors.Is(err, errs.ErrDBNotAffected) {
		return ErrNotAffected
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
