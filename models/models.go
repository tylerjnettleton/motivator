package models

import "github.com/astaxie/beego/orm"

type User struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Password  *UserPassword `orm:"rel(one)"`
	Profile   *Profile      `orm:"rel(one)"`
	Payment   *UserPayment  `orm:"rel(one)"`
}

type UserPassword struct {
	Id       int
	Password string
	Salt     string
}

type Profile struct {
	Id           int
	Age          int16
	HeightFeet   int16
	HeightInches int16
	Gender       string
}

type UserPayment struct {
	Id   int
	Name string
}

// Todo:
// Connected services

func init() {
	// Need to register model in init
	orm.RegisterModel(new(User), new(Profile), new(UserPayment), new(UserPassword))
}
