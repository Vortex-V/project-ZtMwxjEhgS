package requests

import (
	"github.com/beego/beego/v2/core/validation"
)

// Request
// Содаём тип интерфейса и структуру, указатель на которую его реализует
type Request interface {
	implement(_ Request)
}
type request struct{}

func (r *request) implement(_ Request) {}

func init() {
	validation.SetDefaultMessage(map[string]string{
		"Required":     "Должно быть заполнено",
		"Min":          "Минимально допустимое значение %d",
		"Max":          "Максимально допустимое значение %d",
		"Range":        "Должно быть в диапазоне от %d до %d",
		"MinSize":      "Минимально допустимая длина %d",
		"MaxSize":      "Максимально допустимая длина %d",
		"Length":       "Длина должна быть равна %d",
		"Alpha":        "Должно состоять из букв",
		"Numeric":      "Должно состоять из цифр",
		"AlphaNumeric": "Должно состоять из букв или цифр",
		"Match":        "Должно совпадать с %s",
		"NoMatch":      "Не должно совпадать с %s",
		"AlphaDash":    "Должно состоять из букв, цифр или символов (-_)",
		"Email":        "Должно быть в правильном формате email",
		"IP":           "Должен быть правильный IP адрес",
		"Base64":       "Должно быть представлено в правильном формате base64",
		"Mobile":       "Должно быть правильным номером мобильного телефона",
		"Tel":          "Должно быть правильным номером телефона",
		"Phone":        "Должно быть правильным номером телефона или мобильного телефона",
		"ZipCode":      "Должно быть правильным почтовым индексом",
	})
}