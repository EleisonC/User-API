package models

import (

)


type DogOwner struct {
	FirstName string `bson:"firstname,omitempty" validate:"required"`
	LastName string `bson:"lastname,omitempty" validate:"required"`
	TelNumber string `bson:"telNumber,omitempty" validate:"required"`
	Email string `bson:"email,omitempty" validate:"required,email"`
	Password  string `bson:"password,omitempty" validate:"required"`
}




