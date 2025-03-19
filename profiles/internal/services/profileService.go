package services

import (
	"context"
	"fmt"
	"log"
	"profiles/internal/database"
	"profiles/internal/database/models"
	profileLayoutTypes "profiles/internal/types/profileLayout"
	"profiles/internal/types/profileServiceTypes"
	httpErrors "profiles/internal/utils/helpers/httpError"
	"strings"

	"profiles/internal/providers/storage"

	"github.com/google/uuid"
	"github.com/mmcloughlin/geohash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileService struct {
	StorageProvider storage.StorageProvider
}

func (profileService *ProfileService) CreateProfile(ctx context.Context, data profileServiceTypes.CreateProfileType) (string, error) {
	geoHash := geohash.EncodeWithPrecision(*data.Lat, *data.Lng, 5)
	profile, err := models.Create(ctx, database.Mongo().Db(), models.Profile{
		Location: &models.Location{Type: "Point", Coordinates: []float64{*data.Lat, *data.Lng}},
		GeoHash:  geoHash,
		Status:   "active",
		AuthId:   *data.AuthId,
		Category: *data.Category,
	})
	if err != nil {
		return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-create-profile", 400, "Phone number already registered")
	}
	return profile.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (profileService *ProfileService) GetProfile(ctx context.Context, data profileServiceTypes.GetProfileType) (interface{}, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"category": *data.Category, "authId": *data.AuthId}}},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "media",
				"localField":   "media.mediaID",
				"foreignField": "_id",
				"as":           "mediaDetails",
			}},
		},
	}

	profile, err := models.Aggregate(ctx, database.Mongo().Db(), models.Profile{}, pipeline)
	if err != nil {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/profile-not-found", 404, "Profile not found")
	}
	var results []bson.M
	if err := profile.All(ctx, &results); err != nil {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/profile-not-found", 404, "Profile not found")
	}
	if len(results) == 0 {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/profile-not-found", 404, "Profile not found")
	}
	profileToReturn := results[0]
	mediaDetails, ok := profileToReturn["mediaDetails"].(primitive.A)
	if !ok {
		fmt.Println("mediaDetails is not of type []primitive.A")
		return profileToReturn, nil
	}
	// currentMedia := profileToReturn["media"].([]primitive.A)
	mediaArr := []primitive.M{}
	rawMediaList, ok := profileToReturn["media"].(primitive.A)
	if !ok {
		fmt.Println("mediaList is not of type []primitive.A")
		return profileToReturn, nil
	}
	rawMediaListMap := make(map[string]primitive.M)
	for _, rawMedia := range rawMediaList {
		if mediaMap, ok := rawMedia.(primitive.M); ok {
			rawMediaListMap[mediaMap["mediaID"].(primitive.ObjectID).Hex()] = mediaMap
		}
	}
	for _, p := range mediaDetails {
		if mediaMap, ok := p.(primitive.M); ok {
			mediaArr = append(mediaArr, primitive.M{
				"id":       rawMediaListMap[mediaMap["_id"].(primitive.ObjectID).Hex()]["_id"],
				"ext":      mediaMap["ext"],
				"order":    rawMediaListMap[mediaMap["_id"].(primitive.ObjectID).Hex()]["order"],
				"mediaURL": mediaMap["url"],
				"mediaID":  mediaMap["_id"],
			})
		}
	}
	profileToReturn["media"] = mediaArr
	profileToReturn["mediaDetails"] = nil
	return profileToReturn, nil
}

