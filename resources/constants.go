package resources

const DEFAULT_PORT = "8080"
const REGISTRATIONS_PATH = "/dashboard/v1/registrations/"
const DASHBOARDS_PATH = "/dashboard/v1/dashboards/"
const NOTIFICATIONS_PATH = "/dashboard/v1/notifications/"
const STATUS_PATH = "/dashboard/v1/status/"
const REST_COUNTRIES_PATH = "http://129.241.150.113:8080/v3.1"
const OPEN_METEO_PATH = "https://open-meteo.com/en/features#available-apis"
const CURRENCY_PATH = "http://129.241.150.113:9090/currency/"

// Basic error constants
const STANDARD_ERROR = "The request failed with error: "
const DECODING_ERROR = "Error during JSON decoding "
const ENCODING_ERROR = "Error during JSON encoding "

// Firestore collections
const REGISTRATIONS_COLLECTION = "Registrations"

// Example structs
const JSON_STRUCT_POST_AND_PUT = `{
	   "country": "Norway",                                     
	   "isoCode": "NO",                                        
	   "features": {
	                  "temperature": true,                      
	                  "precipitation": true,                    
	                  "capital": true,                          
	                  "coordinates": true,                      
	                  "population": true,                       
	                  "area": true,                             
	                  "targetCurrencies": ["EUR", "USD", "SEK"] 
	               }
}`
