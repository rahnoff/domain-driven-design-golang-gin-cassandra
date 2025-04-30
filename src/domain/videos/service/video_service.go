package service

import (
	"errors"
	"fmt"
	"golang-gin-cassandra/src/domain/videos/model"
	"golang-gin-cassandra/src/domain/videos/repository"
	. "golang-gin-cassandra/src/utils/errors"
	"strings"
)

type VideoService interface {
	GetByTitle(videoTitle string) (*model.Video, *RestErr)
	Create(video model.Video) (*model.Video, *RestErr)
}

type videoService struct {
	repository repository.VideoRepository
}

func (s *videoService) Create(video model.Video) (*model.Video, *RestErr) {
	err := video.ValidateVideo();
	if err != nil {
		return nil, err
	}
	return s.repository.Create(video)
}

func NewService(repository repository.VideoRepository) VideoService {
	return &videoService{
		repository: repository,
	}
}

func (s *videoService) GetByTitle(videoTitle string) (*model.Video, *RestErr) {
	videoTitle = strings.TrimSpace(videoTitle)
	if videoTitle == "" {
		return nil, NewBadRequestError("Title can't be empty")
	}
	video, err := s.repository.GetByTitle(videoTitle)
	if err != nil {
		videoNotFoundErr := fmt.Sprintf("Video not found for title %s", videoTitle)
		return nil, NewInternalServerError(videoNotFoundErr, errors.New("here"))
	}
	return video, nil
}
