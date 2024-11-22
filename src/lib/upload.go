package lib

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ConfigType struct {
	AccessKeyId     string
	SecretAccessKey string
	DefaultRegion   string
	Bucket          string
	CDNUrl          string
}

var S3Client *s3.Client
var S3Config S3ConfigType
var mimeTypeToExt = map[string]string{
	"image/png":          "png",
	"image/jpeg":         "jpg",
	"application/pdf":    "pdf",
	"application/msword": "doc",
	"text/csv":           "csv",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
}

func hashFile(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Upload(file multipart.File, fileName string) (string, string) {

	fileContent, _ := io.ReadAll(file)

	hash, _ := hashFile(bytes.NewReader(fileContent))

	// Read the first 512 bytes of the file & Detect the file's content type
	buf := fileContent[:512]
	mimeType := http.DetectContentType(buf)
	ext := mimeTypeToExt[mimeType]
	if mimeType == "application/octet-stream" {
		splitName := strings.Split(fileName, ".")
		ext = splitName[len(splitName)-1]
	}

	filename := fmt.Sprintf("%s.%s", hash, ext)

	if _, err := S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(S3Config.Bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(fileContent),
	}); err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s/%s", S3Config.CDNUrl, filename), fileName
}

// Initializing
func InitS3Lib(configs S3ConfigType) {
	S3Config = configs

	S3Client = s3.New(s3.Options{
		Credentials:      credentials.NewStaticCredentialsProvider(S3Config.AccessKeyId, S3Config.SecretAccessKey, ""),
		Region:           S3Config.DefaultRegion,
		RetryMaxAttempts: 5,
		RetryMode:        aws.RetryModeStandard,
		HTTPClient:       &http.Client{Timeout: 30 * time.Second},
		ClientLogMode:    aws.LogRequestWithBody | aws.LogResponseWithBody,
	})
}
