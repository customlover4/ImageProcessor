package handlers

import (
	"imager/internal/entities/picture"
)

type srv interface {
	CreatePicture(pc picture.Picture) (uint64, error)
	NewPicture(id uint64) error
	Picture(id uint64) (picture.Picture, error)
	DeletePicture(id uint64) error
}
