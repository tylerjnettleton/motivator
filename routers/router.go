package routers

import (
	"github.com/tylerjnettleton/motivator/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/login", &controllers.UserLoginController{})
    beego.Router("/register", &controllers.UserRegisterController{})
}
