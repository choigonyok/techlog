package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	databaseDriver            = "mysql"
	masterDatabaseServiceName = "mysql-ha-0"
	slave1DatabaseServiceName = "mysql-ha-1"
	slave2DatabaseServiceName = "mysql-ha-2"
	masterSVCName             = "mysql-ha-0.default.svc.cluster.local"
)

func main() {
	godotenv.Load(".env")
	var rootPassword = os.Getenv("DB_PASSWORD")
	var databaseName = os.Getenv("DB_NAME")

	var file string
	var position uint32
	var t1 string
	var t2 string
	var t3 string

	addresses, err := net.LookupHost(masterSVCName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// master setting
	master, err := sql.Open(databaseDriver, "root:"+rootPassword+"@tcp("+masterDatabaseServiceName+")/"+databaseName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer master.Close()

	err = master.QueryRow(`SHOW MASTER STATUS`).Scan(&file, &position, &t1, &t2, &t3)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// slave1 setting
	slave1, err := sql.Open(databaseDriver, "root:"+rootPassword+"@tcp("+slave1DatabaseServiceName+")/"+databaseName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer slave1.Close()

	_, err = slave1.Exec(`STOP REPLICA IO_THREAD FOR CHANNEL ''`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = slave1.Exec(fmt.Sprintf(`CHANGE MASTER TO MASTER_HOST='`+addresses[0]+`', MASTER_PORT=3306, MASTER_USER='replicas', MASTER_PASSWORD='password', MASTER_LOG_FILE='%s', MASTER_LOG_POS=%d, GET_MASTER_PUBLIC_KEY=1`, file, position))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = slave1.Exec(`SET GLOBAL SQL_SLAVE_SKIP_COUNTER = 1;`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = slave1.Exec(`START SLAVE`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// mysql-ha-2 setting
	slave2, err := sql.Open(databaseDriver, "root:"+rootPassword+"@tcp("+slave2DatabaseServiceName+")/"+databaseName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer slave2.Close()

	_, err = slave2.Exec(`STOP REPLICA IO_THREAD FOR CHANNEL ''`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = slave2.Exec(fmt.Sprintf(`CHANGE MASTER TO MASTER_HOST='`+addresses[0]+`', MASTER_PORT=3306, MASTER_USER='replicas', MASTER_PASSWORD='password', MASTER_LOG_FILE='%s', MASTER_LOG_POS=%d, GET_MASTER_PUBLIC_KEY=1`, file, position))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = slave2.Exec(`SET GLOBAL SQL_SLAVE_SKIP_COUNTER = 1;`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = slave2.Exec(`START SLAVE`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
