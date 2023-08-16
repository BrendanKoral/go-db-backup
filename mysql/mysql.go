package mysql

import (
	"backup-db/fs"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jamf/go-mysqldump"
	"log"
	"sync"
)

type Config map[string]SqlConfig

type SqlConfig struct {
	Username  string
	Password  string
	DBName    string
	DBAddress string
	DBPort    string
	DumpDir   string
}

var C Config

func DumpDB(sqlCfg SqlConfig, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Beginning backup for ", sqlCfg.DBName)
	log.Println("Connecting to MySQL instance")

	// Open connection to database
	cfg := mysql.NewConfig()
	cfg.User = sqlCfg.Username
	cfg.Passwd = sqlCfg.Password
	cfg.DBName = sqlCfg.DBName
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%s", sqlCfg.DBAddress, sqlCfg.DBPort)
	cfg.ParseTime = true

	log.Println("Creating backup dir " + sqlCfg.DumpDir)
	err := fs.CreateBackupDir(sqlCfg.DumpDir)

	if err != nil {
		log.Fatal(err)
	}

	dumpFileName := fmt.Sprintf("%s-20060102T150405", sqlCfg.DBName)

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, sqlCfg.DumpDir, dumpFileName)

	defer func() {
		err = dumper.Close()
		if err != nil {
			log.Fatal("Error closing dumper", err)
		}
	}()

	if err != nil {
		log.Fatal("Error registering database to MySQL dump: ", err)
	}

	// Dump database to file
	errMsg := dumper.Dump()

	if errMsg != nil {
		log.Fatal("Error dumping: ", errMsg)
	}

	log.Printf("File %s successfully saved to %s", dumpFileName, sqlCfg.DumpDir)

	ch <- fmt.Sprintf("%s/%s.sql", sqlCfg.DumpDir, dumpFileName)
}
