package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	port                      = os.Getenv("PORT")
	mongoConnUrl              = os.Getenv("MONGO_CONN_URL")
	db                        = os.Getenv("MONGO_DB_DATABASE")
	internalAccessToken       = os.Getenv("INTERNAL_ACCESS_TOKEN")
	firebaseConfigPath        = os.Getenv("FIREBASE_CONFIG_PATH")
	env                       = os.Getenv("APP_ENV")
	googleMapsAPIKey          = os.Getenv("GOOGLE_MAPS_API_KEY")
	googleServiceJsonFilePath = os.Getenv("GOOGLE_SERVICE_JSON_FILE_PATH")
)

type configType struct {
	Port                      string
	MongoConnUrl              string
	Db                        string
	InternalAccessToken       string
	FirebaseConfigPath        string
	Env                       string
	GoogleMapsAPIKey          string
	GoogleServiceJsonFilePath string
}

func GetConfig() configType {
	obj := configType{
		MongoConnUrl:              mongoConnUrl,
		Db:                        db,
		Port:                      port,
		InternalAccessToken:       internalAccessToken,
		FirebaseConfigPath:        firebaseConfigPath,
		Env:                       env,
		GoogleMapsAPIKey:          googleMapsAPIKey,
		GoogleServiceJsonFilePath: googleServiceJsonFilePath,
	}
	if port == "" {
		obj.Port = "8080"
	}
	if googleServiceJsonFilePath == "" {
		obj.GoogleServiceJsonFilePath = "profiles/googleService.json"
	}
	return obj
}
