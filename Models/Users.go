package Models

import (
	"reflect"
	"strings"

	"example.com/example/Utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	Password     string             `json:"password"`
	PasswordHash string             `json:"passwordbase64"`
	Registered   bool               `json:"registered"`
	Admin        bool               `json:"admin"`
}
type UsersPasswordChange struct {
	CurrentPassword string `json:"currentpassword"`
	NewPassword     string `json:"newpassword"`
}
type UsersSearch struct {
	ID               primitive.ObjectID `json:"id"`
	IDIsUsed         bool               `json:"idisused"`
	Username         string             `json:"username"`
	UsernameIsUsed   bool               `json:"usernameisused"`
	Registered       bool               `json:"registered"`
	RegisteredIsUsed bool               `json:"registeredisused"`
}
type UsersLogin struct {
	EmailOrUsername string `json:"usernameoremail"`
	Password        string `json:"password"`
}

func (obj Users) ValidateSignUp() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Username, validation.Required),
		validation.Field(&obj.Email, validation.Required),
		validation.Field(&obj.Password, validation.Required),
	)
}
func (obj UsersLogin) ValidateLogin() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.EmailOrUsername, validation.Required),
	)
}
func (obj Users) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID", "Username", "PasswordHash"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}
func (obj UsersSearch) GetUsersSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.UsernameIsUsed {
		self["username"] = obj.Username
	}
	if obj.RegisteredIsUsed {
		self["registered"] = obj.Registered
	}

	return self
}
