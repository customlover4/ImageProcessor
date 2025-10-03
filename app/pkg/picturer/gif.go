package picturer

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"imager/internal/entities/picture"
	"os"

	"github.com/disintegration/imaging"
)

func cloneRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}

// extractFullFrames — корректная распаковка анимации
func extractFullFrames(g *gif.GIF) []*image.RGBA {
	if len(g.Image) == 0 {
		return nil
	}

	// Убеждаемся, что размеры заданы
	if g.Config.Width == 0 || g.Config.Height == 0 {
		maxX, maxY := 0, 0
		for _, frame := range g.Image {
			b := frame.Bounds()
			if b.Max.X > maxX {
				maxX = b.Max.X
			}
			if b.Max.Y > maxY {
				maxY = b.Max.Y
			}
		}
		g.Config.Width = maxX
		g.Config.Height = maxY
	}

	width, height := g.Config.Width, g.Config.Height
	// Фон: БЕЛЫЙ и НЕПРОЗРАЧНЫЙ (можно заменить на color.RGBA{0,0,0,255} для чёрного)
	bg := color.RGBA{255, 255, 255, 255}

	buffer := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(buffer, buffer.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	var restore *image.RGBA
	var result []*image.RGBA

	for i, frame := range g.Image {
		if i < len(g.Disposal) && g.Disposal[i] == gif.DisposalPrevious {
			restore = cloneRGBA(buffer)
		}

		currentFrame := cloneRGBA(buffer)
		draw.Draw(currentFrame, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		result = append(result, currentFrame)

		if i < len(g.Disposal) {
			switch g.Disposal[i] {
			case gif.DisposalBackground:
				draw.Draw(buffer, frame.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
			case gif.DisposalPrevious:
				if restore != nil {
					draw.Draw(buffer, buffer.Bounds(), restore, image.Point{}, draw.Src)
				}
			}
		}
	}

	return result
}

// processFrames — основная функция обработки
func processFrames(g *gif.GIF, newW, newH int, watermark image.Image) *gif.GIF {
	fullFrames := extractFullFrames(g)
	wmRGBA := imaging.Clone(watermark)

	var newFrames []*image.Paletted

	for _, frame := range fullFrames {
		// Изменяем размер
		resized := imaging.Resize(frame, newW, newH, imaging.Lanczos)

		// Накладываем водяной знак
		wmBounds := wmRGBA.Bounds()
		x := newW - wmBounds.Dx() - 10
		y := newH - wmBounds.Dy() - 10
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		watermarked := imaging.Overlay(resized, wmRGBA, image.Pt(x, y), 1.0)

		// Конвертируем в Paletted с палитрой Plan9 (без прозрачности)
		bounds := image.Rect(0, 0, newW, newH)
		paletted := image.NewPaletted(bounds, palette.Plan9)
		draw.FloydSteinberg.Draw(paletted, bounds, watermarked, image.Point{})

		newFrames = append(newFrames, paletted)
	}

	return &gif.GIF{
		Image:     newFrames,
		Delay:     g.Delay,
		LoopCount: g.LoopCount,
		Disposal:  make([]byte, len(newFrames)), // все кадры полные
		Config: image.Config{
			Width:      newW,
			Height:     newH,
			ColorModel: color.Palette(palette.Plan9),
		},
	}
}

func ProcessGIF(pc picture.Picture) error {
	const op = "internal.service.ProcessGIF"

	fileName := fmt.Sprintf(
		"%s/%d.%s", picture.PicturesFolder, pc.ID, pc.Extension,
	)

	src, err := imaging.Open(fileName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	min := imaging.Thumbnail(src, 200, 200, imaging.Lanczos)

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	g, err := gif.DecodeAll(file)
	_ = file.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	watermarkTMP, err := imaging.Open(WatermarkPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	watermark := imaging.Resize(watermarkTMP, 800, 600, imaging.Lanczos)

	dst := processFrames(g, 800, 600, watermark)

	outputFName := fmt.Sprintf(
		"%s/%d.%s", picture.ProcessedPicturesFolder, pc.ID, pc.Extension,
	)
	outputMinFName := fmt.Sprintf(
		"%s/%d_%s.%s", picture.ProcessedPicturesFolder,
		pc.ID, picture.MiniaturePostfix, pc.Extension,
	)
	outFile, _ := os.Create(outputFName)
	defer outFile.Close()
	if err := gif.EncodeAll(outFile, dst); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	outMinFile, _ := os.Create(outputMinFName)
	defer outMinFile.Close()
	if err := gif.Encode(outMinFile, min, nil); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
