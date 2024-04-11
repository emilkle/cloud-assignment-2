package resources

const DEFAULT_PORT = "8080"
const REGISTRATIONS_PATH = "/dashboard/v1/registrations/"
const DASHBOARDS_PATH = "/dashboard/v1/dashboards/"
const NOTIFICATIONS_PATH = "/dashboard/v1/notifications/"
const STATUS_PATH = "/dashboard/v1/status/"
const REST_COUNTRIES_PATH = "http://129.241.150.113:8080/v3.1"
const OPEN_METEO_PATH = "https://open-meteo.com/en/features#available-apis"
const CURRENCY_PATH = "http://129.241.150.113:9090/currency/"

// Constants for fetching feature values
const GEOCODING_METEO = "https://geocoding-api.open-meteo.com/v1"
const METEO_TEMP_PERCIP = "https://api.open-meteo.com/v1"

// Webhook invocation during development
const TEMP_WEBHOOK_INV = "https://webhook.site/268f5454-08b6-4639-84ee-99381ad547d2"

// Basic error constants
const STANDARD_DATATYPE_ERROR = "Incorrect datatype: "
const DECODING_ERROR = "Error during JSON decoding "
const ENCODING_ERROR = "Error during JSON encoding "

// Firestore collections
const REGISTRATIONS_COLLECTION = "Registrations"
const WEBHOOK_COLLECTION = "Webhooks"

// Example structs
const JSON_STRUCT_POST_AND_PUT = `{
	   "country": "Norway",                                     
	   "isoCode": "NO",                                        
	   "features": {
	                  "temperature": true,                      
	                  "precipitation": true,                    
	                  "capital": false,                          
	                  "coordinates": true,                      
	                  "population": false,                       
	                  "area": true,                             
	                  "targetCurrencies": ["EUR", "USD", "SEK"] 
	               }
}`
