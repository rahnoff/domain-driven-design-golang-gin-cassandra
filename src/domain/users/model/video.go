package model

import "golang-gin-cassandra/src/utils/errors"


type Video struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	Language    string `json:"language"`
	StorageLink string `json:"storage_link"`
}

func (video *Video) ValidateVideo() *errors.RestErr {
	if (video.Title == "") {
	    return errors.NewInternalServerError("Title can't be blank", nil)
	}
	
	return nil
}
