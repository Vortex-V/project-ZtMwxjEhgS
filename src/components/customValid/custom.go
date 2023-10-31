package customValid

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

func UsernameUniqueRequest(v *validation.Validation, username string) {
	query := models.Find(&models.Account{Username: username}, "id").Where("username = ?")
	result, err := models.Raw(query, username).Exec()
	if err != nil {
		return
	}
	count, _ := result.RowsAffected()
	if count > 0 {
		v.SetError("Username", "Указанное имя пользователя уже занято")
	}
}

func TransportTypeExists(v *validation.Validation, transportType string) {
	if transportType != "All" {
		transportType = models.GetTransportType(transportType)
		if transportType == "" {
			v.SetError("transportType", "transportType должен быть одним из [Car, Bike, Scooter, All]")
			return
		}
	}
}

func RentTypeExists(v *validation.Validation, rentType string) {
	if rentType != "All" {
		rentType = models.GetRentType(rentType)
		if rentType == "" {
			v.SetError("rentType", "rentType должен быть одним из [Days, Minutes]")
			return
		}
	}
}
