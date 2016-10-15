package controllers

import (
    . "github.com/newbyebye/mediom/app/models"
    "github.com/revel/revel"
    "github.com/jinzhu/gorm"
    //"fmt"
)

type GameRestful struct {
    App
}

func (c GameRestful) Index() revel.Result {
    var id,skip, limit int
    
    c.Params.Bind(&skip, "skip")
    c.Params.Bind(&limit, "limit")
    c.Params.Bind(&id, "id")
    
    games := FindGame(id, skip, limit)
    
    return c.RenderJson(games)
}

func (c GameRestful) Show() revel.Result {
    g := Game{}
    DB.Preload("GameTemplate").First(&g, c.Params.Get("id"))

    user := c.CurrentUser()
    //fmt.Printf("##################### Show  %v %v\n", user.Id, g.UserID)
    if user.Id == g.UserID {
        g.ShowResult = true
    }

    return c.RenderJson(g)
}

func (c GameRestful) Win() revel.Result {
    var userGames []UserGame

    DB.Model(&UserGame{}).Preload("User", func(db *gorm.DB) *gorm.DB {
        return DB.Select(DefaultSelect())
    }).Preload("GameTemplate").Where("game_id = ? and is_win > 0", c.Params.Get("id")).Find(&userGames)

    return c.RenderJson(userGames)
}

func (c GameRestful) Stat() revel.Result {
    type Result struct{
        Var1           int
        Var2           int
        Count          int
    }

    var results []Result

    DB.Table("user_games").Select("var1, count(var1) as count, var2").Group("var1").Where("game_id = ?", c.Params.Get("id")).Scan(&results)

    return c.RenderJson(results)
}

func (c GameRestful) Template() revel.Result {    
    templates := FindGameTemplate()
    
    return c.RenderJson(templates) 
}

func (c GameRestful) TemplateShow() revel.Result{
    t := GameTemplate{}
    DB.First(&t, c.Params.Get("id"))

    var templates []GameTemplate

    DB.Model(&GameTemplate{}).Where("type = ?", t.Type).Order("id asc").Find(&templates)
    return c.RenderJson(templates) 
}
