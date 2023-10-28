package forms

// Form
// Содаём тип интерфейса и структуру, указатель на которую его реализует.
//
// Можно было бы сделать
//
/*type Form interface {
	*RentEndForm | и т.д.
}*/
// но мне удобней указывать имплементацию непосредсвтенно в типе структуры
type Form interface {
	implement(_ Form)
}
type form struct{}

func (r *form) implement(_ Form) {}
