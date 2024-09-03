package api

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

//func GetImage(bucketName, objectName string) ([]byte, error) {
//	if GoogleStorageClient == nil {
//		return nil, fmt.Errorf("google storage client not initialized")
//	}
//
//	// Get a handle to the bucket
//	bucket := GoogleStorageClient.Bucket(bucketName)
//
//	// Get a handle to the object (image) within the bucket
//	obj := bucket.Object(objectName)
//
//	// Read the object's content
//	reader, err := obj.NewReader(GCloudContext)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create reader for object %s: %v", objectName, err)
//	}
//	defer func(reader *storage.Reader) {
//		err := reader.Close()
//		if err != nil {
//			log.Fatalf("failed to close reader for object %s: %v", objectName, err)
//		}
//	}(reader)
//
//	// Read all the data from the object
//	data, err := io.ReadAll(reader)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read object %s: %v", objectName, err)
//	}
//
//	return data, nil
//}

func AddImageToBucket(bucketName, objectName string, data []byte) (string, error) {
	if GoogleStorageClient == nil {
		return "", fmt.Errorf("google storage client not initialized")
	}

	bucket := GoogleStorageClient.Bucket(bucketName)
	obj := bucket.Object(objectName)
	w := obj.NewWriter(GCloudContext)

	defer func(w *storage.Writer) {
		err := w.Close()
		if err != nil {
			log.Printf("Error when closing the writer %v", err)
		}
	}(w)

	if _, err := w.Write(data); err != nil {
		return "", fmt.Errorf("failed to write object %s: %v", objectName, err)
	}
	var urlString = GetPublicURL(bucketName, objectName)
	log.Printf("wrote object %s with URL %s", objectName, urlString)
	return urlString, nil
}

// GetPublicURL generates the public URL for an object in a bucket
func GetPublicURL(bucketName, objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
}
