package responses

import (
	"app/src/models"
	"reflect"
)

// Response
// Содаём тип интерфейса и структуру, указатель на которую его реализует
type Response interface {
	implement(_ Response)
}
type response struct{}

func (r *response) implement(_ Response) {}

func MapTo(r Response, m models.Model) map[string]interface{} {
	rValue := reflect.ValueOf(m).Elem()
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

	var data = make(map[string]interface{})
	for i := 0; i < rResponse.NumField(); i++ {
		responseField := rResponse.Field(i)
		if responseField.CanInterface() {
			data[rResponseType.Field(i).Name] = responseField.Interface()
		}
	}

	return data
}
