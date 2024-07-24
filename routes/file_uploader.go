package routes

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/LitPad/backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
)

var sess *session.Session

func init() {
	var err error
	// Configure the S3 client
	sess, err = session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"), // Dummy region, replace if necessary
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Endpoint:    aws.String(cfg.S3EndpointUrl),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
}

// uploadToCloudinary uploads the file to Cloudinary and returns the URL of the uploaded file
func UploadFile(fileHeader *multipart.FileHeader, key string, folder string) {
	file, err := fileHeader.Open()
	if err != nil {
		log.Println("failed to open file")
	}
	defer file.Close()
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(folder),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Println("failed to upload file")
	}	
}

func ValidateImage(c *fiber.Ctx, name string, required bool) (*multipart.FileHeader, *utils.ErrorResponse) {
	file, err := c.FormFile(name)

	data := map[string]string{
		name: "Invalid image type",
	}
	errData := utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data)

	if required && err != nil {
		data[name] = "Image is required"
		errData = utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data)
		return nil, &errData
	}

	// Open the file
	if file != nil {
		fileHandle, err := file.Open()
		if err != nil {
			return nil, &errData
		}
		
		defer fileHandle.Close()

		// Read the first 512 bytes for content type detection
		buffer := make([]byte, 512)
		_, err = fileHandle.Read(buffer)
		if err != nil {
			return nil, &errData
		}

		// Detect the content type
		contentType := http.DetectContentType(buffer)
		switch contentType {
			case "image/jpeg", "image/png", "image/gif":
				return file, nil
		}
		return nil, &errData
	}
	return nil, nil
}
