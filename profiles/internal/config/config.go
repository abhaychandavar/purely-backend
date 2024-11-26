package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	port                = os.Getenv("PORT")
	mongoConnUrl        = os.Getenv("MONGO_CONN_URL")
	db                  = os.Getenv("MONGO_DB_DATABASE")
	internalAccessToken = os.Getenv("INTERNAL_ACCESS_TOKEN")
	firebaseConfigPath  = os.Getenv("FIREBASE_CONFIG_PATH")
	env                 = os.Getenv("APP_ENV")
	googleMapsAPIKey    = os.Getenv("GOOGLE_MAPS_API_KEY")
)

type configType struct {
	Port                string
	MongoConnUrl        string
	Db                  string
	InternalAccessToken string
	FirebaseConfigPath  string
	Env                 string
	GoogleMapsAPIKey    string
}

func GetConfig() configType {
	obj := configType{
		MongoConnUrl:        mongoConnUrl,
		Db:                  db,
		Port:                "8080",
		InternalAccessToken: internalAccessToken,
		FirebaseConfigPath:  firebaseConfigPath,
		Env:                 env,
		GoogleMapsAPIKey:    googleMapsAPIKey,
	}
	if port != "" {
		obj.Port = port
	}
	return obj
}
