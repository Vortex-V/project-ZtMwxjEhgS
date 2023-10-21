package responses

import (
	"app/src/models"
	"fmt"
	"reflect"
)

// Response
// Содаём тип интерфейса и структуру, указатель на которую его реализует
type Response interface {
	implement(_ Response)
}
type response struct{}

func (r *response) implement(_ Response) {}

func MapTo(r Response, args ...interface{}) (map[string]interface{}, error) {
	var (
		data    = make(map[string]interface{})
		message string
	)

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]interface{}:
			for key, value := range v {
				data[key] = value
			}
		case string:
			message = v
		case models.Model:
			rValue := reflect.ValueOf(v).Elem()
			rValueType := rValue.Type()
			rResponse := reflect.ValueOf(r).Elem()
			rResponseType := rResponse.Type()

			// Проходит по полям модели и записывает их в структуру ответа, если это возможно
			// По сути отфильтровывает поля модели, которые не нужно отдавать в ответе
			for i := 0; i < rValue.NumField(); i++ {
				valueField := rValue.Field(i)
				responseField := rResponse.FieldByName(rValueType.Field(i).Name)
				if responseField.CanSet() {
					responseField.Set(valueField)
				}
			}

			for i := 0; i < rResponse.NumField(); i++ {
				responseField := rResponse.Field(i)
				if responseField.CanInterface() {
					data[rResponseType.Field(i).Name] = responseField.Interface()
				}
			}
		default:
			return nil, fmt.Errorf("can't map response")
		}
	}

	resultData := map[string]interface{}{
		"data": data,
	}
	if message != "" {
		resultData["message"] = message
	}

	return resultData, nil
}
