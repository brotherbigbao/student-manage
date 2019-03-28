package main

import (
	"database/sql"
	_ "database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/modood/table"
	"log"
	"student/model"
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
	var listFlag,addFlag,updateFlag bool
	flag.BoolVar(&listFlag, "l", false, "展示学生列表")
	flag.BoolVar(&addFlag, "a", false, "添加新学生信息")
	flag.BoolVar(&updateFlag, "u", false, "更新学生信息")
	flag.Parse()

	if listFlag {
		studentList := studentList()
		table.Output(studentList)
	} else if addFlag {
		fmt.Println("This is a init command")
	} else if updateFlag {
		fmt.Println("This is a test")
	} else {
		flag.Usage()
	}
}

func studentList() (userList []model.Student) {
	rows, err := Db.Query("SELECT * FROM student ORDER BY id DESC")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		student := model.Student{}
		err = rows.Scan(&student.Id, &student.No, &student.Name, &student.CScore, &student.MathScore,
			&student.EnglishScore, &student.TotalScore, &student.AverageScore, &student.Ranking,
			&student.UpdatedTime, &student.CreatedTime)
		if err != nil {
			panic(err)
		}
		userList = append(userList, student)
	}

	return
}