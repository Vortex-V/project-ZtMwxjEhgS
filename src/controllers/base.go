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

type controller struct {
	web.Controller
}

type dataMap map[string]interface{}

// Принимает map[string]interface{} data и int status
func (c *controller) response(data ...interface{}) {
	for _, arg := range data {
		switch v := arg.(type) {
		case dataMap, map[string]interface{}:
			c.Data["json"] = v
		case int:
			c.Ctx.Output.Status = v
		}
	}
	_ = c.ServeJSON()
}

func (c *controller) responseMapTo(r responses.Response, data ...interface{}) {
	var (
		result  = make(dataMap)
		message string
		status  int
	)
	for _, arg := range data {
		switch v := arg.(type) {
		case dataMap:
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

	responseData := map[string]interface{}{
		"data": result,
	}
	if message != "" {
		responseData["message"] = message
	}

	c.response(responseData, status)
}

func (c *controller) responseError(data interface{}, status int) {
	switch data.(type) {
	case dataMap, string:
		c.response(dataMap{"error": data}, status)
	default:
		c.response(status)
	}
}

func (c *controller) load(data requests.Request) bool {
	err := c.parseRequestBody(data)
	if err != nil {
		c.responseError(err, 500)
		return false
	}

	if validationErrors := validateRequest(data); len(validationErrors) > 0 {
		c.responseError(validationErrors, 400)
		return false
	}
	return true
}

func (c *controller) parseRequestBody(data requests.Request) (err error) {
	err = json.Unmarshal(c.Ctx.Input.RequestBody, data)
	if err != nil {
		return err
	}

	return nil
}

func validateRequest(data requests.Request) dataMap {
	var (
		errors = make(dataMap)
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
