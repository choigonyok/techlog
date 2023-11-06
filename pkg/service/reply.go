package service

import (
	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/model"
)

func (s *Service) GetRepliesByPostID(postID string) ([]model.Reply, error) {
	replies, err := s.provider.GetRepliesByPostID(postID)
	return data.UnmarshalReplyToDatabaseFmt(replies), err
}

func (s *Service) GetReplyPasswordByReplyID(replyID string) (string, error) {
	password, err := s.provider.GetReplyPasswordByReplyID(replyID)
	return data.DecodeBase64(password), err
}

func (s *Service) DeleteReplyByReplyID(replyID string) error {
	return s.provider.DeleteReplyByReplyID(replyID)
}

func (s *Service) CreateReply(reply model.Reply) error {
	reply = data.MarshalReplyToDatabaseFmt(reply)
	reply.WriterPW = data.EncodeBase64(reply.WriterPW)
	return s.provider.CreateReply(reply)
}
