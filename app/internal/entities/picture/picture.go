package picture

const (
	PicturesFolder          = "pictures"
	ProcessedPicturesFolder = "processed_pictures"

	MiniaturePostfix = "min"

	StatusComplete   = "complete"
	StatusProcessing = "processing"
	StatusError      = "error"
)

type Picture struct {
	ID        uint64
	Extension string
	Status    string
}

func ValidateExtension(extension string) bool {
	switch extension {
	case "jpg", "png", "gif":
		return true
	}

	return false
}

type PictureInfoPage struct {
	Status        string
	FilePath      string
	Info          string
	ThumbnailPath string
}

var (
	ProcessingPage = func() PictureInfoPage {
		return PictureInfoPage{
			Status:        StatusProcessing,
			FilePath:      "/processed/default.jpg",
			Info:          "image still in processing, wait",
			ThumbnailPath: "-",
		}
	}
	ErrorProcessingPage = func() PictureInfoPage {
		return PictureInfoPage{
			Status:        StatusError,
			FilePath:      "/processed/default.jpg",
			Info:          "ERROR ON PROCESSING IMAGE, TRY AGAIN",
			ThumbnailPath: "-",
		}
	}
	CompleteProcessingPage = func(filePath, minPath string) PictureInfoPage {
		return PictureInfoPage{
			Status:        StatusComplete,
			FilePath:      filePath,
			Info:          "image processed, this is result",
			ThumbnailPath: minPath,
		}
	}
)
