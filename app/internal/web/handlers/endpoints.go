package handlers

import (
	"errors"
	"fmt"
	"imager/internal/entities/picture"
	"imager/internal/entities/response"
	"imager/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func logErr(op string, err error) {
	zlog.Logger.Error().Err(fmt.Errorf("%s: %w", op, err)).Send()
}

func UploadFile(srv srv) ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		const op = "internal.web.handlers.uploadFile"

		ctx.Writer.Header().Set("Content-Type", "application/json")

		formFile, err := ctx.FormFile("image")
		if err != nil {
			logErr(op, err)
			ctx.JSON(http.StatusInternalServerError, response.Error(
				"can't get your image from form",
			))
			return
		}

		contentType := strings.Split(formFile.Header.Get("Content-Type"), "/")
		if contentType[0] != "image" {
			ctx.JSON(http.StatusBadRequest, response.Error(
				"you can upload only images",
			))
			return
		}

		tmp := strings.Split(formFile.Filename, ".")
		extension := tmp[len(tmp)-1]
		if !picture.ValidateExtension(extension) {
			ctx.JSON(http.StatusBadRequest, response.Error(
				"unknown extension of file",
			))
			return
		}

		id, err := srv.CreatePicture(
			picture.Picture{Extension: extension},
		)
		if err != nil {
			logErr(op, err)
			ctx.JSON(http.StatusInternalServerError, response.Error(
				"internal server error",
			))
			return
		}

		fileName := fmt.Sprintf(
			"%s/%d.%s", picture.PicturesFolder, id, extension,
		)
		err = ctx.SaveUploadedFile(formFile, fileName)
		if err != nil {
			logErr(op, err)
			ctx.JSON(http.StatusInternalServerError, response.Error(
				"can't upload this file",
			))
			return
		}

		err = srv.NewPicture(id)
		if err != nil {
			// Можно не удалять, но делаем ради экономии памяти.
			go func() {
				_ = srv.DeletePicture(id)
			}()
			logErr(op, err)
			ctx.JSON(http.StatusInternalServerError, response.Error(
				"internal server error",
			))
			return
		}

		ctx.JSON(http.StatusOK, response.OK(id))
	}
}

func DeleteFile(srv srv) ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		const op = "internal.web.handlers.DeleteFile"

		ctx.Writer.Header().Set("Content-Type", "application/json")

		tmp := ctx.Param("id")
		id, err := strconv.ParseUint(tmp, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Error(
				"wrong id of picture",
			))
			return
		}

		err = srv.DeletePicture(id)
		if errors.Is(err, service.ErrNotValidData) {
			ctx.JSON(http.StatusServiceUnavailable, response.Error(
				err.Error(),
			))
			return
		} else if errors.Is(err, service.ErrNotAffected) {
			ctx.JSON(http.StatusBadRequest, response.Error(
				"not found picture by this id",
			))
			return
		} else if err != nil {
			logErr(op, err)
			ctx.JSON(http.StatusInternalServerError, response.Error(
				"internal server error",
			))
			return
		}

		ctx.JSON(http.StatusOK, response.OK("ok"))
	}
}
