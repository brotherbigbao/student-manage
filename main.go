package main

import (
	"database/sql"
	_ "database/sql"
	"flag"
	"fmt"
	"gopkg.in/AlecAivazis/survey.v1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/modood/table"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
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
	var addFlag,addFromFileFlag,listFlag,updateFlag,deleteFlag,reportFlag,loginFlag,logoutFlag bool
	var no int
	var name,filePath string

	flag.BoolVar(&addFlag, "a", false, "从键盘添加新学生信息")

	flag.BoolVar(&addFromFileFlag, "af", false, "从文件添加新学生信息")
	flag.StringVar(&filePath, "f", "", "文件路径")

	flag.BoolVar(&listFlag, "l", false, "展示学生列表")
	flag.IntVar(&no, "no", 0, "按学号查询")
	flag.StringVar(&name, "name", "", "按姓名查询")

	flag.BoolVar(&updateFlag, "u", false, "修改姓名")
	flag.BoolVar(&deleteFlag, "d", false, "删除记录")

	flag.BoolVar(&reportFlag, "r", false, "查看报表统计")

	flag.BoolVar(&loginFlag, "login", false, "登录")
	flag.BoolVar(&logoutFlag, "logout", false, "退出")

	flag.Parse()

	if addFlag {
		studentAdd()
	} else if addFromFileFlag {
		studentAddFromFile(filePath)
	} else if listFlag {
		studentList := studentList(no, name)
		table.Output(studentList)
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

func studentList(no int, name string) (userList []model.Student) {
	sql := "SELECT * FROM student ORDER BY id DESC"
	if no != 0 {
		sql = "SELECT * FROM student WHERE `no`=" + strconv.Itoa(no) + " ORDER BY id DESC"
	}
	if len(name) != 0 {
		sql = "SELECT * FROM student WHERE `name` LIKE '%" + name + "%' ORDER BY id DESC"
	}

	rows, err := Db.Query(sql)
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

func studentAddFromFile(filePath string) {
	if len(filePath) == 0 {
		fmt.Println("文件地址不能为空")
		return
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	dataStr := string(b[:])
	dataRows := strings.Split(dataStr, "\n")
	totalAffectedNum := 0
	for _, v := range dataRows {
		if len(v) < 1 {
			continue
		}
		fields := strings.Split(v, ",")
		if len(fields) != 5 {
			fmt.Println(v)
			fmt.Println("数据格式不正确")
		}

		cScore, err := strconv.Atoi(fields[2])
		if err != nil {
			panic(err)
		}
		mathScore, err := strconv.Atoi(fields[3])
		if err != nil {
			panic(err)
		}
		englishScore, err := strconv.Atoi(fields[4])
		if err != nil {
			panic(err)
		}
		totalScore := cScore + mathScore + englishScore
		averageScore := totalScore/3
		timeStr := time.Now().Format("2006-01-02 15:04:05")

		insertSql := "INSERT INTO student(no,name,c_score,math_score,english_score,total_score,average_score,ranking,updated_time,created_time) VALUES(?,?,?,?,?,?,?,?,?,?)"
		res, err := Db.Exec(insertSql, fields[0], fields[1], fields[2], fields[3], fields[4], totalScore, averageScore, 0, timeStr, timeStr)
		if err != nil {
			panic(err)
		}

		affectedNum, err := res.RowsAffected()
		totalAffectedNum += int(affectedNum)
	}
	fmt.Println("成功创建" + strconv.Itoa(int(totalAffectedNum)) + "条数据")
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



	updateSql := "UPDATE student SET name=?,c_score=?,math_score=?,english_score=?,total_score=?,average_score=?,ranking=?,updated_time=? WHERE no=?"
	res, err := Db.Exec(updateSql, answers.Name, answers.CScore, answers.MathScore, answers.EnglishScore, totalScore, averageScore, 0, timeStr, answers.No)
	if err != nil {
		panic(err)
	}

	affectedNum, err := res.RowsAffected()
	fmt.Println("成功更新" + strconv.Itoa(int(affectedNum)) + "条数据")

	updateRanking()
}

func studentReport() {
	//todo 报表信息
}