package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"imager/internal/entities/picture"
	"imager/internal/storage"
	"imager/pkg/picturer"

	"github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/zlog"
)

func logErr(op string, err error) {
	zlog.Logger.Error().Err(fmt.Errorf("%s: %w", op, err)).Send()
}

func Processing(pc picture.Picture) error {
	const op = "internal.service.Processing"

	if pc.Extension != "gif" {
		err := picturer.ProcessDefaultImage(pc)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		err := picturer.ProcessGIF(pc)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Service) processImage(msg kafka.Message) {
	const op = "internal.service.ProcessImage"

	var id uint64
	buffer := bytes.NewBuffer(msg.Value)
	binary.Read(buffer, binary.LittleEndian, &id)

	pc, err := s.Picture(id)
	if errors.Is(err, storage.ErrNotFound) {
		_ = s.storage.Commit(msg)
	} else if err != nil {
		logErr(op, err)
		return
	}

	if pc.Status == picture.StatusComplete {
		_ = s.storage.Commit(msg)
		return
	}

	err = Processing(pc)
	if err != nil {
		// try to set error status for picture, if can't don't commit kafka.
		err := s.UpdatePictureStatus(id, true)
		if err != nil {
			logErr(op, err)
			return
		}
		_ = s.storage.Commit(msg)
		logErr(op, err)
		return
	}

	err = s.UpdatePictureStatus(id, false)
	if err != nil {
		logErr(op, err)
		return
	}

	err = s.storage.Commit(msg)
	if err != nil {
		logErr(op, err)
	}
}

func (s *Service) StartConsuming(ctx context.Context, out chan kafka.Message) {
	go func() {
		for {
			select {
			case v := <-out:
				go s.processImage(v)
			case <-ctx.Done():
				return
			}
		}
	}()
}
