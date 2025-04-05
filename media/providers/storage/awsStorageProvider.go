package storage

import (
	"bytes"
	"fmt"
	"math"
	"media/internal/utils/constants"
	"net/http"
	"strings"
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

func (provider *AWSStorageProvider) GenerateSignedUrl(bucket string, filePath string, fileName string, contentType string, fileSize int) (*UploadSignedUrl, error) {
	ext := constants.FileExtMap[contentType]
	formattedFilePath := filePath + "/" + fileName + "." + ext
	client := provider.clientInstance
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &formattedFilePath,
	})
	expiry := 10 * time.Minute
	urlStr, err := req.Presign(expiry)
	if err != nil {
		return nil, err
	}
	return &UploadSignedUrl{
		Bucket:    bucket,
		FilePath:  formattedFilePath,
		SignedUrl: urlStr,
		Expires:   time.Now().Add(expiry),
	}, nil
}

func (provider *AWSStorageProvider) InitiateMultipartUpload(bucket string, filePath string, fileName string, contentType string, fileSize int) (*InitiateMultipartUpload, error) {
	ext := constants.FileExtMap[contentType]
	formattedFilePath := filePath + "/" + fileName + "." + ext
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(formattedFilePath),
		ContentType: aws.String(contentType),
	}
	client := provider.clientInstance
	resp, err := client.CreateMultipartUpload(input)
	if err != nil {
		return nil, err
	}
	return &InitiateMultipartUpload{
		Bucket:   bucket,
		FilePath: formattedFilePath,
		UploadId: *resp.UploadId,
	}, nil
}

func (provider *AWSStorageProvider) GenerateSignedURLsForParts(bucket string, filePath string, fileName string, uploadID string, contentType string, fileSize int) (*GenerateSignedURLsForPartsResType, error) {
	presignedURLs := make(map[int]string)
	client := provider.clientInstance
	expiry := 10 * time.Minute
	partSize := 5 * 1024 * 1024
	partsCount := int(math.Ceil(float64(fileSize) / float64(partSize)))
	ext := constants.FileExtMap[contentType]
	formattedFilePath := filePath + "/" + fileName + "." + ext
	for partNumber := 1; partNumber <= partsCount; partNumber++ {
		input := &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(formattedFilePath),
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
		PartsCount: partsCount,
		URL:        formattedFilePath,
	}, nil
}

func (provider *AWSStorageProvider) CompleteMultipartUpload(bucket string, uploadID string, filePath string, fileName string, contentType string, parts map[int]string) (*CompletedMultipartUploadResponseType, error) {
	client := provider.clientInstance
	partList := []*s3.CompletedPart{}
	ext := constants.FileExtMap[contentType]
	formattedFilePath := filePath + "/" + fileName + "." + ext

	for key, value := range parts {
		partList = append(partList, &s3.CompletedPart{
			PartNumber: aws.Int64(int64(key)),
			ETag:       aws.String(value),
		})
	}
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		UploadId: aws.String(uploadID),
		Key:      aws.String(formattedFilePath),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: partList,
		},
	}
	_, err := client.CompleteMultipartUpload(input)
	if err != nil {
		return nil, err
	}
	awsBaseURL := "https://dl1b79m70nfwv.cloudfront.net"
	objUrl := fmt.Sprintf("%s/%s", awsBaseURL, formattedFilePath)

	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    &formattedFilePath,
	}
	headObjectOutput, err := client.HeadObject(headObjectInput)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve object metadata: %v", err)
	}

	fileSize := *headObjectOutput.ContentLength

	res := CompletedMultipartUploadResponseType{
		URL:      objUrl,
		Path:     filePath,
		Domain:   awsBaseURL,
		FileSize: fileSize,
	}
	return &res, nil
}

func uploadFile(signedURL string, data []byte, contentType string) (*string, error) {
	req, err := http.NewRequest("PUT", signedURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status code %d", resp.StatusCode)
	}

	etag := resp.Header.Get("ETag")
	if etag == "" {
		return nil, fmt.Errorf("ETag not found in response headers")
	}

	etag = strings.Replace(etag, "\"", "", -1)

	return &etag, nil
}

func (provider *AWSStorageProvider) UploadFile(signedUrls map[int]string, file []byte, partSize int, contentType string) (map[int]string, error) {
	var partPromises []chan string
	var wg sync.WaitGroup
	for partIndex, signedURL := range signedUrls {
		wg.Add(1)
		start := (partIndex - 1) * partSize
		end := int(math.Min(float64(start+partSize), float64(len(file))))

		filePart := file[start:end]

		errChan := make(chan string, 1)
		partPromises = append(partPromises, errChan)

		go func(signedURL string, filePart []byte, errChan chan string) {
			defer wg.Done()
			etagPtr, err := uploadFile(signedURL, filePart, contentType)
			if err != nil {
				errChan <- ""
			} else {
				errChan <- *etagPtr
			}
		}(signedURL, filePart, errChan)
	}

	wg.Wait()

	var result map[int]string = map[int]string{}
	for idx, resChan := range partPromises {
		etag := <-resChan
		result[idx+1] = etag
	}

	return result, nil
}
