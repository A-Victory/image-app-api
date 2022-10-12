// Package utils provides the connection to the cloudinary server.
package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func connect() *cloudinary.Cloudinary {
	cld, err := cloudinary.New()
	// For connecting to cloudinary without loading from .env file.
	//cld, err := cloudinary.NewFromParams("cloud_id", "Key", "secret")
	if err != nil {
		log.Fatal("Error getting cloudinary connection: ", err)
	}
	return cld
}

// UploadImage uploads an image to the cloudinary server.
func UploadImage(name string, file multipart.File) error {
	cloud := connect()
	ctx := context.Background()

	_, err := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: name,
	})
	if err != nil {
		fmt.Println("Error uploading file: ", err)
	}
	return err
}

// GetImage returns the image url for the specified image(name).
func GetImage(name string) (string, error) {
	cloud := connect()
	ctx := context.Background()

	res, err := cloud.Admin.Asset(ctx, admin.AssetParams{
		PublicID: name,
	})
	if err != nil {
		fmt.Println("Error getting file: ", err)
	}

	return res.URL, nil
}

// DeleteImage provides the neccesary functionality to delete an image from the server.
func DeleteImage(name string) (string, error) {
	cloud := connect()
	ctx := context.Background()

	res, err := cloud.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: name,
	})
	if err != nil {
		fmt.Println("Error destroying image: ", err)
	}

	return res.Result, nil

}
