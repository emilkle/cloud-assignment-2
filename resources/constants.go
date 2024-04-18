package resources

const RootPath = "http://localhost:8080"
const DefaultPort = "8080"
const RegistrationsPath = "/dashboard/v1/registrations/"
const DashboardsPath = "/dashboard/v1/dashboards/"
const NotificationsPath = "/dashboard/v1/notifications/"
const StatusPath = "/dashboard/v1/status/"
const RestCountriesPath = "http://129.241.150.113:8080/v3.1"
const OpenMeteoPath = "https://open-meteo.com/en/features#available-apis"
const CurrencyPath = "http://129.241.150.113:9090/currency/"

// Constants for fetching feature values
const GeocodingMeteo = "https://geocoding-api.open-meteo.com/v1"
const MeteoTempPercip = "https://api.open-meteo.com/v1"

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

// Basic error constants
const StandardDatatypeError = "Incorrect datatype: "
const DecodingError = "Error during JSON decoding "
const EncodingError = "Error during JSON encoding "

// Firestore collections
const RegistrationsCollection = "Registrations"
const WebhookCollection = "Webhooks"

// Example structs
const JsonStructPostAndPut = `{
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
