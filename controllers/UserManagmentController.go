package controllers

import (
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/tylerjnettleton/motivator/models"
)

type UserRegistration struct {
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
	Email     string `form:"email"`
	Password  string `form:"password"`
}

type UserLogin struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

// Register
type UserRegisterController struct {
	beego.Controller
}

func (this *UserRegisterController) Get() {

	this.TplName = "register.tpl"
}

func (this *UserRegisterController) Post() {

	// This is the logger to the console of 'BeeGo'
	l := logs.GetLogger()

	// This is the 'UserRegistration' structure to parse the form body into
	userToRegister := UserRegistration{}
	if err := this.ParseForm(&userToRegister); err != nil {
		// If this does not parse correctly, send the user to an error page
		// Todo: Redirect the user to an actual error page
		this.Redirect("/error", 404)
		return
	}

	// We now have a valid User Registration object which contains all the information needed
	// to register a user.
	// We need to validate a few things
	// 1.) Email has not been taken

	// Let's search for the email address in the database
	// If we find an email already associated with an account
	// We need to notify the end-user!

	// Create an ORM instance
	o := orm.NewOrm()
	o.Using("default")

	// Query the users table and create a query seter
	// Set the query seter filter to user.email == the email provided by the form request

	user := new(models.User)
	err := o.QueryTable(&user).Filter("email", userToRegister.Email).One(user)
	if err != orm.ErrNoRows {
		// Present an error to the user that the mail provided has already been taken
		this.Data["email_error"] = "The email address has already been used."
		this.TplName = "register.tpl"
		return
	}

	// Now lets check the password requirements
	// If they are incorrect, we will notify the user
	sevenOrMore, number, upper, special := VerifyPassword(userToRegister.Password)
	if sevenOrMore != true {
		this.Data["password_error"] = "Your password needs to be at least 7 characters long"
		this.TplName = "register.tpl"
		return
	}

	if number != true {
		this.Data["password_error"] = "Your password needs to contain at least one number"
		this.TplName = "register.tpl"
		return
	}

	if upper != true {
		this.Data["password_error"] = "Your password needs to contain at least one upper character"
		this.TplName = "register.tpl"
		return
	}

	if special != true {
		this.Data["password_error"] = "Your password needs to contain at least one special character"
		this.TplName = "register.tpl"
		return
	}

	// Let's create a salt and hash the users password before we store it in the database
	salt, _ := randbytes()
	hashedPassword, _ := HashPassword([]byte(userToRegister.Password), salt)

	l.Print(string(salt))
	passwordModel := new(models.UserPassword)
	passwordModel.Password = hex.EncodeToString(hashedPassword)
	passwordModel.Salt = hex.EncodeToString(salt)

	profileModel := new(models.Profile)
	profileModel.Age = 21
	profileModel.HeightFeet = 5
	profileModel.HeightInches = 8
	profileModel.Gender = "male"

	paymentModel := new(models.UserPayment)
	paymentModel.Name = "CC"

	userModel := new(models.User)
	userModel.FirstName = userToRegister.FirstName
	userModel.LastName = userToRegister.LastName
	userModel.Email = userToRegister.Email
	userModel.Password = passwordModel
	userModel.Profile = profileModel
	userModel.Payment = paymentModel

	_, err = o.Insert(passwordModel)
	// Todo: Whatever failed, we need to delete the previous
	//  creations because of there relationship with each other
	if err != nil {
		l.Print(err)
		return
	}

	_, err = o.Insert(profileModel)
	// Todo: Whatever failed, we need to delete the previous
	//  creations because of there relationship with each other
	if err != nil {
		l.Print(err)
		return
	}

	_, err = o.Insert(paymentModel)
	// Todo: Whatever failed, we need to delete the previous
	//  creations because of there relationship with each other
	if err != nil {
		l.Print(err)
		return
	}

	_, err = o.Insert(userModel)
	// Todo: Whatever failed, we need to delete the previous
	//  creations because of there relationship with each other
	if err != nil {
		l.Print(err)
		return
	}

	this.Redirect("/", 302)
}

// Login
type UserLoginController struct {
	beego.Controller
}

func (this *UserLoginController) Get() {
	this.TplName = "login.tpl"
}

func (this *UserLoginController) Post() {

	l := logs.GetLogger()

	// This is the 'UserRegistration' structure to parse the form body into
	userToLogin := UserRegistration{}
	if err := this.ParseForm(&userToLogin); err != nil {
		// If this does not parse correctly, send the user to an error page
		// Todo: Redirect the user to an actual error page
		this.Redirect("/error", 404)
		return
	}

	// Create an ORM instance
	o := orm.NewOrm()
	o.Using("default")

	// Query the user
	user := new(models.User)
	err := o.QueryTable(&user).Filter("email", userToLogin.Email).One(user)
	if err == orm.ErrNoRows {
		// Present an error to the user that the mail provided has already been taken
		this.Data["email_error"] = "The email address provided was not found."
		this.TplName = "login.tpl"
		return
	}

	// Read the password relation in the database
	if user.Password != nil {
		err := o.Read(user.Password)
		if err != nil {
			// Failed to fetch users profile
			// Todo: Lets handle this error better!
			return
		}
	}

	// Attempt to authenticate the user
	decodedHashedPassword, _ := hex.DecodeString(user.Password.Password)
	decodedSalt, _ := hex.DecodeString(user.Password.Salt)

	authenticated, err := Authenticate([]byte(userToLogin.Password), decodedSalt, decodedHashedPassword)
	if authenticated == true {
		l.Print("Successfully authenticated user")
		// Return a valid token to the user and store it in database with a time stamp to expire
	} else {
		l.Print("Failed to authenticated user")
	}
}
