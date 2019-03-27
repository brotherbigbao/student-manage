package main

import (
	_ "database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

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
