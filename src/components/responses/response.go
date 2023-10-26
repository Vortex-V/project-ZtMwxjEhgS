package responses

import (
	"app/src/models"
	"encoding/json"
	"reflect"
)

// Response
// Содаём тип интерфейса и структуру, указатель на которую его реализует.
//
// Можно было бы сделать
//
/*type Response interface {
	*AccountMeResponse | *AccountSignUpResponse | и т.д.
}*/
// но мне удобней указывать имплементацию непосредсвтенно в типе структуры
type Response interface {
	implement(_ Response)
}
type response struct{}

func (r *response) implement(_ Response) {}

type ResponseConstructable interface {
	Response
	CustomFields(m models.Model) interface{}
}

// MapTo
// Deprecated
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

func New[rT ResponseConstructable](r rT, m models.Model, onlyCustomFields bool) rT {
	rtResponse := reflect.ValueOf(r).Elem()

	if !onlyCustomFields {
		rtModel := reflect.ValueOf(m).Elem()
		rtModelT := rtModel.Type()
		setFields(rtModelT, rtModel, rtResponse)
	}

	customFields := r.CustomFields(m)
	rtCustomFields := reflect.ValueOf(customFields)
	rtCustomFieldsT := rtCustomFields.Type()
	setFields(rtCustomFieldsT, rtCustomFields, rtResponse)

	return r
}

func setFields(rtSrcType reflect.Type, rtSrc, rtDst reflect.Value) {
	for i := 0; i < rtSrc.NumField(); i++ {
		rtSrcFieldValue := rtSrc.Field(i)
		rtDstField := rtDst.FieldByName(rtSrcType.Field(i).Name)
		if rtDstField.CanSet() {
			rtDstField.Set(rtSrcFieldValue)
		}
	}
}

func Collection[rT ResponseConstructable, mT models.Model](r rT, models []mT, onlyCustomFields bool) []rT {
	var collection = make([]rT, 0, len(models))
	for _, v := range models {
		newR := clone[rT](r)
		collection = append(collection, New[rT](newR, v, onlyCustomFields))
	}

	return collection
}

func clone[T any](el T) (newEl T) {
	tmpStr, _ := json.Marshal(el)
	_ = json.Unmarshal(tmpStr, &newEl)
	return
}
