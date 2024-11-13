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
const TMP_VISITOR_DB_DATABASE = "visitor_db"
const TMP_POST_DB_DATABASE = "post_db"
const TMP_DB_PORT = 5432
