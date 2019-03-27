package main

import (
	"database/sql"
	_ "database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/student_manage?maxAllowedPacket=0&interpolateParams=true")
	if err != nil {
		panic(err)
	}

	log.SetPrefix("Trace: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

func main() {
	var isHelp,isInit,isTest bool
	flag.BoolVar(&isHelp, "h", false, "this help")
	flag.BoolVar(&isInit, "i", false, "init database")
	flag.BoolVar(&isTest, "t", false, "test sql")
	flag.Parse()

	if isHelp {
		fmt.Println("This is a help message")
	} else if isInit {
		fmt.Println("This is a init command")
	} else if isTest {
		fmt.Println("This is a test")
	} else {
		fmt.Println("other usage")
	}
}

func list() {
	rows, err := Db.Query("SELECT * FROM student ORDER BY id DESC")
	if err != nil {

	}
}