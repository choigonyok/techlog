package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/choigonyok/techlog/pkg/time"
)

type VisitorRepository struct {
	db *sql.DB
}

func NewVisitorRepository() *VisitorRepository {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", TMP_DB_HOST, TMP_DB_PORT, TMP_DB_USERNAME, TMP_DB_PASSWORD, TMP_VISITOR_DB_DATABASE)

	cli, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("TEST CREATE POSTGRESQL CLIENT ERROR: ", err)
	}

	return &VisitorRepository{
		db: cli,
	}
}

func (repo *VisitorRepository) GetToday() (int, error) {
	r := repo.db.QueryRow(`SELECT count FROM visitor WHERE date = '` + time.GetCurrentTimeByFormat("2006/01/02") + `'`)

	count := 0
	err := r.Scan(&count)

	return count, err
}

func (repo *VisitorRepository) GetTotal() (*[]int, error) {
	counts := []int{}
	r, err := repo.db.Query(`SELECT count FROM visitor`)
	if err != nil {
		return nil, err
	}

	for r.Next() {
		c := 0
		r.Scan(&c)
		counts = append(counts, c)
	}

	return &counts, nil
}

func (repo *VisitorRepository) UpdateToday(today int) error {
	_, err := repo.db.Exec(`UPDATE visitor SET count = ` + strconv.Itoa(today) + ` WHERE date = '` + time.GetCurrentTimeByFormat("2006/01/02") + `'`)
	return err
}

func (repo *VisitorRepository) CreateToday() error {
	_, err := repo.db.Exec(`INSERT INTO visitor (date, count) VALUES ('` + time.GetCurrentTimeByFormat("2006/01/02") + `', 1)`)
	return err
}
