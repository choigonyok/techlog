package service

import "github.com/choigonyok/techlog/pkg/database"

type Service struct {
	provider database.Provider
}

// NewService creates new service to connect handler with database provider
func NewService(prov database.Provider) *Service {
	return &Service{
		provider: prov,
	}
}
