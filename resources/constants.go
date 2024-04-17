package resources

const RootPath = "http://localhost:8080"
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

// Webhook event constants
const (
	POSTTitle   = "New country data is registered to dashboard"
	PUTTitle    = "New updates in dashboard"
	GETTitle    = "Invoked country data from dashboard"
	DELETETitle = "Deleted country data from dashboard"

	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventDelete   = "DELETE"
	EventInvoke   = "INVOKE"
)

// Webhook invocation during development
const WebhookInv = "/invocation/"

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
