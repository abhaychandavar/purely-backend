package authService

import (
	"auth/internal/database"
	"auth/internal/database/models"
	firebaseHelper "auth/internal/utils/helpers/firebaseHelpers"
	httpErrors "auth/internal/utils/helpers/httpError"
	"context"
	"log"

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

func GetAuthToken(uid *string) (string, error) {
	auth := models.FindOne(context.Background(), database.Mongo().Db(), models.Auth{Identifier: *uid})
	if auth.Err() != nil {
		log.Default().Println(auth.Err().Error())
		return "", httpErrors.HydrateHttpError("purely/requests/get-auth-token/errors/invalid-user", 400, "Could not find user")
	}
	firebaseAuth, err := firebaseHelper.App().Auth(context.Background())
	if err != nil {
		log.Default().Println(err)
		return "", httpErrors.HydrateHttpError("purely/requests/get-auth-token/errors/internal_server_error", 500, "Internal Server Error")
	}
	var authRecord models.Auth

	if err := auth.Decode(&authRecord); err != nil {
		log.Default().Println(err)
		return "", httpErrors.HydrateHttpError("purely/requests/get-auth-token/errors/internal_server_error", 500, "Internal Server Error")
	}

	token, err := firebaseAuth.CustomTokenWithClaims(context.Background(), *uid, map[string]interface{}{"id": authRecord.ID.Hex()})
	if err != nil {
		log.Default().Println(err)
		return "", httpErrors.HydrateHttpError("purely/requests/get-auth-token/errors/internal_server_error", 500, "Internal Server Error")
	}

	return token, nil
}
