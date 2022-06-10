package Models

import (
	"fmt"
	"reflect"
	"strings"

	"example.com/example/Utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exams struct {
	ID                       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseRef                primitive.ObjectID `json:"courseref"`
	ExamName                 string             `json:"examname"`
	NumberOfQuestionsContest int                `json:"numberofquestionscontest"`
	NumberOfQuestionsWritten int                `json:"numberofquestionswritten"`
}
type ExamsPopulated struct {
	ID                       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseRef                Courses            `json:"courseref"`
	ExamName                 string             `json:"examname"`
	NumberOfQuestionsContest int                `json:"numberofquestionscontest"`
	NumberOfQuestionsWritten int                `json:"numberofquestionswritten"`
}
type ExamsSearch struct {
	ID                             primitive.ObjectID `json:"id"`
	IDIsUsed                       bool               `json:"idisused"`
	CourseRef                      primitive.ObjectID `json:"courseref"`
	CourseRefIsUsed                bool               `json:"courserefisused"`
	ExamName                       string             `json:"examname"`
	ExamNameIsUsed                 bool               `json:"examnameisused"`
	NumberOfQuestionsContest       int                `json:"numberofquestionscontest"`
	NumberOfQuestionsContestIsUsed bool               `json:"numberofquestionscontestisused"`
	NumberOfQuestionsWritten       int                `json:"numberofquestionswritten"`
	NumberOfQuestionsWrittenIsUsed bool               `json:"numberofquestionswrittenisused"`
}

func (obj Exams) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.ExamName, validation.Required),
	)
}
func (obj Exams) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}
func (obj ExamsSearch) GetUsersSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}
	if obj.CourseRefIsUsed {
		self["courseref"] = obj.CourseRef
	}
	if obj.NumberOfQuestionsContestIsUsed {
		self["numberofquestionscontest"] = obj.NumberOfQuestionsContest
	}
	if obj.NumberOfQuestionsWrittenIsUsed {
		self["numberofquestionswritten"] = obj.NumberOfQuestionsWritten
	}

	if obj.ExamNameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.ExamName)
		self["examname"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	return self
}
func (obj *ExamsPopulated) CloneFrom(other Exams) {
	obj.ID = other.ID
	obj.CourseRef = Courses{}
	obj.ExamName = other.ExamName
	obj.NumberOfQuestionsContest = other.NumberOfQuestionsContest
	obj.NumberOfQuestionsWritten = other.NumberOfQuestionsWritten
}
