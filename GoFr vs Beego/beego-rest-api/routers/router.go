package routers

import (
	"beego-rest-api/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/users", &controllers.UserController{}, "get:GetAll;post:Post")
	web.Router("/users/:id", &controllers.UserController{}, "get:Get;put:Put;delete:Delete")
}
