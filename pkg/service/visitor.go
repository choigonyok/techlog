package service

import (
	"github.com/choigonyok/techlog/pkg/database"
)

type VisitorService struct {
	provider database.Provider
}

// NewService creates new service to connect handler with database provider
func NewService(prov database.Provider) *VisitorService {
	return &VisitorService{
		provider: prov,
	}
}

// GetDate returns stored visitor/date
func (svc *VisitorService) GetDate() (string, error) {
	m, err := svc.provider.GetVisitor()
	return m.Date.Format("2006-01-02"), err
}

// ResetToday resets stored visitor/today to 1 and visitor/date to today
func (svc *VisitorService) ResetToday(today string) error {
	return svc.provider.ResetVisitorTodayAndDate(today)
}

// GetCounts returns stored visitor/today and visitor/total
func (svc *VisitorService) GetCounts() (int, int, error) {
	m, err := svc.provider.GetVisitor()
	return m.Today, m.Total, err
}

// AddToday updates stored visitor/today to visitor/today + 1
func (svc *VisitorService) AddToday() error {
	m, err := svc.provider.GetVisitor()
	newToday := m.Today + 1
	if err != nil {
		return err
	}
	return svc.provider.UpdateVisitorToday(newToday)
}
