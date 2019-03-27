package model

type Student struct {
	Id uint8
	No int8
	Name string
	CScore float32
	MathScore float32
	EnglishScore float32
	TotalScore float32
	AverageScore float32
	Ranking int8
	CreatedTime string
	UpdatedTime string
}