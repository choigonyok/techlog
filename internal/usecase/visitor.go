package usecase

import (
	repo "github.com/choigonyok/techlog/internal/repository"
	"github.com/choigonyok/techlog/pkg/time"
	_ "github.com/lib/pq"
)

type VisitorUsecase struct {
	visitorRepository *repo.VisitorRepository
}

func NewVisitorUsecase() *VisitorUsecase {
	return &VisitorUsecase{
		visitorRepository: repo.NewVisitorRepository(),
	}
}

func (u *VisitorUsecase) GetTodayCount() int {
	t := u.visitorRepository.Get("count", "date", time.GetCurrentTimeByFormat("2006/01/02")).(int)
	return t
}

func (u *VisitorUsecase) GetTotalCount() (int, error) {
	counts, err := u.visitorRepository.GetAllCount()
	total := 0

	for _, c := range *counts {
		total += c
	}

	return total, err
}
