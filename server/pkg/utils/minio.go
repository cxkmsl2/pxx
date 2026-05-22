package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var minioBucket = "pxx"

func InitMinio(endpoint, accessKey, secretKey string) error {
	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	return err
}

func GetMinioPresignedURL(objectName string, expiry time.Duration) (string, error) {
	if minioClient == nil {
		return "", fmt.Errorf("minio not initialized")
	}
	u, err := minioClient.PresignedPutObject(context.Background(), minioBucket, objectName, expiry)
	if err != nil { return "", err }
	return u.String(), nil
}

func GetMinioPublicURL(objectName string) string {
	return fmt.Sprintf("http://localhost:9000/%s/%s", minioBucket, objectName)
}
