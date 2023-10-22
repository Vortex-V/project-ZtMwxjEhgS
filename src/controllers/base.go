package controllers

import (
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/models"
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
)

var (
	ErrorNotFound = errors.New("not Found")
)

type Controller struct {
	web.Controller
}

type DataMap map[string]interface{}

// Response Принимает map[string]interface{} data и int status
func (c *Controller) Response(data ...interface{}) {
	for _, arg := range data {
		switch v := arg.(type) {
		case DataMap, map[string]interface{}:
			c.Data["json"] = v
		case int:
			c.Ctx.Output.Status = v
		}
	}
	_ = c.ServeJSON()
}

func (c *Controller) ResponseMapTo(r responses.Response, data ...interface{}) {
	var (
		result  = make(DataMap)
		message string
		status  int
	)
	for _, arg := range data {
		switch v := arg.(type) {
		case DataMap:
			for key, value := range v {
				result[key] = value
			}
		case string:
			message = v
		case int:
			status = v
		case models.Model:
			result = responses.MapTo(r, v)
		}
	}

	responseData := DataMap{
		"data": result,
	}
	if message != "" {
		responseData["message"] = message
	}

	c.Response(responseData, status)
}

func (c *Controller) ResponseError(data interface{}, status int) {
	switch data.(type) {
	case DataMap, string:
		c.Response(DataMap{"error": data}, status)
	default:
		c.Response(status)
	}
}

func (c *Controller) Load(data requests.Request) bool {
	err := c.parseRequestBody(data)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return false
	}

	if validationErrors := validateRequest(data); len(validationErrors) > 0 {
		c.ResponseError(validationErrors, 400)
		return false
	}
	return true
}

func (c *Controller) parseRequestBody(data requests.Request) (err error) {
	err = json.Unmarshal(c.Ctx.Input.RequestBody, data)
	if err != nil {
		return err
	}

	return nil
}

func validateRequest(data requests.Request) DataMap {
	var (
		errors = make(DataMap)
		valid  = validation.Validation{}
	)
	result, err := valid.Valid(data)
	if err != nil { // Ошибки валидатора
		errors["error"] = err.Error()
	} else if !result { // Запрос невалиден
		for _, err := range valid.Errors {
			errors[err.Field] = err.Message
		}
	}
	return errors
}
