package profileService

import (
	"auth/internal/database"
	"auth/internal/database/models"
	"auth/internal/types/profileServiceTypes"
	httpErrors "auth/internal/utils/helpers/httpError"
	"context"
	"log"

	"github.com/mmcloughlin/geohash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProfile(data profileServiceTypes.CreateProfileType) (string, error) {
	geoHash := geohash.EncodeWithPrecision(*data.Lat, *data.Lng, 5)
	profile, err := models.Create(context.Background(), database.Mongo().Db(), models.Profile{
		Location: models.Location{Type: "Point", Coordinates: []float64{*data.Lat, *data.Lng}},
		GeoHash:  geoHash,
		Status:   "active",
		AuthId:   *data.AuthId,
	})
	if err != nil {
		log.Fatal(err)
		return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-create-profile", 400, "Phone number already registered")
	}
	return profile.InsertedID.(primitive.ObjectID).Hex(), nil
}
