package main

import (
	"database/sql"
	_ "database/sql"
	"flag"
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
	var addFlag,addFromFileFlag,listFlag,listByNoFlag,listByNameFlag,updateFlag,deleteFlag,reportFlag,loginFlag,logoutFlag bool
	flag.BoolVar(&addFlag, "a", false, "从键盘添加新学生信息")
	flag.BoolVar(&addFromFileFlag, "af", false, "从文件添加新学生信息")

	flag.BoolVar(&listFlag, "l", false, "展示学生列表")
	flag.BoolVar(&listByNoFlag, "lno", false, "按学号查询")
	flag.BoolVar(&listByNameFlag, "lname", false, "按姓名查询")

	flag.BoolVar(&updateFlag, "u", false, "修改姓名")
	flag.BoolVar(&deleteFlag, "d", false, "删除记录")

	flag.BoolVar(&reportFlag, "r", false, "查看报表统计")

	flag.BoolVar(&loginFlag, "login", false, "登录")
	flag.BoolVar(&logoutFlag, "logout", false, "退出")

	flag.Parse()

	if addFlag {
		studentAdd()
	} else if addFromFileFlag {
		//todo
	} else if listFlag {
		studentList := studentList()
		table.Output(studentList)
	} else if listByNoFlag {
		//todo
	} else if listByNameFlag {
		//todo
	} else if updateFlag {
		studentUpdate()
	} else if deleteFlag {
		//todo
	} else if reportFlag {
		studentReport()
	} else if loginFlag {
		//todo
	} else if logoutFlag {
		//todo
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

func studentAdd() {
	//todo 学生信息新增
}

func studentUpdate() {
	//todo 学生信息编辑
}

func studentReport() {
	//todo 报表信息
}