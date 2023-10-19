package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Rental struct {
	Id          int         `orm:"column(id);pk"`
	UserId      *Accounts   `orm:"column(user_id);rel(fk)"`
	TypeId      *RentTypes  `orm:"column(type_id);rel(fk)"`
	TransportId *Transports `orm:"column(transport_id);rel(fk)"`
	TimeStart   time.Time   `orm:"column(time_start);type(timestamp without time zone);null;auto_now_add"`
	TimeEnd     time.Time   `orm:"column(time_end);type(timestamp without time zone);null"`
	PriceOfUnit float64     `orm:"column(price_of_unit)"`
	FinalPrice  float64     `orm:"column(final_price);null"`
	Status      int         `orm:"column(status)"`
	CreatedAt   time.Time   `orm:"column(created_at);type(timestamp without time zone);auto_now_add"`
	UpdatedAt   time.Time   `orm:"column(updated_at);type(timestamp without time zone);auto_now_add"`
}

func (t *Rental) TableName() string {
	return "rental"
}

func init() {
	orm.RegisterModel(new(Rental))
}

// AddRental insert a new Rental into database and returns
// last inserted Id on success.
func AddRental(m *Rental) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRentalById retrieves Rental by Id. Returns error if
// Id doesn't exist
func GetRentalById(id int) (v *Rental, err error) {
	o := orm.NewOrm()
	v = &Rental{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRental retrieves all Rental matches certain condition. Returns empty list if
// no records exist
func GetAllRental(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rental))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Rental
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateRental updates Rental by Id and returns error if
// the record to be updated doesn't exist
func UpdateRentalById(m *Rental) (err error) {
	o := orm.NewOrm()
	v := Rental{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRental deletes Rental by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRental(id int) (err error) {
	o := orm.NewOrm()
	v := Rental{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Rental{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
