package controllers

import (
    . "github.com/newbyebye/mediom/app/models"
    "github.com/revel/revel"
)

type UsersRestful struct {
    App
    user User
}

func init() {
    revel.InterceptMethod((*UsersRestful).Before, revel.BEFORE)
    // revel.InterceptMethod((*Users).After, revel.AFTER)
}

func (c *UsersRestful) Before() revel.Result {
    var id int
    c.Params.Bind(&id, "id")
    var err error
    c.user, err = FindUserById(id)
    if err != nil {
        c.Finish(c.NotFound("页面不存在。"))
    }
    
    c.user.Password = ""
    
    return nil
}

func (c UsersRestful) Show() revel.Result {
    return c.RenderJson(c.user)
}


