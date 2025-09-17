package controllers

import (
	"album/api"
	"album/app/models"
	"album/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct{}

func (a *AlbumHandler) CreateAlbum(c *gin.Context) {
	var requestBody api.CreateAlbumJSONRequestBody
	if err := c.ShouldBind(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
	}

	createAlbum, err := models.CreateAlbum(
		requestBody.Title,
		requestBody.ReleaseDate.Time,
		string(requestBody.Category.Name),
	)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createAlbum)
}

func (a *AlbumHandler) GetAlbumById(c *gin.Context, ID int) {
	album, err := models.GetAlbum(ID)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, album)
}

func (a *AlbumHandler) UpdateAlbumById(c *gin.Context, ID int) {
	var requestBody api.UpdateAlbumByIdJSONRequestBody
	if err := c.ShouldBind(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	album, err := models.GetAlbum(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	if requestBody.Category != nil {
		album.Category.Name = string(requestBody.Category.Name)
	}
	if requestBody.Title != nil {
		album.Title = *requestBody.Title
	}

	if err := album.Save(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

func (a *AlbumHandler) DeleteAlbumById(c *gin.Context, ID int) {
	album := models.Album{ID: ID}

	if err := album.Delete(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
