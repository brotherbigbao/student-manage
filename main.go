package main

import (
	"database/sql"
	_ "database/sql"
	"flag"
	"fmt"
	"gopkg.in/AlecAivazis/survey.v1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/modood/table"
	"log"
	"strconv"
	"student/model"
	"time"
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
	var qs = []*survey.Question{
		{
			Name: "no",
			Prompt: &survey.Input{Message: "请输入学号"},
			Validate: survey.Required,
		},
		{
			Name: "name",
			Prompt: &survey.Input{Message: "请输入姓名"},
			Validate: survey.Required,
		},
		{
			Name: "c_score",
			Prompt: &survey.Input{Message: "请输入C语言成绩"},
			Validate: survey.Required,
		},
		{
			Name: "math_score",
			Prompt: &survey.Input{Message: "请输入数学成绩"},
			Validate: survey.Required,
		},
		{
			Name: "english_score",
			Prompt: &survey.Input{Message: "请输入英语成绩"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		No uint32
		Name string
		CScore float32	`survey:"c_score"`
		MathScore float32	`survey:"math_score""`
		EnglishScore float32	`survey:"english_score"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	totalScore := answers.CScore + answers.MathScore + answers.EnglishScore
	averageScore := totalScore/3
	timeStr := time.Now().Format("2006-01-02 15:04:05")



	insertSql := "INSERT INTO student(no,name,c_score,math_score,english_score,total_score,average_score,ranking,updated_time,created_time) VALUES(?,?,?,?,?,?,?,?,?,?)"
	res, err := Db.Exec(insertSql, answers.No, answers.Name, answers.CScore, answers.MathScore, answers.EnglishScore, totalScore, averageScore, 0, timeStr, timeStr)
	if err != nil {
		panic(err)
	}

	affectedNum, err := res.RowsAffected()
	fmt.Println("成功创建" + strconv.Itoa(int(affectedNum)) + "条数据")

	updateRanking()
}

func updateRanking() {
	rows, err := Db.Query("SELECT id FROM student ORDER BY total_score DESC")
	if err != nil {
		panic(err)
	}

	increment := 1

	for rows.Next() {
		var studentId uint32
		err = rows.Scan(&studentId)
		if err != nil {
			panic(err)
		}
		Db.Exec("UPDATE student SET ranking=? WHERE id=?", increment, studentId)
		increment++
	}

	return
}

func studentUpdate() {
	updateRanking()
}

func studentReport() {
	//todo 报表信息
}