func (profileService *ProfileService) GetProfileLayout(ctx context.Context, data profileServiceTypes.GetProfileLayoutType) (interface{}, error) {
	return []profileLayoutTypes.LayoutElement{
		profileLayoutTypes.ElementGroup{
			Id:    "basicDetails",
			Label: "Let us know you a bit",
			Elements: []profileLayoutTypes.LayoutElement{
				profileLayoutTypes.InputElement{
					Element: profileLayoutTypes.Element{
						Id:       "name",
						Type:     profileLayoutTypes.InputElementType,
						Required: true,
						Label:    "Your name",
					},
					Placeholder: "Abhay Chandavar",
					InputType:   "text",
				},
				profileLayoutTypes.SearchableSelectElement{
					Element: profileLayoutTypes.Element{
						Id:       "gender",
						Type:     profileLayoutTypes.SearchableSelectElementType,
						Required: true,
						Label:    "I am",
					},
					Options: []profileLayoutTypes.SelectOption{
						{
							Label: "Man",
							Value: "man",
							Id:    "man",
						},
						{
							Label: "Woman",
							Value: "woman",
							Id:    "woman",
						},
						{
							Label: "More",
							Value: "female",
							Id:    "female",
						},
					},
					DefaultOptionIds: []string{"man", "woman"},
				},
				profileLayoutTypes.SearchableSelectElement{
					Element: profileLayoutTypes.Element{
						Id:       "lookingFor",
						Type:     profileLayoutTypes.SearchableSelectElementType,
						Required: true,
						Label:    "I am looking for",
					},
					Options: []profileLayoutTypes.SelectOption{
						{
							Label: "Man",
							Value: "man",
							Id:    "man",
						},
						{
							Label: "Woman",
							Value: "woman",
							Id:    "woman",
						},
						{
							Label: "More",
							Value: "female",
							Id:    "female",
						},
					},
					DefaultOptionIds: []string{"man", "woman"},
				},
				profileLayoutTypes.SelectElement{
					Element: profileLayoutTypes.Element{
						Id:       "hereFor",
						Type:     profileLayoutTypes.SelectElementType,
						Required: true,
						Label:    "I am here for",
					},
					Options: []profileLayoutTypes.SelectOption{
						{
							Label: "Relationship",
							Value: "relationship",
							Id:    "relationship",
						},
						{
							Label: "Don't know yet",
							Value: "dontKnowYet",
							Id:    "dontKnowYet",
						},
					},
					Placeholder:  "Select one",
					InitialValue: "relationship",
				},
			},
		},
		profileLayoutTypes.ElementGroup{
			Id:    "bioAndPrompts",
			Label: "Grab their attention!",
			Elements: []profileLayoutTypes.LayoutElement{
				profileLayoutTypes.InputElement{
					Element: profileLayoutTypes.Element{
						Id:       "bio",
						Type:     profileLayoutTypes.InputElementType,
						Required: true,
						Label:    "Bio",
					},
					Placeholder: "Write something catchy about yourself",
					InputType:   "text",
				},
				profileLayoutTypes.Prompt{
					Element: profileLayoutTypes.Element{
						Id:       "prompts",
						Type:     profileLayoutTypes.PromptElementType,
						Required: true,
						Label:    "Prompts",
					},
					PromptOptions: []string{"The key to my heart is", "My kind of date is", "Fun according to me", "Most spontaneous thing I've done", "Dating me is like", "An unpopular opinion of mine is"},
					Count:         3,
					UniquePrompts: true,
					InputElement: profileLayoutTypes.PromptInput{
						InputType:   "text",
						Placeholder: "Enter a prompt",
					},
				},
			},
		},
		profileLayoutTypes.Images{
			Element: profileLayoutTypes.Element{
				Id:    "images",
				Type:  profileLayoutTypes.ImageElementType,
				Label: "Your pics won't pick matches, but they'll keep it real after you vibe",
			},
			Count:         4,
			RequiredCount: 1,
		},
		profileLayoutTypes.ElementGroup{
			Id:    "location",
			Label: "One last thing!",
			Elements: []profileLayoutTypes.LayoutElement{
				profileLayoutTypes.Location{
					Element: profileLayoutTypes.Element{
						Id:       "location",
						Type:     profileLayoutTypes.LocationElementType,
						Required: true,
						Label:    "We will need your location to find the best dates for you",
					},
				},
				profileLayoutTypes.DistanceStepper{
					Element: profileLayoutTypes.Element{
						Id:       "distanceStepper",
						Type:     profileLayoutTypes.LocationStepperElementType,
						Required: true,
						Label:    "We will try to match you with people within this distance",
					},
					MinDistance: 20,
					MaxDistance: 100,
					Unit:        "km",
				},
			},
		},
	}, nil
}

func (profileService *ProfileService) computeProfileCompletionScore(profile *models.Profile) int {
	score := 0
	if profile.Name != "" {
		score++
	}
	if profile.Age > 0 {
		score++
	}
	if profile.Gender != primitive.NilObjectID {
		score++
	}
	if profile.Location != nil {
		score++
	}
	if len(profile.Prompts) > 0 {
		score++
	}
	if profile.LocationLabel != "" {
		score++
	}
	return score
}

