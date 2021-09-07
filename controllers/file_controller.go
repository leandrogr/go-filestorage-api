package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

type File struct {
	Name         string
	Type         string
	Size         string
	LastModified time.Time
}

func CreateFile(c *gin.Context) {

	sess := c.MustGet("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)

	r := c.Request

	directory := r.FormValue("directory")
	r.ParseMultipartForm(32 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	contentType := http.DetectContentType(buffer)
	filepath := directory + "/" + fileHeader.Filename

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(os.Getenv("AWS_BUCKET")),
		Key:                aws.String(filepath),
		Body:               file,
		ACL:                aws.String("public-read"),
		ContentDisposition: aws.String("attachment"),
		ContentType:        aws.String(contentType),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to upload file",
			"message":  err.Error(),
			"uploader": up,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filepath": filepath,
	})

}

func ShowFiles(c *gin.Context) {

	sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess)

	directory := c.Query("directory")

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Prefix: aws.String(directory),
	})
	if err != nil {
		message := fmt.Sprintf("Unable to list items in bucket %q \n %v", os.Getenv("AWS_BUCKET"), err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}

	mySlice := []File{}
	for _, item := range resp.Contents {

		head, err := svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET")),
			Key:    aws.String(*item.Key),
		})
		if err != nil {
			message := fmt.Sprintf("Unable to get head items in bucket %q \n %v", os.Getenv("AWS_BUCKET"), err)
			c.JSON(400, gin.H{
				"message": message,
			})
			return
		}

		mySlice = append(mySlice, File{
			Name:         aws.StringValue(item.Key),
			Type:         *aws.String(*head.ContentType),
			Size:         strconv.FormatInt(*item.Size, 10) + " Bytes",
			LastModified: *item.LastModified,
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"files": mySlice,
	})

}

func GetFile(c *gin.Context) {

	sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess)

	key := c.Query("key")

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Prefix: aws.String(key),
	})
	if err != nil {
		message := fmt.Sprintf("Unable to list items in bucket %q \n %v", os.Getenv("AWS_BUCKET"), err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}

	mySlice := []File{}
	for _, item := range resp.Contents {

		head, err := svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET")),
			Key:    aws.String(*item.Key),
		})
		if err != nil {
			message := fmt.Sprintf("Unable to get head items in bucket %q \n %v", os.Getenv("AWS_BUCKET"), err)
			c.JSON(400, gin.H{
				"message": message,
			})
			return
		}

		mySlice = append(mySlice, File{
			Name:         aws.StringValue(item.Key),
			Type:         *aws.String(*head.ContentType),
			Size:         strconv.FormatInt(*item.Size, 10) + " Bytes",
			LastModified: *item.LastModified,
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"files": mySlice,
	})

}

func DeleteFile(c *gin.Context) {

	sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess)

	r := c.Request
	filePath := r.FormValue("filepath")

	_, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(filePath),
	})
	if err != nil {
		message := fmt.Sprintf("File not exist at bucket %q \n %v", os.Getenv("AWS_BUCKET"), err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(filePath),
	})
	if err != nil {
		message := fmt.Sprintf("File not deleted at bucket %q \n %v", os.Getenv("AWS_BUCKET"), err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted with sussess!",
	})

}

func CopyFile(c *gin.Context) {

	sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess)

	r := c.Request
	origem := r.FormValue("origem")
	destino := r.FormValue("destino")
	bucket := os.Getenv("AWS_BUCKET")
	destinoPath := destino + "/" + origem
	origemPath := bucket + "/" + origem
	_, err := svc.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(os.Getenv("AWS_BUCKET")),
		CopySource: aws.String(url.PathEscape(origemPath)),
		Key:        aws.String(destinoPath),
	})
	if err != nil {
		message := fmt.Sprintf("Unable to copy item from directory %q to directory %q, %v", destino, origem, err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(destinoPath),
	})
	if err != nil {
		message := fmt.Sprintf("Item not copy from directory %q to directory %q, %v", destino, origem, err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(origem),
	})
	if err != nil {
		message := fmt.Sprintf("File not deleted at bucket  %q \n %v", os.Getenv("AWS_BUCKET"), err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File copied with sussess!",
	})

}

func DownloadFile(c *gin.Context) {

	sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess)

	item := c.Query("item")

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(item),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		message := fmt.Sprintf("Unable to url item %q, %v", item, err)
		c.JSON(400, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dowload_url": urlStr,
	})

}
