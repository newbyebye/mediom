package controllers

import (
	//"fmt"
	//"strconv"
	//"strings"
	"time"
	. "github.com/newbyebye/mediom/app/models"
	"github.com/revel/revel"
	"github.com/jinzhu/gorm"
)

type TopicsRestful struct {
	RestController
}

func (c TopicsRestful)Count(channel string) revel.Result{
	var m map[string]interface{}

	m = make(map[string]interface{})

	m["count"] = TopicsCountCached()

	return c.RenderJson(m)
}

func (c TopicsRestful) Index(channel string) revel.Result {
	var skip, limit int
	
	c.Params.Bind(&skip, "skip")
	c.Params.Bind(&limit, "limit")
	
	topics := FindTopic(c.CurrentUser(), channel, skip, limit)
	
	return c.RenderJson(topics)
}

func (c TopicsRestful) Show() revel.Result {
	t := Topic{}
	DB.Preload("User", func(db *gorm.DB) *gorm.DB {
        return db.Select(DefaultSelect())
	}).Preload("Node").First(&t, c.Params.Get("id"))

	if DB.Unscoped().Where("topic_id = ? AND date(start_time) = date(?) ", t.Id, time.Now().Format("2006-01-02")).First(&t.Lesson).RecordNotFound(){
		t.Lesson = nil
	}

	DB.Table("user_topics").Where("topic_id = ?", t.Id).Count(&t.RegisterCount)

	return c.RenderJson(t)
}