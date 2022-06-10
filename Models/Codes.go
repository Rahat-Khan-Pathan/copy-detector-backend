package Models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Codes struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseRef    primitive.ObjectID `json:"courseref"`
	ExamRef      primitive.ObjectID `json:"examref"`
	QuestionNo   int                `json:"questionno"`
	QuestionType string             `json:"questiontype"`
	Email        string             `json:"email"`
	Code         string             `json:"code"`
	SubmitDate   primitive.DateTime `json:"submitdate"`
}
type CodesCopy struct {
	Email      string             `json:"email"`
	Percentage float64            `json:"percentage"`
	Code       string             `json:"code"`
	SubmitDate primitive.DateTime `json:"submitdate"`
}
type CodesSearch struct {
	CourseRef  primitive.ObjectID `json:"courseref"`
	ExamRef    primitive.ObjectID `json:"examref"`
	QuestionNo int                `json:"questionno"`
	Email      string             `json:"email"`
}

func (obj CodesSearch) GetUsersSearchBSONObj() bson.M {
	self := bson.M{}

	self["courseref"] = obj.CourseRef
	self["examref"] = obj.ExamRef
	self["questionno"] = obj.QuestionNo
	self["email"] = obj.Email

	return self
}
