package controllers

import (
	"app/src/components/requests"
	"encoding/json"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
)

type controller struct {
	web.Controller
}

type dataMap map[string]interface{}

// Принимает dataMap data и int status
func (c *controller) response(data ...interface{}) {
	for _, arg := range data {
		switch v := arg.(type) {
		case dataMap:
			c.Data["json"] = v
		case int:
			c.Ctx.Output.Status = v
		}
	}
	_ = c.ServeJSON()
}

func (c *controller) responseError(err error, status int) {
	c.response(dataMap{"error": err.Error()}, status)
}

func (c *controller) parseRequestBody(data requests.Request) (err error) {
	err = json.Unmarshal(c.Ctx.Input.RequestBody, data)
	if err != nil {
		return err
	}

	return nil
}

func validateRequest(data requests.Request) (errors dataMap) {
	valid := validation.Validation{}
	result, err := valid.Valid(data)
	if err != nil { // Ошибки валидатора
		errors["error"] = err.Error()
	} else if !result { // Запрос невалиден
		for _, err := range valid.Errors {
			errors[err.Key] = err.Message
		}
	}
	return errors
}
