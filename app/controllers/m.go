package controllers

import (
    "github.com/revel/revel"
)

type M struct {
    App
}

//func init() {
//  revel.InterceptMethod((*Home).Before, revel.BEFORE)
//  revel.InterceptMethod((*Home).After, revel.AFTER)
//}

func (c M) Index() revel.Result {
    user := c.CurrentUser()

    c.RenderArgs["user"] = user
    return c.Render()
}