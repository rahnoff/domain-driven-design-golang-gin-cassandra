package repository

import (
	"golang-gin-cassandra/src/domain/videos/model"
	"golang-gin-cassandra/src/utils/errors"
)

type VideoRepository interface {
	GetByTitle(videoTitle string) (*model.Video, *errors.RestErr)
	Create(video model.Video) (*model.Video, *errors.RestErr)
}
