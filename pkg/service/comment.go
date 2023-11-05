package service

import (
	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/model"
)

// GetComments returns every comments in blog
func (svc *Service) GetComments() ([]model.Comment, error) {
	comments, err := svc.provider.GetComments()
	return data.UnmarshalCommentToDatabaseFmt(comments), err
}

// GetCommentsByPostID returns every comments in a post by post id
func (svc *Service) GetCommentsByPostID(postID string) ([]model.Comment, error) {
	comments, err := svc.provider.GetCommentsByPostID(postID)
	return data.UnmarshalCommentToDatabaseFmt(comments), err
}

// GetCommentPasswordByCommentID returns specific comment's password
func (svc *Service) GetCommentPasswordByCommentID(commentID string) (string, error) {
	password, err := svc.provider.GetCommentPasswordByCommentID(commentID)
	return data.DecodeBase64(password), err
}

// DeleteCommentByCommentID deletes a stored comment by comment id
func (svc *Service) DeleteCommentByCommentID(commentID string) error {
	return svc.provider.DeleteCommentByCommentID(commentID)
}

// CreateComment stores new comment into database
func (svc *Service) CreateComment(comment model.Comment) error {
	var admin string
	comment = data.MarshalCommentToDatabaseFmt(comment)

	if comment.Admin {
		admin = "1"
	} else {
		admin = "0"
	}
	return svc.provider.CreateComment(comment, admin)
}
