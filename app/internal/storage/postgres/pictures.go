package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"imager/internal/entities/picture"
	"imager/pkg/errs"
)

func (p *Postgres) CreatePicture(pc picture.Picture) (uint64, error) {
	const op = "internal.storage.postgres.CreatePicture"

	q := fmt.Sprintf(
		"insert into %s (extension, status) values ($1, $2) returning id;",
		PicturesTable,
	)

	var id uint64
	row, err := p.db.QueryRowWithRetry(
		context.Background(), strategy,
		q, pc.Extension, picture.StatusProcessing,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errs.ErrDBNotFound
	} else if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *Postgres) Picture(id uint64) (picture.Picture, error) {
	const op = "internal.storage.postgres.Picture"

	q := fmt.Sprintf("select * from %s where id = $1;", PicturesTable)

	var pc picture.Picture
	row, err := p.db.QueryRowWithRetry(context.Background(), strategy, q, id)
	if err != nil {
		return pc, fmt.Errorf("%s: %w", op, err)
	}

	err = row.Scan(&pc.ID, &pc.Extension, &pc.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return picture.Picture{}, errs.ErrDBNotFound
	} else if err != nil {
		return picture.Picture{}, fmt.Errorf("%s: %w", op, err)
	}

	return pc, nil
}

func (p *Postgres) UpdatePictureStatus(id uint64, status string) error {
	const op = "internal.storage.postgres.UpdatePictureStatus"

	q := fmt.Sprintf(
		"update %s set status = $1 where id = $2;", PicturesTable,
	)

	res, err := p.db.ExecWithRetry(
		context.Background(), strategy, q, status, id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if affected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else if affected == 0 {
		return errs.ErrDBNotAffected
	}

	return nil
}

func (p *Postgres) DeletePicture(id uint64) error {
	const op = "internal.storage.postgres.DeletePicture"

	q := fmt.Sprintf("delete from %s where id = $1;", PicturesTable)

	res, err := p.db.ExecWithRetry(context.Background(), strategy, q, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if affected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else if affected == 0 {
		return errs.ErrDBNotAffected
	}

	return nil
}
