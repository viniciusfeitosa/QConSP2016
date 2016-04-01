package models

import (
	"encoding/json"
	"errors"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// UserModel is a struct
type UserModel struct {
	ID    bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string        `json:"name" bson:"name"`
	Email string        `json:"email" bson:"email"`
}

// GetUserModel is a method to return an instance of UserModel
// with all fields not empty
func (u *UserModel) GetUserModel() (*UserModel, error) {
	u.ID = bson.NewObjectId()
	if strings.TrimSpace(string(u.ID)) == "" {
		return nil, errors.New("[UserModel GetUserModel] User ID is can't be empty")
	}
	if strings.TrimSpace(u.Name) == "" {
		return nil, errors.New("[UserModel GetUserModel] User Name is can't be empty")
	}
	if strings.TrimSpace(u.Email) == "" {
		return nil, errors.New("[UserModel GetUserModel] User Email is can't be empty")
	}
	return u, nil
}

// ToJSONString generate a json string from an instance of UserModel
func (u *UserModel) ToJSONString() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSONToObject generate a json string from an instance of UserModel
func (u *UserModel) FromJSONToObject(value string) (UserModel, error) {
	var user UserModel
	if err := json.Unmarshal([]byte(value), &user); err != nil {
		return user, err
	}
	return user, nil
}
