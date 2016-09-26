package controllers

import (
	. "github.com/newbyebye/mediom/app/models"
	"github.com/revel/revel"
	"regexp"
)

type Accounts struct {
	App
}

func init() {
	revel.InterceptMethod((*Accounts).Before, revel.BEFORE)
	// revel.InterceptMethod((*Accounts).After, revel.AFTER)
}

var (
	regexRequireUserActions, _ = regexp.Compile("Edit|Update|Password|UpdatePassword")
)

func (c *Accounts) Before() revel.Result {
	if regexRequireUserActions.MatchString(c.Action) {
		c.requireUser()
	}
	return nil
}

func (c Accounts) New() revel.Result {
	return c.Render()
}

func (c Accounts) Create() revel.Result {
	u := User{}
	newUser := User{}

	v := revel.Validation{}

	if !c.validateCaptcha(c.Params.Get("captcha")) {
		v.Error("验证码不正确")
		return c.renderValidation("accounts/new.html", v)
	}

	newUser, v = u.Signup(c.Params.Get("login"), c.Params.Get("email"), c.Params.Get("password"), c.Params.Get("password-confirm"))
	if v.HasErrors() {
		return c.renderValidation("accounts/new.html", v)
	}

	c.storeUser(&newUser)
	c.Flash.Success("注册成功")
	return c.Redirect(Home.Index)
}

func (c Accounts) Login() revel.Result {
	return c.Render()
}

func (c Accounts) LoginCreate() revel.Result {
	u := User{}
	newUser := User{}
	v := revel.Validation{}

	if !c.validateCaptcha(c.Params.Get("captcha")) {
		v.Error("验证码不正确")
		return c.renderValidation("accounts/login.html", v)
	}

	newUser, v = u.Signin(c.Params.Get("login"), c.Params.Get("password"))
	if v.HasErrors() {
		return c.renderValidation("accounts/login.html", v)
	}

	c.storeUser(&newUser)
	c.Flash.Success("登录成功，欢迎再次回来。")
	return c.Redirect(Home.Index)
}

func (c Accounts) Logout() revel.Result {
	c.clearUser()
	c.Flash.Success("登出成功")
	return c.Redirect(Home.Index)
}

func (c Accounts) Edit() revel.Result {
	return c.Render()
}

func (c Accounts) Update() revel.Result {
	c.currentUser.Email = c.Params.Get("email")
	c.currentUser.Fullname = c.Params.Get("fullname")
	c.currentUser.StudentNo = c.Params.Get("studentNo")
	c.currentUser.Profession = c.Params.Get("profession")
	c.currentUser.School = c.Params.Get("school")

	c.currentUser.Tagline = c.Params.Get("tagline")
	c.currentUser.Location = c.Params.Get("location")
	c.currentUser.Description = c.Params.Get("description")
	var u User
	u = *c.currentUser
	_, v := UpdateUserProfile(u)
	if v.HasErrors() {
		return c.renderValidation("accounts/edit.html", v)
	}
	c.Flash.Success("个人信息修改成功")
	return c.Redirect("/account/edit")
}

func (c Accounts) Password() revel.Result {
	return c.Render()
}

func (c Accounts) UpdatePassword() revel.Result {
	v := c.currentUser.UpdatePassword(c.Params.Get("password"), c.Params.Get("new-password"), c.Params.Get("confirm-password"))
	if v.HasErrors() {
		return c.renderValidation("accounts/password.html", v)
	}
	c.Flash.Success("密码修改成功")
	return c.Redirect("/account/password")
}
