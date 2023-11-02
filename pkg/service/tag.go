package service

import (
	"strings"

	"github.com/choigonyok/techlog/pkg/model"
)

func (svc *Service) GetEveryTag() ([]string, error) {
	result := []string{}
	tags, err := svc.provider.GetEveryTag()
	m := make(map[string]bool)

	for _, v := range tags {
		separateTags := strings.Split(v.Tags, " ")
		for _, tag := range separateTags {
			pureTag := strings.ReplaceAll(tag, " ", "")
			if !m[pureTag] {
				result = append(result, pureTag)
				m[tag] = true
			}
		}
	}
	return result, err
}

func (svc *Service) GetEveryCardByTag(tag string) ([]model.PostCard, error) {
	switch tag {
	case "ALL":
		return svc.provider.GetEveryCard()
	default:
		return svc.provider.GetEveryCardByTag(tag)
	}
}
