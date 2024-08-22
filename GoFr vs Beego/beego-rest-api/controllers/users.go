package controllers

import (
	"beego-rest-api/models"
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"io/ioutil"
	"strconv"
)

type UserController struct {
	web.Controller
}

func (c *UserController) GetAll() {
	o := orm.NewOrm()
	var users []*models.Users
	_, err := o.QueryTable("users").All(&users)
	if err != nil {
		c.Ctx.WriteString(err.Error())
	} else {
		c.Data["json"] = users
		c.ServeJSON()
	}
}

func (c *UserController) Get() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	o := orm.NewOrm()
	user := models.Users{Id: id}
	err := o.Read(&user)
	if errors.Is(err, orm.ErrNoRows) {
		c.Ctx.WriteString("User not found")
	} else if err != nil {
		c.Ctx.WriteString("Error reading data: " + err.Error())
	} else {
		c.Data["json"] = user
		c.ServeJSON()
	}
}

func (c *UserController) Post() {
	var user models.Users

	// Read the request body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.WriteString("Error reading request body: " + err.Error())
		return
	}

	// Unmarshal the request body into the user struct
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.Ctx.WriteString("Error unmarshalling data: " + err.Error())
		return
	}

	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		c.Ctx.WriteString("Error inserting data: " + err.Error())
	} else {
		c.Data["json"] = map[string]int64{"id": id}
		c.ServeJSON()
	}
}

func (c *UserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	var user models.Users

	// Read the request body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.WriteString("Error reading request body: " + err.Error())
		return
	}

	// Unmarshal the request body into the user struct
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.Ctx.WriteString("Error unmarshalling data: " + err.Error())
		return
	}

	user.Id = id
	o := orm.NewOrm()
	_, err = o.Update(&user)
	if err != nil {
		c.Ctx.WriteString("Error updating data: " + err.Error())
	} else {
		c.Data["json"] = user
		c.ServeJSON()
	}
}

func (c *UserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	o := orm.NewOrm()
	if num, err := o.Delete(&models.Users{Id: id}); err == nil {
		c.Data["json"] = map[string]int64{"num": num}
		c.ServeJSON()
	} else {
		c.Ctx.WriteString(err.Error())
	}
}
