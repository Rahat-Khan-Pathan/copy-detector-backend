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

type Courses struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseName string               `json:"coursename"`
	Exams      []primitive.ObjectID `json:"exams"`
}
type CoursesPopulated struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseName string             `json:"coursename"`
	Exams      []Exams            `json:"exams"`
}
type CoursesSearch struct {
	ID               primitive.ObjectID `json:"id"`
	IDIsUsed         bool               `json:"idisused"`
	CourseName       string             `json:"coursename"`
	CourseNameIsUsed bool               `json:"coursenameisused"`
}

func (obj Courses) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.CourseName, validation.Required),
	)
}
func (obj Courses) GetModifcationBSONObj() bson.M {
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
func (obj CoursesSearch) GetUsersSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.CourseNameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.CourseName)
		self["coursename"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	return self
}
func (obj *CoursesPopulated) CloneFrom(other Courses) {
	obj.ID = other.ID
	obj.CourseName = other.CourseName
	obj.Exams = []Exams{}
}
