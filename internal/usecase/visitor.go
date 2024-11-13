package usecase

import (
	"database/sql"

	repo "github.com/choigonyok/techlog/internal/repository"
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

func (u *VisitorUsecase) GetTodayCount() (int, error) {
	today, err := u.visitorRepository.GetToday()
	if err == sql.ErrNoRows {
		u.visitorRepository.CreateToday()
		return 1, nil
	}
	if err := u.visitorRepository.UpdateToday(today + 1); err != nil {
		return 0, err
	}

	return today + 1, nil
}

func (u *VisitorUsecase) GetTotalCount() (int, error) {
	counts, err := u.visitorRepository.GetTotal()
	total := 0
	if err == sql.ErrNoRows {
		return 0, nil
	}

	for _, c := range *counts {
		total += c
	}

	return total, err
}