func (profileService *ProfileService) UpsertDatingProfile(ctx context.Context, profile *profileServiceTypes.UpsertDatingProfileType) (string, error) {
	// Validate input
	if profile.AuthId == nil {
		return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-input", 400, "AuthId cannot be null")
	}

	filter := models.Profile{
		AuthId:   *profile.AuthId,
		Category: "date",
	}
	upsertData := models.Profile{
		AuthId:   *profile.AuthId,
		Category: "date",
	}
	if profile.Name != nil {
		upsertData.Name = *profile.Name
	}

	if profile.Bio != nil {
		upsertData.Bio = *profile.Bio
	}

	if profile.Age != nil {
		upsertData.Age = *profile.Age
	}
	if profile.HereFor != nil {
		upsertData.HereFor = *profile.HereFor
	}
	if profile.LookingFor != nil {
		upsertData.LookingFor = *profile.LookingFor
	}
	if profile.PreferredMatchDistance != nil {
		upsertData.PreferredMatchDistance = *profile.PreferredMatchDistance
	}
	if profile.LocationLabel != nil {
		upsertData.LocationLabel = *profile.LocationLabel
	}
	if profile.Gender != nil {
		genderId := *profile.Gender
		gender, err := primitive.ObjectIDFromHex(genderId)
		if err != nil {
			return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-gender-id", 400, "Invalid gender ID")
		}
		upsertData.Gender = gender
	}

	if profile.Prompts != nil {
		var prompts []models.PromptElementType
		for _, prompt := range *profile.Prompts {
			promptId, err := primitive.ObjectIDFromHex(prompt.PromptId)
			if err != nil {
				return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-prompt-id", 400, "Invalid prompt ID")
			}
			prompts = append(prompts, models.PromptElementType{
				Prompt: promptId,
				Answer: prompt.Answer,
			})
		}
		upsertData.Prompts = prompts
	}
	if profile.Media != nil {
		var mediaElements []models.MediaType
		for _, media := range *profile.Media {
			mediaID, err := primitive.ObjectIDFromHex(media.MediaID)
			if err != nil {
				return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-image-id", 400, "Invalid image ID")
			}
			mediaElements = append(mediaElements, models.MediaType{
				MediaID: mediaID,
				Order:   media.Order,
			})
		}
		upsertData.Media = mediaElements
	}
	if profile.Location != nil {
		upsertData.Location = &models.Location{
			Type:        "Point",
			Coordinates: []float64{profile.Location.Lat, profile.Location.Lng},
		}
	}

	upsertData.ProfileCompletionScore = profileService.computeProfileCompletionScore(&upsertData)

	upsertResult, err := models.Upsert(ctx, database.Mongo().Db(), filter, upsertData)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/duplicate-entry", 400, "Phone number already registered")
		}
		log.Printf("Error during upsert: %v", err)
		return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-create-profile", 500, "Failed to create or update profile")
	}
	if upsertResult.UpsertedID != nil {
		return upsertResult.UpsertedID.(primitive.ObjectID).Hex(), nil
	}
	existingProfile := models.Profile{}
	err = models.FindOne(ctx, database.Mongo().Db(), filter).Decode(&existingProfile)
	if err != nil {
		log.Printf("Error fetching existing profile: %v", err)
		return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/not-found", 404, "Profile not found after upsert")
	}
	return existingProfile.ID.Hex(), nil
}

func (profileService *ProfileService) GetPrompts(ctx context.Context, data profileServiceTypes.GetPromptsType) (*profileServiceTypes.GetPromptsResponse, error) {
	limit := 20
	prompts, err := models.Find(
		ctx,
		database.Mongo().Db(),
		models.Prompt{
			Category: *data.Category,
		},
		options.Find().SetSort(bson.D{{Key: "order", Value: 1}}).SetLimit(int64(limit)),
	)
	if err != nil {
		log.Printf("Error getting prompts: %v", err)
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-get-prompts", 500, "Failed to get prompts")
	}
	var promptsData []models.Prompt
	if err := prompts.All(ctx, &promptsData); err != nil {
		log.Printf("Error getting prompts 2: %v", err)
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-get-prompts", 500, "Failed to get prompts")
	}
	return &profileServiceTypes.GetPromptsResponse{
		Page:    data.Page,
		Limit:   &limit,
		Records: promptsData,
	}, nil
}

