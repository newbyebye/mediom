package models

import (
    //"errors"
    //"fmt"
    //"github.com/revel/revel"
    //"github.com/revel/revel/cache"
    //"time"
    //"github.com/jinzhu/gorm"
)

type GameTemplate struct{
    BaseModel

    Name            string      `gorm:"not null"`
    Type            int
    Subtype         int
    Subname         string
    RuleLabel       string

    Var1Label       string
    Var1Help        string
    Var1Type        int
    Var1Range       string
    Var1Select      string

    Var2Label       string
    Var2Help        string
    Var2Type        int
    Var2Range       string
    Var2Select      string
}

type Game struct{
    BaseModel

    User            *User         `gorm:"not null"`
    UserID          int32

    Topic           *Topic         `gorm:"not null"`
    PostID          uint

    GameTemplate    *GameTemplate      `gorm:"not null"`
    GameTemplateID  uint

    Status         int
    Reward         int
    GameTime       int
    PlayerNum      int
    ShowResult     bool
}

type UserGame struct{
    BaseModel

    User            *User          `gorm:"not null"`
    UserID          uint

    Game            *Game          `gorm:"not null"`
    GameID          uint
    
    Var1           int
    Var2           int
    IsWin          int
}

func FindGame(id, skip, limit int) (games []Game) {
    var pageInfo = Pagination{}
    pageInfo.Query = db.Model(&Game{}).Preload("GameTemplate")

    pageInfo.Query = pageInfo.Query.Where("post_id = ? ", id).Order("id asc")
    
    pageInfo.Offset(skip, limit).Find(&games)
    return 
}

func FindGameTemplate() (templates []GameTemplate){
    db.Model(&GameTemplate{}).Group("type").Order("id asc").Find(&templates)
    return
}