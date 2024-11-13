package usecase

import (
	"fmt"
	"time"

	repo "github.com/choigonyok/techlog/internal/repository"
	_ "github.com/lib/pq"
)

type VisitorUsecase struct {
	repositories map[int]repo.Repository
}

func NewVisitorUsecase() *VisitorUsecase {
	return &VisitorUsecase{
		repositories: setRepositories(),
	}
}

func (u *VisitorUsecase) GetVisitorCount() int {
	t := u.repositories[repo.VISITOR_REPOSITORY].Get("count", "date", time.Now().Format("2006/01/02")).(int)
	return t
}

func setRepositories() map[int]repo.Repository {
	m := make(map[int]repo.Repository)
	fmt.Println("TEST1")
	m[repo.VISITOR_REPOSITORY] = repo.NewVisitorRepository()
	fmt.Println("TEST2")
	return m
}