func (profileService *ProfileService) GetGenders(ctx context.Context, data profileServiceTypes.GetGendersType) (interface{}, error) {
	limit := 20

	genders, err := models.Find(
		ctx,
		database.Mongo().Db(),
		models.Gender{},
		options.Find().SetSort(bson.D{{Key: "order", Value: 1}}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-get-genders", 500, "Failed to get genders")
	}
	var gendersData []models.Gender
	if err := genders.All(ctx, &gendersData); err != nil {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-get-genders", 500, "Failed to get genders")
	}
	return &profileServiceTypes.GetGendersResponseType{
		Page:    data.Page,
		Limit:   &limit,
		Records: gendersData,
	}, nil
}

func (profileService *ProfileService) GetProfiles(ctx context.Context, data profileServiceTypes.GetProfilesType) ([]models.Profile, error) {
	limit := 20
	var profileData models.Profile

	// Fetch the self profile to get the location coordinates
	err := models.FindOne(ctx, database.Mongo().Db(), models.Profile{
		AuthId:   data.AuthId,
		Category: data.Category,
	}).Decode(&profileData)
	if err != nil {
		log.Printf("Error fetching self profile: %v", err)
		return nil, err
	}

	if profileData.Location == nil {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/no-location-provided", 400, "Location required to find relevant matches")
	}
	location := profileData.Location
	latLng := location.Coordinates

	radius := 10000 // in meters
	if profileData.PreferredMatchDistance > 0 {
		radius = profileData.PreferredMatchDistance * 1000
	}

	// Geospatial query to find profiles within the radius
	filter := bson.M{
		"authId": bson.M{
			"$ne": data.AuthId,
		},
		"location": bson.M{
			"$geoWithin": bson.M{
				"$centerSphere": []interface{}{
					latLng,
					radius / 6378100.0,
				},
			},
		},
		"category": data.Category,
	}

	cursor, err := database.Mongo().Db().Collection("profiles").Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		log.Printf("Error fetching profiles: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []models.Profile
	if err := cursor.All(ctx, &profiles); err != nil {
		log.Printf("Error decoding profiles: %v", err)
		return nil, err
	}

	log.Printf("Found %d profiles within 60 km radius", len(profiles))
	return profiles, nil
}

func (profileService *ProfileService) GenerateMediaUploadSignedUrl(ctx context.Context, mediaUploadData profileServiceTypes.GenerateMediaUploadSignedUrlType) (*profileServiceTypes.GenerateMediaUploadSignedUrlResType, error) {
	id := uuid.New()
	signedUrlData, error := profileService.StorageProvider.GenerateSignedUrl(
		"purely-profiles",
		fmt.Sprintf("profiles/%s/media/%s/%s/%s/%s",
			mediaUploadData.AuthId,
			mediaUploadData.Purpose,
			mediaUploadData.MimeType,
			id.String(),
			mediaUploadData.FileName),
		mediaUploadData.MimeType,
		mediaUploadData.FileSize)
	if error != nil {
		log.Printf("Error generating signed URL: %v", error)
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-generate-signed-url", 500, "Failed to generate signed URL")
	}
	return &profileServiceTypes.GenerateMediaUploadSignedUrlResType{
		SignedUrl: signedUrlData.SignedUrl,
		Expiry:    signedUrlData.Expires.Unix(),
	}, nil
}

func (profileService *ProfileService) GenerateMultipartUploadUrls(mediaUploadData profileServiceTypes.GenerateMultipartUploadUrlsType) (*profileServiceTypes.GenerateMultipartUploadUrlsResType, error) {
	id := uuid.New()
	bucket := "purely-public-assets"
	filePath := fmt.Sprintf("profiles/%s/media/%s/%s/%s/%s",
		mediaUploadData.AuthId,
		mediaUploadData.Purpose,
		mediaUploadData.MimeType,
		id.String(),
		mediaUploadData.FileName)
	uploadData, err := profileService.StorageProvider.InitiateMultipartUpload(
		bucket,
		filePath,
		mediaUploadData.MimeType,
		mediaUploadData.FileSize,
	)
	if err != nil {
		return nil, err
	}
	res, err := profileService.StorageProvider.GenerateSignedURLsForParts(bucket, filePath, uploadData.UploadId, int(mediaUploadData.PartsCount))
	if err != nil {
		return nil, err
	}
	return &profileServiceTypes.GenerateMultipartUploadUrlsResType{
		SignedUrls: res.SignedUrls,
		Expiry:     res.Expiry.Unix(),
		UploadID:   uploadData.UploadId,
		FilePath:   filePath,
	}, nil
}

func (profileService *ProfileService) CompleteMultipartUpload(ctx context.Context, mediaUploadData profileServiceTypes.CompleteMultipartUploadType) (*profileServiceTypes.CompleteMultipartUploadResType, error) {
	res, err := profileService.StorageProvider.CompleteMultipartUpload("purely-public-assets", mediaUploadData.UploadID, mediaUploadData.FilePath, mediaUploadData.Parts)
	if err != nil {
		return nil, err
	}
	pathSplits := strings.Split(mediaUploadData.FilePath, "/")
	mimeType := pathSplits[len(pathSplits)-3]
	media, err := models.Create(ctx, database.Mongo().Db(), models.Media{
		ID:  primitive.NewObjectID(),
		Url: *res,
		EXT: mimeType,
	})
	if err != nil {
		log.Printf("Error creating media entry: %v", err)
		return nil, err
	}
	return &profileServiceTypes.CompleteMultipartUploadResType{
		URL: *res,
		ID:  media.InsertedID.(primitive.ObjectID).Hex(),
	}, nil
}
