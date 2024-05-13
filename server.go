package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tj/go-dropbox"
)

var dropboxToken = ""

func uploadToDropbox(file io.Reader, filename string) error {
	client := dropbox.New(dropbox.NewConfig(dropboxToken))
	args := dropbox.UploadInput{
		Path:   "/" + filename,
		Mode:   dropbox.WriteModeAdd,
		Mute:   true,
		Reader: file,
	}
	_, err := client.Files.Upload(&args)
	return err
}

func upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Check MIME type for image validation
	buffer := make([]byte, 512) // Only need the first 512 bytes to sniff the content type
	if _, err := src.Read(buffer); err != nil {
		return err
	}
	src.Seek(0, io.SeekStart) // Reset the read pointer after sniffing

	contentType := http.DetectContentType(buffer)
	if !isImage(contentType) {
		return c.HTML(http.StatusBadRequest, "<p>Only image files are allowed.</p>")
	}

	// Create a unique filename using timestamp
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s-%s", timestamp, file.Filename)

	// Upload to Dropbox
	if err := uploadToDropbox(src, filename); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully.</p>", filename))
}

func isImage(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/jpg", "image/gif", "image/png", "image/bmp", "image/svg+xml":
		return true
	default:
		return false
	}
}

func main() {
	dropboxToken = os.Getenv("DROPBOX_TOKEN")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":1323"))
}
