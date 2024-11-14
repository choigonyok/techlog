package env

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

// The type of a variable's value
type VarType byte

const (
	STRING VarType = iota
	BOOL
	INT
	FLOAT
	DURATION
	OTHER
)

const (
	DB_PASSWORD    = "DB_PASSWORD"
	DB_MASTER_USER = "DB_MASTER_USER"
	DB_SLAVE_USER  = "DB_SLAVE_USER"
	DB_PORT        = "DB_PORT"
	DB_HOST_READ   = "DB_HOST_READ"
	DB_HOST_WRITE  = "DB_HOST_WRITE"
	DB_NAME        = "DB_NAME"
	GITHUB_TOKEN   = "GITHUB_TOKEN"
)

var EnvVars map[string]string
var mutex sync.Mutex

func Register(name string, description string) {
	mutex.Lock()
	fmt.Println("GG")
	mutex.Unlock()
}

func LoadDotEnvFile(path string) {
	godotenv.Load(path + ".env")

}

func GetEnvVarList() {

}
