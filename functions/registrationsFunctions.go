package functions

import (
	"cloud.google.com/go/firestore"
	"countries-dashboard-service/database"
	"countries-dashboard-service/resources"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"strconv"
	"time"
)

func CreateRegistrationsGET(idParam string) (resources.RegistrationsGET, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	idNumber, err2 := strconv.Atoi(idParam)
	if err2 != nil {
		log.Println("This id could not be parsed, try another id.", err2.Error())
		return resources.RegistrationsGET{}, err2
	}
	query := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==", idNumber).Limit(1)
	documents, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Println("Failed to fetch documents:", err)
		return resources.RegistrationsGET{}, err
	}

	// Check if any documents were found
	if len(documents) == 0 {
		err4 := fmt.Errorf("No document found with ID: %s", idParam)
		log.Println(err4)
		return resources.RegistrationsGET{}, err4
	}

	for _, document := range documents {
		data := document.Data()
		featuresData := data["features"].(map[string]interface{})

		return resources.RegistrationsGET{
			Id:      idNumber,
			Country: data["country"].(string),
			IsoCode: data["isoCode"].(string),
			Features: resources.Features{
				Temperature:      featuresData["temperature"].(bool),
				Precipitation:    featuresData["precipitation"].(bool),
				Capital:          featuresData["capital"].(bool),
				Coordinates:      featuresData["coordinates"].(bool),
				Population:       featuresData["population"].(bool),
				Area:             featuresData["area"].(bool),
				TargetCurrencies: getTargetCurrencies(featuresData),
			},
			LastChange: data["lastChange"].(string),
		}, nil
	}
	log.Println("Document with ID", idParam, "was not found.")
	return resources.RegistrationsGET{}, nil
}

func getTargetCurrencies(featuresData map[string]interface{}) []string {
	targetCurrenciesInterface, ok := featuresData["targetCurrencies"].([]interface{})
	if !ok {
		// TargetCurrencies is not an array of interfaces
		log.Println("The requested currency data was not found.")
	}

	// Converting []interface{} to []string
	var targetCurrencies []string
	for _, currency := range targetCurrenciesInterface {
		if c, ok1 := currency.(string); ok1 {
			targetCurrencies = append(targetCurrencies, c)
		} else {
			// TargetCurrencies is not a string
			log.Println("The requested currency data is invalid.")
			return nil
		}
	}
	return targetCurrencies
}

func GetAllRegisteredDocuments() ([]resources.RegistrationsGET, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	iter := client.Collection(resources.REGISTRATIONS_COLLECTION).OrderBy("lastChange", firestore.Asc).Documents(ctx)

	var registrationsResponses []resources.RegistrationsGET

	idIndex := 1
	for {
		document, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := document.Data()

		// Retrieve the lastChange timestamp from the document
		lastChange, ok := data["lastChange"].(string)
		if !ok {
			log.Printf("The timestamp of the last change"+
				" %v could not be converted to string.", data["lastChange"])
			continue
		}

		registrationsResponse := createRegistrationsResponse(data, lastChange, idIndex)
		registrationsResponses = append(registrationsResponses, registrationsResponse)

		idIndex++
	}

	return registrationsResponses, nil
}

func createRegistrationsResponse(data map[string]interface{}, lastChange string, idIndex int) resources.RegistrationsGET {
	featuresData := data["features"].(map[string]interface{})

	return resources.RegistrationsGET{
		Id:      idIndex,
		Country: data["country"].(string),
		IsoCode: data["isoCode"].(string),
		Features: resources.Features{
			Temperature:      featuresData["temperature"].(bool),
			Precipitation:    featuresData["precipitation"].(bool),
			Capital:          featuresData["capital"].(bool),
			Coordinates:      featuresData["coordinates"].(bool),
			Population:       featuresData["population"].(bool),
			Area:             featuresData["area"].(bool),
			TargetCurrencies: getTargetCurrencies(featuresData),
		},
		LastChange: lastChange,
	}
}

func ParsePostResponse(client *firestore.Client) (resources.RegistrationsPOSTResponse, error) {
	ctx := database.GetFirestoreContext()

	// Creating a query to find the highest existing id field in the Registration collection
	query := client.Collection(resources.REGISTRATIONS_COLLECTION).OrderBy("id", firestore.Desc).Limit(1)
	documents, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Println("Failed to fetch documents:", err)
		return resources.RegistrationsPOSTResponse{}, err
	}

	var highestId int64
	idFound := false
	for _, document := range documents {
		data := document.Data()

		id, ok := data["id"].(int64)
		if id == 0 {
			continue
		}
		if !ok {
			err4 := fmt.Errorf("ID %v could not be used.", id)
			log.Println(err4)
		}
		if id > highestId && id != 0 {
			highestId = id
			idFound = true
		}
	}

	nextId := 1
	if idFound && highestId != 0 {
		nextId = int(highestId + 1)
	}

	return resources.RegistrationsPOSTResponse{
		Id:         nextId,
		LastChange: time.Now().Format("20060102 15:04"),
	}, nil
}
