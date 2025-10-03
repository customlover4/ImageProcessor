package picturer

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"imager/internal/entities/picture"
	"os"

	"github.com/disintegration/imaging"
)


func SaveFile(pc picture.Picture, dst *image.NRGBA, min *image.NRGBA) error {
	const op = "internal.service.SaveFile"

	outputFName := fmt.Sprintf(
		"%s/%d.%s", picture.ProcessedPicturesFolder, pc.ID, pc.Extension,
	)
	outputFile, err := os.Create(outputFName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	outputMinFName := fmt.Sprintf(
		"%s/%d_%s.%s", picture.ProcessedPicturesFolder,
		pc.ID, picture.MiniaturePostfix, pc.Extension,
	)
	outputMinFile, err := os.Create(outputMinFName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = outputMinFile.Close()
	}()

	switch pc.Extension {
	case "jpg", "jpeg":
		err = jpeg.Encode(outputFile, dst, &jpeg.Options{Quality: 90})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = jpeg.Encode(outputMinFile, min, &jpeg.Options{Quality: 90})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case "png":
		err = png.Encode(outputFile, dst)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = png.Encode(outputMinFile, min)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func ProcessDefaultImage(pc picture.Picture) error {
	const op = "internal.service.ProcessDefaultImage"

	fileName := fmt.Sprintf(
		"%s/%d.%s", picture.PicturesFolder, pc.ID, pc.Extension,
	)

	src, err := imaging.Open(fileName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	watermarkTMP, err := imaging.Open(WatermarkPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	watermark := imaging.Resize(watermarkTMP, 400, 200, imaging.Lanczos)

	tmp := imaging.Overlay(
		src, watermark, image.Pt(
			src.Bounds().Dx()-watermark.Bounds().Dx()-20,
			src.Bounds().Dy()-watermark.Bounds().Dy()-20,
		), 1.0,
	)

	dst := imaging.Resize(tmp, 1200, 600, imaging.Lanczos)

	min := imaging.Thumbnail(dst, 200, 200, imaging.Lanczos)
	err = SaveFile(pc, dst, min)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
