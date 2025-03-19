package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	instance *AWSStorageProvider
	once     sync.Once
)

type AWSStorageProvider struct {
	clientInstance *s3.S3
}

func NewAWSStorageProvider(region string, accessKey string, secretKey string) (*AWSStorageProvider, error) {
	var err error
	once.Do(func() {
		sess, sessionErr := session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		})
		if sessionErr != nil {
			err = sessionErr
			return
		}

		instance = &AWSStorageProvider{
			clientInstance: s3.New(sess),
		}
	})

	return instance, err
}

type AWSInitiateMultipartUploadResult struct {
	UploadID string `xml:"UploadId"`
}

func (provider *AWSStorageProvider) GenerateSignedUrl(bucket string, filePath string, MimeType string, fileSize int64) (*UploadSignedUrl, error) {
	client := provider.clientInstance
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &filePath,
	})
	expiry := 10 * time.Minute
	urlStr, err := req.Presign(expiry)
	if err != nil {
		fmt.Println("Error generating presigned URL:", err)
		return nil, err
	}
	return &UploadSignedUrl{
		Bucket:    bucket,
		FilePath:  filePath,
		SignedUrl: urlStr,
		Expires:   time.Now().Add(expiry),
	}, nil
}

func (provider *AWSStorageProvider) InitiateMultipartUpload(bucket string, filePath string, MimeType string, fileSize int64) (*InitiateMultipartUpload, error) {
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filePath),
		ContentType: aws.String(MimeType),
	}
	client := provider.clientInstance
	resp, err := client.CreateMultipartUpload(input)
	if err != nil {
		return nil, err
	}
	return &InitiateMultipartUpload{
		Bucket:   bucket,
		FilePath: filePath,
		UploadId: *resp.UploadId,
	}, nil
}

func (provider *AWSStorageProvider) GenerateSignedURLsForParts(bucket string, filePath string, uploadID string, partsCount int) (*GenerateSignedURLsForPartsResType, error) {
	presignedURLs := make(map[int]string)
	client := provider.clientInstance
	expiry := 10 * time.Minute
	for partNumber := 1; partNumber <= partsCount; partNumber++ {
		input := &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(filePath),
			PartNumber: aws.Int64(int64(partNumber)),
			UploadId:   aws.String(uploadID),
		}

		presignedReq, _ := client.UploadPartRequest(input)

		url, err := presignedReq.Presign(expiry)

		if err != nil {
			return nil, err
		}

		presignedURLs[partNumber] = url
	}

	return &GenerateSignedURLsForPartsResType{
		SignedUrls: presignedURLs,
		Expiry:     time.Now().Add(expiry),
	}, nil
}

func (provider *AWSStorageProvider) CompleteMultipartUpload(bucket string, uploadID string, filePath string, parts map[int]string) (*string, error) {
	client := provider.clientInstance
	partList := []*s3.CompletedPart{}
	for key, value := range parts {
		partList = append(partList, &s3.CompletedPart{
			PartNumber: aws.Int64(int64(key)),
			ETag:       &value,
		})
	}
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		UploadId: aws.String(uploadID),
		Key:      &filePath,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: partList,
		},
	}
	_, err := client.CompleteMultipartUpload(input)
	if err != nil {
		return nil, err
	}
	objUrl := fmt.Sprintf("https://dl1b79m70nfwv.cloudfront.net/%s", filePath)
	return &objUrl, nil
}
