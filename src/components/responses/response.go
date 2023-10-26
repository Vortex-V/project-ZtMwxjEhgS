package responses

import (
	"app/src/models"
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
	Construct(m models.Model) interface{}
}

// MapTo
// Deprecated
func MapTo(r Response, m models.Model) map[string]interface{} {
	rValue := reflect.ValueOf(m).Elem()
	rValueType := rValue.Type()
	rResponse := reflect.ValueOf(r).Elem()
	rResponseType := rResponse.Type()

	setFields(rValueType, rValue, rResponse)

	var data = make(map[string]interface{})
	for i := 0; i < rResponse.NumField(); i++ {
		responseField := rResponse.Field(i)
		if responseField.CanInterface() {
			data[rResponseType.Field(i).Name] = responseField.Interface()
		}
	}

	return data
}

func New[rT Response](r Response, m models.Model) rT {
	if rConstruct, ok := r.(ResponseConstructable); ok {
		r = rConstruct.Construct(m).(rT)
	} else {
		rtResponse := reflect.ValueOf(r).Elem()
		rtModel := reflect.ValueOf(m).Elem()
		rtModelT := rtModel.Type()
		setFields(rtModelT, rtModel, rtResponse)
	}

	return r.(rT)
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

func Collection[rT Response, mT models.Model](r rT, models []mT) []rT {
	var collection = make([]rT, 0, len(models))
	for _, v := range models {
		collection = append(collection, New[rT](r, v))
	}

	return collection
}
