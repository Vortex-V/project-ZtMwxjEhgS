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
	ErrorNotFound   = errors.New("not Found")
	ErrorBadRequest = errors.New("bad Request")
)

type Controller struct {
	web.Controller
}

type DataMap map[string]interface{}

// TODO отрефакторить методы response, а то расплодились, надо быть проще и меньше экспериментов.

// ResponseJson Принимает DataMap и int status
func (c *Controller) ResponseJson(data ...interface{}) {
	for _, arg := range data {
		switch v := arg.(type) {
		case DataMap:
			c.Data["json"] = v
		case int:
			c.Ctx.Output.Status = v
		}
	}
	_ = c.ServeJSON()
}

// Response Формирует DataMap из переданного responses.Response
func (c *Controller) Response(args ...interface{}) {
	var (
		responseData = make(DataMap)
		status       int
	)
	for _, arg := range args {
		switch v := arg.(type) {
		case DataMap:
			for key, value := range v {
				responseData[key] = value
			}
		case responses.Response:
			data, err := toMap(v)
			if err != nil {
				responseData["error"] = err.Error()
				continue
			}
			responseData["data"] = data
		case []responses.Response:
			data, err := toCollection(v)
			if err != nil {
				responseData["error"] = err.Error()
				continue
			}
			responseData["data"] = data
		case int:
			status = v
		default:

		}
	}

	c.ResponseJson(responseData, status)
}

// ResponseMapTo Фильтрует models.Model по переданному responses.Response и формирует DataMap
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

	c.ResponseJson(responseData, status)
}

func (c *Controller) ResponseError(data interface{}, status int) {
	switch data.(type) {
	case DataMap, string:
		c.ResponseJson(DataMap{"error": data}, status)
	default:
		c.ResponseJson(status)
	}
}

func (c *Controller) LoadAndValidate(data requests.Request) bool {
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

func toCollection(data interface{}) ([]DataMap, error) {
	list := make([]DataMap, 0)
	tmpStr, _ := json.Marshal(data)
	err := json.Unmarshal(tmpStr, &list)
	return list, err
}

func toMap(data interface{}) (DataMap, error) {
	result := make(DataMap)
	tmpStr, _ := json.Marshal(data)
	err := json.Unmarshal(tmpStr, &result)
	return result, err
}
