package handlers

import (
	"errors"
	"fmt"
	"imager/internal/entities/picture"
	"imager/internal/service"
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func MainPage(ctx *ginext.Context) {
	pictureID := ctx.Query("picture-id")
	if pictureID != "" {
		ctx.Redirect(http.StatusSeeOther, "/image/"+pictureID)
		return
	}

	ctx.HTML(http.StatusOK, "index.html", nil)
}

func UploadPage(ctx *ginext.Context) {
	ctx.HTML(http.StatusOK, "upload.html", nil)
}

func GetFile(srv srv) ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		const op = "internal.web.handlers.GetFile"

		ctx.Writer.Header().Set("Content-Type", "text/html")

		tmp := ctx.Param("id")
		id, err := strconv.ParseUint(tmp, 10, 64)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "400.html", nil)
			return
		}

		pc, err := srv.Picture(id)
		if errors.Is(err, service.ErrNotValidData) {
			ctx.HTML(http.StatusServiceUnavailable, "503.html", nil)
			return
		} else if errors.Is(err, service.ErrNotFound) {
			ctx.HTML(http.StatusNotFound, "404.html", nil)
			return
		} else if err != nil {
			zlog.Logger.Error().Err(fmt.Errorf("%s: %w", op, err))
			ctx.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}

		if pc.Status == picture.StatusProcessing {
			ctx.HTML(
				http.StatusOK, "picture.html", picture.ProcessingPage(),
			)
			return
		}
		if pc.Status == picture.StatusError {
			ctx.HTML(
				http.StatusOK, "picture.html", picture.ErrorProcessingPage(),
			)
			return
		}

		filePath := fmt.Sprintf("/processed/%d.%s", pc.ID, pc.Extension)
		ThumbnailPath := fmt.Sprintf(
			"/processed/%d_%s.%s", pc.ID, picture.MiniaturePostfix,
			pc.Extension,
		)
		ctx.HTML(
			http.StatusOK, "picture.html", picture.CompleteProcessingPage(
				filePath, ThumbnailPath,
			),
		)
	}
}
