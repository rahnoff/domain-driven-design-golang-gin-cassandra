package http

import (
	"github.com/gin-gonic/gin"
	"golang-gin-cassandra/src/domain/videos/model"
	"golang-gin-cassandra/src/domain/videos/service"
	"golang-gin-cassandra/src/utils/errors"
	"net/http"
	"strings"
)

type VideoHandler interface {
	GetByTitle(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type videoHandler struct {
	videoService service.VideoService
}

func (videoHandler videoHandler) GetByTitle(ctx *gin.Context) {
	videoTitle := strings.TrimSpace(ctx.Param("video_title"))
	video, err := videoHandler.videoService.GetByTitle(videoTitle)
	if err != nil {
		ctx.JSON(err.ErrStatus, err)
		return
	}
	ctx.JSON(http.StatusOK, video)
}

func NewHandler(videoService service.VideoService) VideoHandler {
	return &videoHandler{
		videoService: videoService,
	}
}

func (videoHandler *videoHandler) Create(ctx *gin.Context)  {
	var video model.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		restErr := errors.NewBadRequestError("Invalid json body")
		ctx.JSON(restErr.ErrStatus, restErr)
	}
	_, err := videoHandler.videoService.Create(video)
	if err != nil {
		ctx.JSON(err.ErrStatus, err)
		return
	}
	ctx.JSON(http.StatusCreated, video)
}
