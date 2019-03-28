package model

type Student struct {
	Id uint32
	No uint32
	Name string
	CScore float32
	MathScore float32
	EnglishScore float32
	TotalScore float32
	AverageScore float32
	Ranking uint32
	CreatedTime string
	UpdatedTime string
}