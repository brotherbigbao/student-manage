package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/modood/table"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
	"log"
	"os"
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

	flag.BoolVar(&addFromFileFlag, "af", false, "从文件添加新学生信息 eg: ./student -af -f 文件路径")
	flag.StringVar(&filePath, "f", "", "文件路径, 和-af一起使用")

	flag.BoolVar(&listFlag, "l", false, "展示学生列表 eg: ./student -l")
	flag.IntVar(&no, "no", 0, "按学号查询 eg: ./student -l -no 1234")
	flag.StringVar(&name, "name", "", "按姓名查询 eg: ./student -l -name 朱志庭")

	flag.BoolVar(&updateFlag, "u", false, "修改学生信息（使用学号）eg: ./student -u -no 1234")
	flag.BoolVar(&deleteFlag, "d", false, "删除学生信息（使用学号）eg: ./student -d -no 1234")

	flag.BoolVar(&reportFlag, "r", false, "查看报表统计")

	flag.BoolVar(&loginFlag, "login", false, "登录")
	flag.BoolVar(&logoutFlag, "logout", false, "退出")

	flag.Parse()

	if !loginFlag && !logoutFlag {
		if !checkLogin() {
			fmt.Println("请先登录！")
			flag.Usage()
			return
		}
	}

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
		studentDelete(no)
	} else if reportFlag {
		studentReport()
	} else if loginFlag {
		login()
	} else if logoutFlag {
		logout()
	} else {
		flag.Usage()
	}
}

//展示学生列表
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

//从命令行添加学生
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

//从文件添加学生,一定要这样每个学生占一行,列用逗号分隔,依次是: "学号,姓名,c语言成绩,数学成绩,英语成绩", eg:83423,朱志庭,88,77,66
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

	updateRanking()
}

//更新学生排名 内部调用
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

//命令行手动更新学生信息
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

//删除学生 输入学生学号删除
func studentDelete(no int)  {
	if no == 0 {
		fmt.Println("请输入正确的学号")
		return
	}

	deleteSql := "DELETE FROM student WHERE no=?"
	res, err := Db.Exec(deleteSql, no)
	if err != nil {
		panic(err)
	}

	affectedNum, err := res.RowsAffected()
	fmt.Println("成功删除" + strconv.Itoa(int(affectedNum)) + "条数据")

	updateRanking()
}

//查看学生报表
func studentReport() {
	type StatResult struct {
		MaxCScore float32
		MaxMathScore float32
		MaxEnglishScore float32
		FailStudentNum int
	}

	maxScoreSql := "SELECT MAX(c_score) AS max_c_score, MAX(math_score) AS max_math_score, MAX(english_score) AS max_english_score FROM student"
	failStudents := "SELECT COUNT(*) AS num FROM student WHERE c_score < 60 OR math_score < 60 OR english_score < 60"

	maxScoreRow := Db.QueryRow(maxScoreSql)
	failStudentRow := Db.QueryRow(failStudents)

	stat := StatResult{}
	maxScoreRow.Scan(&stat.MaxCScore, &stat.MaxMathScore, &stat.MaxEnglishScore)
	failStudentRow.Scan(&stat.FailStudentNum)

	result := []StatResult{
		stat,
	}
	table.Output(result)
}

//登入 当前目录要有写权限
func login() {
	var qs = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{Message: "请输入姓名"},
			Validate: survey.Required,
		},
		{
			Name: "password",
			Prompt: &survey.Input{Message: "请输入密码"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name string
		Password string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	name := answers.Name
	password := answers.Password

	if name == "admin" && password == "123456" {
		err := ioutil.WriteFile("login.data", []byte("1"), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("登录成功！")
	} else {
		fmt.Println("用户名或密码错误！")
	}
}

//退出 当前目前要有写权限
func logout() {
	err := os.Remove("login.data")
	if err != nil {
		fmt.Println("退出失败，请重试！")
	}
	fmt.Println("退出成功！")
}

//判断登录状态 内部调用
func checkLogin() bool  {
	ioutil.ReadFile("login.data")
	b, err := ioutil.ReadFile("login.data")
	if err != nil {
		return false
	}

	if string(b[:]) == "1" {
		return true
	}

	return false
}