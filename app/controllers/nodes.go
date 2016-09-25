package controllers

import (
    "strconv"
    //"strings"
    "fmt"
    . "github.com/newbyebye/mediom/app/models"
    "github.com/revel/revel"
)

type Nodes struct {
    App
}

func (c Nodes) Index() revel.Result {
    c.RenderArgs["nodes"] = FindAllNodes() 
    c.RenderArgs["groups"] = FindAllNodeRoots() 
    return c.Render()
}


func (c Nodes) Create() revel.Result {
    c.requireAdmin()

    parentId, _ := strconv.Atoi(c.Params.Get("parentId"))
    t := &Node{
        Name:  c.Params.Get("name"),
        ParentId: &parentId,
    }

    //t.UserId = c.currentUser.Id
    v := CreateNode(t)
    if v.HasErrors() {
        return c.renderValidation("nodes/index.html", v)
    }
    return c.Redirect("/nodes")
}

func (c Nodes) Edit() revel.Result {
    c.requireUser()
    t := &Node{}
    DB.Where("id = ?", c.Params.Get("id")).First(t)
    
    c.RenderArgs["node"] = t
    c.RenderArgs["groups"] = FindAllNodeRoots() 

    return c.Render()
}

func (c Nodes) Update() revel.Result {
    c.requireAdmin()
    t := Node{}
    DB.First(&t, c.Params.Get("id"))
    
    t.Name = c.Params.Get("name")
    t.Summary = c.Params.Get("summary")
    parentId, _ := strconv.Atoi(c.Params.Get("parentId"))
    t.ParentId = &parentId;

    v := UpdateNode(&t)
    if v.HasErrors() {
        return c.renderValidation("nodes/edit.html", v)
    }
    return c.Redirect(fmt.Sprintf("/nodes/%v/edit", t.Id))
}

func (c Nodes) Delete() revel.Result {
    c.requireAdmin()
    t := Node{}
    DB.First(&t, c.Params.Get("id"))
    if !c.isOwner(t) {
        c.Flash.Error("没有修改的权限")
        return c.Redirect("/")
    }

    err := DB.Delete(&t).Error
    if err != nil {
        c.RenderError(err)
    }
    return c.Redirect("/nodes")
}
