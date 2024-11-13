package repository

type Repository interface {
	Get(column, conditionKey, conditionValue string) any
}

const (
	VISITOR_REPOSITORY = iota
)

const TMP_DB_HOST = "localhost"
const TMP_DB_PASSWORD = "tester1234"
const TMP_DB_USERNAME = "postgres"
const TMP_DB_DATABASE = "test_db"
const TMP_DB_PORT = 5432
