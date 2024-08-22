package models

import "github.com/beego/beego/v2/client/orm"

type Users struct {
	Id    int    `orm:"auto" json:"id"`
	Name  string `orm:"size(100)" json:"name"`
	Email string `orm:"size(100)" json:"email"`
}

func init() {
	orm.RegisterModel(new(Users))
}
