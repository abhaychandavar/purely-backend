package services

import (
	"context"
	"fmt"
	"log"
	"media/internal/database"
	"media/internal/database/models"
	"media/internal/types/mediaServiceTypes"
	"media/internal/utils/constants"
	httpErrors "media/internal/utils/helpers/httpError"
	"media/internal/utils/helpers/httpHelper"
	mediahelpers "media/internal/utils/helpers/mediaHelpers"
	PubSub "media/providers/pubSub"
	"media/providers/storage"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaService struct {
	StorageProvider storage.StorageProvider
}

func (mediaService *MediaService) BlurImage(ctx context.Context, imageID string, profileID *string) (*string, error) {
	fmt.Println("BlurImage Image ID", imageID)
	fmt.Println("BlurImage Profile ID", profileID)
	imageIDPrimitive, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		return nil, err
	}
	imageMediaDataCur := models.FindOne(ctx, database.Mongo().Db(), models.Media{
		ID: imageIDPrimitive,
	})
	if imageMediaDataCur.Err() != nil {
		return nil, httpErrors.HydrateHttpError("purely/media/notFound", 404, "Not found")
	}
	var imageMediaData models.Media
	if err := imageMediaDataCur.Decode(&imageMediaData); err != nil {
		return nil, err
	}
	_, image, err := httpHelper.DownloadImageFromSignedURL(imageMediaData.URL)
	if err != nil {
		return nil, err
	}
	blurredImageBytes, _, err := mediahelpers.BlurImage(image, 20)
	if err != nil {
		return nil, err
	}

	rawFilePath := "blurred/" + imageMediaData.Path
	fileName := imageMediaData.FileName
	bucketName := "purely-public-assets"
	fileSize := len(blurredImageBytes)
	blurredImageType := "image/jpeg"
	filePathSplits := strings.Split(rawFilePath, imageMediaData.ContentType)
	filePath := filePathSplits[0] + blurredImageType + filePathSplits[1]

	initUploadRes, err := mediaService.StorageProvider.InitiateMultipartUpload(
		bucketName,
		filePath,
		fileName,
		blurredImageType,
		fileSize,
	)

	if err != nil {
		return nil, err
	}

	signedURLsRes, err := mediaService.StorageProvider.GenerateSignedURLsForParts(
		bucketName,
		filePath,
		fileName,
		initUploadRes.UploadId,
		blurredImageType,
		fileSize)
	if err != nil {
		return nil, err
	}
	uploadRes, err := mediaService.StorageProvider.UploadFile(
		signedURLsRes.SignedUrls,
		blurredImageBytes,
		len(blurredImageBytes),
		blurredImageType,
	)
	if err != nil {
		return nil, err
	}
	uploadCompleteRes, err := mediaService.StorageProvider.CompleteMultipartUpload(
		bucketName,
		initUploadRes.UploadId,
		filePath,
		fileName,
		blurredImageType,
		uploadRes,
	)
	if err != nil {
		return nil, err
	}
	fmt.Println("URL <><><><><>", uploadCompleteRes.URL)
	savedImage, err := models.Create(ctx, database.Mongo().Db(), models.Media{
		ID:          primitive.NewObjectID(),
		URL:         uploadCompleteRes.URL,
		EXT:         constants.FileExtMap[blurredImageType],
		Path:        uploadCompleteRes.Path,
		Domain:      uploadCompleteRes.Domain,
		ContentType: blurredImageType,
		FileName:    fileName,
		Size:        fileSize,
	})
	if err != nil {
		return nil, err
	}
	insertedID := savedImage.InsertedID.(primitive.ObjectID).Hex()
	if profileID != nil {
		NotifyImageBlurred(ctx, imageID, insertedID, *profileID)
	}
	return &insertedID, nil
}

func NotifyImageBlurred(ctx context.Context, mediaID string, blurredImageID string, profileID string) {
	pubsub := *PubSub.GetClient()
	pubsub.PublishToService(ctx, "profiles", PubSub.PublishMessageType{
		Type: "imageBlurred",
		Data: map[string]interface{}{
			"mediaID":        mediaID,
			"profileID":      profileID,
			"blurredImageID": blurredImageID,
		},
	})
}

