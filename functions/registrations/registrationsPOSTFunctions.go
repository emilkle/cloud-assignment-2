package registrations

import (
	"countries-dashboard-service/resources"
	"time"
)

func CreatePOSTResponse() (resources.RegistrationsPOSTResponse, error) {
	allDocuments, _ := GetAllRegisteredDocuments()
	nextId := len(allDocuments) + 1

	return resources.RegistrationsPOSTResponse{
		Id:         nextId,
		LastChange: time.Now().Format("20060102 15:04"),
	}, nil
}
