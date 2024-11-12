package authService

import (
	"auth/internal/database"
	"auth/internal/database/models"
	httpErrors "auth/internal/utils/helpers/httpError"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertAuth(auth models.Auth) (string, error) {
	existingAuth := models.FindOne(context.Background(), database.Mongo().Db(), models.Auth{Phone: auth.Phone})

	if existingAuth.Err() == nil {
		return "", httpErrors.HydrateHttpError("purely/requests/errors/phone-already-registered", 400, "Phone number already registered")
	}

	data, err := models.Create(context.Background(), database.Mongo().Db(), auth)
	if err != nil {
		return "", err
	}
	return data.InsertedID.(primitive.ObjectID).Hex(), nil
}
