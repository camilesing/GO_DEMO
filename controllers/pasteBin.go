package controllers

import (
	"encoding/json"
	"errors"
	"go_demo/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// PasteBinController operations for PasteBin
type PasteBinController struct {
	beego.Controller
}

// URLMapping ...
func (c *PasteBinController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create PasteBin
// @Param	body		body 	models.PasteBin	true		"body for PasteBin content"
// @Success 201 {int} models.PasteBin
// @Failure 403 body is empty
// @router / [post]
func (c *PasteBinController) Post() {
	var v models.PasteBin
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddPasteBin(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	httpPostForm(v.Poster, v.Syntax, v.Content)
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get PasteBin by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.PasteBin
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PasteBinController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPasteBinById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get PasteBin
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.PasteBin
// @Failure 403
// @router / [get]
func (c *PasteBinController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllPasteBin(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the PasteBin
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.PasteBin	true		"body for PasteBin content"
// @Success 200 {object} models.PasteBin
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PasteBinController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.PasteBin{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdatePasteBinById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the PasteBin
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PasteBinController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeletePasteBin(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func httpPostForm(poster string, syntax string, cotent string) (string, error) {
	urlStr := "http://pastebin.ubuntu.com/"
	resp, err := http.PostForm(urlStr,
		url.Values{"poster": {poster}, "syntax": {syntax}, "content": {cotent}})

	if err != nil {
		errorStr := "error happend when post form, error info: "
		return errorStr, errors.New(errorStr)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorStr := "error happend when read HTTP Request body , error info: "
		return errorStr, errors.New(errorStr)
	}

	str_body := (string(body))
	preStr := "pturl\" href=\""
	nextStr := "/plain/\">Download as text"
	index1 := strings.Index(str_body, preStr)
	index2 := strings.Index(str_body, nextStr)
	preStrSize := len(preStr)
	if index1 < 0 || index2 < 0 {
		errorStr := "please set correctly syntax type "
		return errorStr, errors.New(errorStr)
	}
	str := str_body[index1+preStrSize : index2]
	resultStr := "http://pastebin.ubuntu.com" + str
	log.Println(resultStr)
	return resultStr, nil

}
