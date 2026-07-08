package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func readBlob(problemID, fileName string) (string, error) {
	account := os.Getenv("AZURE_STORAGE_ACCOUNT")
	key := os.Getenv("AZURE_STORAGE_KEY")
	container := os.Getenv("AZURE_CONTAINER_NAME")
	folder := os.Getenv("FOLDER")

	blobName := fmt.Sprintf("%s/%s/%s", folder, problemID, fileName)

	cred, err := azblob.NewSharedKeyCredential(account, key)
	if err != nil {
		return "", err
	}

	client, err := azblob.NewClientWithSharedKeyCredential(
		fmt.Sprintf("https://%s.blob.core.windows.net/", account),
		cred,
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, err := client.DownloadStream(context.Background(), container, blobName, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}