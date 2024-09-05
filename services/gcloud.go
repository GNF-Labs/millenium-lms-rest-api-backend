package services

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"path/filepath"
)

var GCloudContext = context.Background()

var GoogleStorageClient *storage.Client

func InitGoogleStorageClient() error {
	var err error

	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	credentialsFile := filepath.Join(projectRoot, "keys", "gnf-aic-6eca004e0090.json")
	client, err := storage.NewClient(GCloudContext, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return err
	}
	GoogleStorageClient = client
	return nil
}

func AddImageToBucket(bucketName, objectName string, data []byte) (string, error) {
	if GoogleStorageClient == nil {
		return "", fmt.Errorf("google storage client not initialized")
	}

	bucket := GoogleStorageClient.Bucket(bucketName)
	obj := bucket.Object(objectName)
	w := obj.NewWriter(GCloudContext)
	w.ContentType = "image/jpeg"
	w.CacheControl = "private"
	// Write the image data to the object
	if _, err := w.Write(data); err != nil {
		return "", fmt.Errorf("failed to write object %s: %v", objectName, err)
	}

	// Close the writer to finalize the upload
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer for object %s: %v", objectName, err)
	}

	// Get the public URL for the object
	urlString := GetPublicURL(bucketName, objectName)
	log.Printf("Wrote object %s with URL %s", objectName, urlString)
	return urlString, nil
}

// GetPublicURL generates the public URL for an object in a bucket
func GetPublicURL(bucketName, objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
}
