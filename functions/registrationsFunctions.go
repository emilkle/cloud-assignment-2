package functions

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/resources"
	"log"
	"strconv"
)

func CreateRegistrationsGET(idParam string) (resources.RegistrationsGET, error) {
	// Loop through all entries in collection "Registrations"
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
		// Handle error
		log.Println("Failed to fetch documents:", err)
		return resources.RegistrationsGET{}, err
	}

	// Check if any documents were found
	if len(documents) == 0 {
		log.Printf("No document found with ID: %s", idParam)
		return resources.RegistrationsGET{}, nil
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