func (profileService *MediaService) GenerateMediaUploadSignedUrl(ctx context.Context, mediaUploadData mediaServiceTypes.GenerateMediaUploadSignedUrlType) (*mediaServiceTypes.GenerateMediaUploadSignedUrlResType, error) {
	id := uuid.New()
	signedUrlData, error := profileService.StorageProvider.GenerateSignedUrl(
		"purely-profiles",
		fmt.Sprintf("profiles/%s/media/%s/%s/%s",
			mediaUploadData.AuthId,
			mediaUploadData.Purpose,
			mediaUploadData.ContentType,
			id.String()),
		mediaUploadData.FileName,
		mediaUploadData.ContentType,
		mediaUploadData.FileSize)
	if error != nil {
		log.Printf("Error generating signed URL: %v", error)
		return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/could-not-generate-signed-url", 500, "Failed to generate signed URL")
	}
	return &mediaServiceTypes.GenerateMediaUploadSignedUrlResType{
		SignedUrl: signedUrlData.SignedUrl,
		Expiry:    signedUrlData.Expires.Unix(),
	}, nil
}

func (profileService *MediaService) GenerateMultipartUploadUrls(mediaUploadData mediaServiceTypes.GenerateMultipartUploadUrlsType) (*mediaServiceTypes.GenerateMultipartUploadUrlsResType, error) {
	id := uuid.New()
	bucket := "purely-public-assets"
	filePath := fmt.Sprintf("profiles/%s/media/%s/%s/%s",
		mediaUploadData.AuthId,
		mediaUploadData.Purpose,
		mediaUploadData.ContentType,
		id.String())

	rawFileName := mediaUploadData.FileName
	fileNameSplit := strings.Split(rawFileName, ".")[0]
	fileName := fileNameSplit

	uploadData, err := profileService.StorageProvider.InitiateMultipartUpload(
		bucket,
		filePath,
		fileName,
		mediaUploadData.ContentType,
		mediaUploadData.FileSize,
	)
	if err != nil {
		return nil, err
	}

	res, err := profileService.StorageProvider.GenerateSignedURLsForParts(bucket, filePath, fileName, uploadData.UploadId, mediaUploadData.ContentType, int(mediaUploadData.FileSize))
	if err != nil {
		return nil, err
	}
	return &mediaServiceTypes.GenerateMultipartUploadUrlsResType{
		SignedUrls: res.SignedUrls,
		Expiry:     res.Expiry.Unix(),
		UploadID:   uploadData.UploadId,
		FilePath:   filePath,
		PartsCount: res.PartsCount,
		URL:        res.URL,
	}, nil
}

func (profileService *MediaService) CompleteMultipartUpload(ctx context.Context, mediaUploadData mediaServiceTypes.CompleteMultipartUploadType) (*mediaServiceTypes.CompleteMultipartUploadResType, error) {
	pathSplits := strings.Split(mediaUploadData.URL, "/")
	mimeType := pathSplits[len(pathSplits)-3]
	contentType := pathSplits[len(pathSplits)-4] + "/" + pathSplits[len(pathSplits)-3]
	filePath := strings.Join(pathSplits[:len(pathSplits)-1], "/")
	fileName := strings.Split(pathSplits[len(pathSplits)-1], ".")[0]

	res, err := profileService.StorageProvider.CompleteMultipartUpload("purely-public-assets", mediaUploadData.UploadID, filePath, fileName, contentType, mediaUploadData.Parts)
	if err != nil {
		return nil, err
	}

	media, err := models.Create(ctx, database.Mongo().Db(), models.Media{
		ID:          primitive.NewObjectID(),
		URL:         res.URL,
		EXT:         mimeType,
		ContentType: contentType,
		Path:        filePath,
		FileName:    fileName,
		Domain:      res.Domain,
		Size:        int(res.FileSize),
	})
	if err != nil {
		log.Printf("Error creating media entry: %v", err)
		return nil, err
	}
	return &mediaServiceTypes.CompleteMultipartUploadResType{
		URL: res.URL,
		ID:  media.InsertedID.(primitive.ObjectID).Hex(),
	}, nil
}

func (i *MediaService) HandlePubSubMessage(ctx context.Context, data PubSub.PublishMessageType) bool {
	fmt.Println("handlePubSubMessage data", data)
	fmt.Println("handlePubSubMessage data type", data.Type)
	switch data.Type {
	case "blurImage":
		{
			fmt.Println("handlePubSubMessage blurImage")
			refID := data.Data["refID"].(string)
			_, err := i.BlurImage(ctx, data.Data["imageID"].(string), &refID)
			if err != nil {
				log.Printf("Error handling pubsub message: %v", err)
				return false
			}
			return true
		}
	}
	return true
}
