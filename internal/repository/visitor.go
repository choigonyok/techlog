package repository

import (
	"database/sql"
	"fmt"
)

type VisitorRepository struct {
	db *sql.DB
}

func NewVisitorRepository() *VisitorRepository {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", TMP_DB_HOST, TMP_DB_PORT, TMP_DB_USERNAME, TMP_DB_PASSWORD, TMP_DB_DATABASE)

	cli, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("TEST CREATE POSTGRESQL CLIENT ERROR: ", err)
	}

	return &VisitorRepository{
		db: cli,
	}
}

func (repo *VisitorRepository) Get(column, conditionKey, conditionValue string) any {
	r := repo.db.QueryRow(`SELECT ` + column + ` FROM visitor WHERE ` + conditionKey + ` = '` + conditionValue + `'`)

	count := 0
	r.Scan(&count)

	return count
}

func (repo *VisitorRepository) GetAllCount() (*[]int, error) {
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
