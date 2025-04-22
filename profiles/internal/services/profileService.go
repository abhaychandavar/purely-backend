package services

import (
	"context"
	"fmt"
	"log"
	"profiles/internal/database"
	"profiles/internal/database/models"
	PubSub "profiles/internal/providers/pubSub"
	profileLayoutTypes "profiles/internal/types/profileLayout"
	"profiles/internal/types/profileServiceTypes"
	httpErrors "profiles/internal/utils/helpers/httpError"
	"unicode"

	"github.com/mmcloughlin/geohash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileService struct {
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

	existingProfile := models.Profile{}
	err := models.FindOne(ctx, database.Mongo().Db(), filter).Decode(&existingProfile)
	if err != nil {
		log.Printf("Error fetching existing profile: %v", err)
		return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/not-found", 404, "Profile not found after upsert")
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
	var mediaIDsToBlur []string
	if profile.Media != nil {
		var mediaElements []models.MediaType
		for _, media := range *profile.Media {
			mediaID, err := primitive.ObjectIDFromHex(media.MediaID)
			if err != nil {
				return "", httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-image-id", 400, "Invalid image ID")
			}
			mediaIDsToBlur = append(mediaIDsToBlur, mediaID.Hex())
			var blurredImageID *primitive.ObjectID
			blurredImageObjectID, err := primitive.ObjectIDFromHex(media.BlurredImageID)
			if err != nil {
				blurredImageID = &blurredImageObjectID
			}
			toAppendMediaEle := models.MediaType{
				MediaID: mediaID,
				Order:   media.Order,
			}
			if blurredImageID != nil {
				toAppendMediaEle.BlurredImageID = *blurredImageID
			}
			mediaElements = append(mediaElements, toAppendMediaEle)
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
	for _, mediaID := range mediaIDsToBlur {
		PubSub.GetClient().PublishToService(ctx, "media", PubSub.PubSubMessageType{
			Type: "blurImage",
			Data: map[string]interface{}{
				"mediaID":   mediaID,
				"profileID": existingProfile.ID.Hex(),
			},
		})
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

func (profileService *ProfileService) GetProfiles(ctx context.Context, data profileServiceTypes.GetProfilesType) ([]primitive.M, error) {
	// limit := 20
	var profileData models.Profile
	fmt.Println("Called Get Profiles")
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
	// latLng := profileData.Location.Coordinates

	// radius := 10000.0
	// if profileData.PreferredMatchDistance > 0 {
	// 	radius = float64(profileData.PreferredMatchDistance * 1000)
	// }

	// pipeline := mongo.Pipeline{
	// 	// 1. Geo filter
	// 	// {{Key: "$geoNear", Value: bson.M{
	// 	// 	"near":          bson.M{"type": "Point", "coordinates": latLng},
	// 	// 	"distanceField": "distance",
	// 	// 	"maxDistance":   radius,
	// 	// 	"spherical":     true,
	// 	// 	"query": bson.M{
	// 	// 		"authId":   bson.M{"$ne": data.AuthId},
	// 	// 		"category": data.Category,
	// 	// 		"status":   "active",
	// 	// 	},
	// 	// }}},
	// 	// 2. Unwind media
	// 	{{Key: "$unwind", Value: bson.M{
	// 		"path":                       "$media",
	// 		"preserveNullAndEmptyArrays": true,
	// 	}}},
	// 	// 3. Lookup media.mediaID
	// 	{{Key: "$lookup", Value: bson.M{
	// 		"from":         "media",
	// 		"localField":   "media.mediaID",
	// 		"foreignField": "_id",
	// 		"as":           "mediaData",
	// 	}}},
	// 	// 4. Lookup media.blurredImageID
	// 	{{Key: "$lookup", Value: bson.M{
	// 		"from":         "media",
	// 		"localField":   "media.blurredImageID",
	// 		"foreignField": "_id",
	// 		"as":           "blurredMediaData",
	// 	}}},
	// 	// 5. Merge looked-up fields
	// 	{{Key: "$addFields", Value: bson.M{
	// 		"media.mediaData":        bson.M{"$arrayElemAt": []interface{}{"$mediaData", 0}},
	// 		"media.blurredMediaData": bson.M{"$arrayElemAt": []interface{}{"$blurredMediaData", 0}},
	// 	}}},
	// 	// 6. Group back media array
	// 	{{Key: "$group", Value: bson.M{
	// 		"_id":     "$_id",
	// 		"profile": bson.M{"$first": "$$ROOT"},
	// 		"media":   bson.M{"$push": "$media"},
	// 	}}},
	// 	// 7. Merge grouped media back into profile
	// 	{{Key: "$addFields", Value: bson.M{
	// 		"profile.media": "$media",
	// 	}}},
	// 	{{Key: "$replaceRoot", Value: bson.M{
	// 		"newRoot": "$profile",
	// 	}}},
	// 	// 8. Limit
	// 	{{Key: "$limit", Value: limit}},
	// }
	fmt.Println("Get profiles called")
	pipeline := mongo.Pipeline{
		// Match by category
		{{Key: "$match", Value: bson.M{"category": data.Category}}},

		// Lookup mediaDetails for each media.mediaID
		{{Key: "$lookup", Value: bson.M{
			"from": "media",
			"let":  bson.M{"mediaArray": "$media"},
			"pipeline": mongo.Pipeline{
				{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$in": bson.A{"$_id", "$$mediaArray.mediaID"}},
				}}},
			},
			"as": "mediaDetails",
		}}},

		// Lookup blurredMediaDetails for each media.blurredImageID
		{{Key: "$lookup", Value: bson.M{
			"from": "media",
			"let": bson.M{"blurredIds": bson.M{
				"$map": bson.M{
					"input": "$media",
					"as":    "m",
					"in":    "$$m.blurredImageID",
				},
			}},
			"pipeline": mongo.Pipeline{
				{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$in": bson.A{"$_id", "$$blurredIds"}},
				}}},
			},
			"as": "blurredMediaDetails",
		}}},

		// Lookup promptDetails for each prompts.prompt
		{{Key: "$lookup", Value: bson.M{
			"from": "prompts",
			"let": bson.M{"promptIds": bson.M{
				"$map": bson.M{
					"input": "$prompts",
					"as":    "p",
					"in":    "$$p.prompt",
				},
			}},
			"pipeline": mongo.Pipeline{
				{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$in": bson.A{"$_id", "$$promptIds"}},
				}}},
			},
			"as": "promptDetails",
		}}},

		// Merge prompts with promptDetails, preserving other prompt fields like answer
		{{Key: "$addFields", Value: bson.M{
			"prompts": bson.M{
				"$map": bson.M{
					"input": "$prompts",
					"as":    "p",
					"in": bson.M{
						"$mergeObjects": bson.A{
							"$$p",
							bson.M{
								"prompt": bson.M{
									"$arrayElemAt": bson.A{
										bson.M{
											"$filter": bson.M{
												"input": "$promptDetails",
												"as":    "pd",
												"cond": bson.M{
													"$eq": bson.A{"$$pd._id", "$$p.prompt"},
												},
											},
										},
										0,
									},
								},
							},
						},
					},
				},
			},
		}}},

		// Merge mediaDetails and blurredMediaDetails for each media item
		{{Key: "$addFields", Value: bson.M{
			"mediaDetails": bson.M{
				"$map": bson.M{
					"input": "$media",
					"as":    "m",
					"in": bson.M{
						"media": bson.M{
							"$arrayElemAt": bson.A{
								bson.M{
									"$filter": bson.M{
										"input": "$mediaDetails",
										"as":    "md",
										"cond":  bson.M{"$eq": bson.A{"$$md._id", "$$m.mediaID"}},
									},
								},
								0,
							},
						},
						"blurredImage": bson.M{
							"$arrayElemAt": bson.A{
								bson.M{
									"$filter": bson.M{
										"input": "$blurredMediaDetails",
										"as":    "bd",
										"cond":  bson.M{"$eq": bson.A{"$$bd._id", "$$m.blurredImageID"}},
									},
								},
								0,
							},
						},
					},
				},
			},
		}}},

		// Final cleanup
		{{Key: "$project", Value: bson.M{
			"media":               0,
			"blurredMediaDetails": 0,
			"promptDetails":       0,
		}}},
	}

	cursor, err := models.Aggregate(ctx, database.Mongo().Db(), models.Profile{}, pipeline)
	if err != nil {
		log.Printf("Error aggregating profiles: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/profile-not-found", 404, "Profile not found")
	}
	if len(results) == 0 {
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/profile-not-found", 404, "Profile not found")
	}
	var profiles []primitive.M

	for _, profile := range results {
		runes := []rune(profile["name"].(string))
		firstNameChar := unicode.ToUpper(runes[0])
		profile["name"] = fmt.Sprintf("%s...", string(firstNameChar))
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (profileService *ProfileService) UpsertProfileBlurredImage(ctx context.Context, mediaID string, blurredImageID string, profileID string) {
	mediaObjectID, err := primitive.ObjectIDFromHex(mediaID)
	if err != nil {
		log.Printf("Invalid mediaObjectID: %v", err)
		return
	}
	profileObjectID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		log.Printf("Invalid blurredMediaID: %v", err)
		return
	}
	blurredImageObjectID, err := primitive.ObjectIDFromHex(blurredImageID)
	if err != nil {
		log.Printf("Invalid blurredImageID: %v", err)
		return
	}
	profile := models.FindOne(ctx, database.Mongo().Db(), models.Profile{
		ID: profileObjectID,
	})
	if profile.Err() != nil {
		log.Printf("Error fetching profile: %v", profile.Err())
		return
	}
	var profileData models.Profile
	if err := profile.Decode(&profileData); err != nil {
		log.Printf("Error decoding profile: %v", err)
		return
	}
	media := profileData.Media
	mediaArr := []models.MediaType{}
	for _, mediaEle := range media {
		currMediaEle := mediaEle
		if mediaEle.MediaID.Hex() == mediaObjectID.Hex() {
			currMediaEle.BlurredImageID = blurredImageObjectID
		}
		mediaArr = append(mediaArr, currMediaEle)
	}
	profileData.Media = mediaArr
	_, err = models.Upsert(ctx, database.Mongo().Db(), bson.M{"_id": profileObjectID}, profileData)
	if err != nil {
		log.Printf("Error updating profile: %v", err)
	}
}
