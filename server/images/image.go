// +build !vips

package images

import (
	"github.com/joonazan/imagick/imagick"
	"github.com/valyala/fasthttp"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Initialize the image library
func init() {
	imagick.Initialize()
}

// Imageblob is just an image file dumped, byte by byte to an byte array.
type ImageBlob []byte

// Image is an uncompressed image that must be convertd to blob before serving to a client.
type Image struct {
	*imagick.MagickWand
}

func NewImage() Image {
	return Image{imagick.NewMagickWand()}
}

// Clone an image. Remember images and made clones need to be destroyed using Destroy().
func (img Image) Clone() Image {
	return Image{img.MagickWand.Clone()}
}

// Converts the image to different format. Takes extension as parameter.
func (img Image) Convert(ext string) (err error) {
	err = img.SetImageFormat(ext)
	if err != nil {
		err = errors.New("Could not convert image: " + err.Error())
	}
	return
}

// Returns image width
func (img Image) GetWidth() (width int) {
	return int(img.GetImageWidth())
}

// Returns image height
func (img Image) GetHeight() (height int) {
	return int(img.GetImageHeight())
}

// Returns filename extension of the image e.g. jpg, gif, webp
func (img Image) GetExtension() (ext string) {
	format := img.GetImageFormat()
	ext = strings.ToLower(format)
	if strings.EqualFold(ext, "jpeg") {
		ext = "jpg"
	}
	return
}

// Method for converting Image to ImageBlob.
func (img Image) Blob() ImageBlob {
	return img.GetImageBlob()
}

// Returns Image from file.
func (img Image) FromFile(filename string) error {

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	// Remember to free resources after you're done
	defer reader.Close()

	buffer := bytes.NewBuffer([]byte{})
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return err
	}
	blob := ImageBlob(buffer.Bytes())

	return img.FromBlob(blob)
}

// Return binary ImageBlob of an image from web.
func (img Image) FromWeb(url string) error {

	//resp, err := http.Get(url)
	statuscode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return err
	}
	//defer resp.Body.Close()

	if statuscode != 200 {
		return errors.New(fmt.Sprintf("Couldn't load image. Server returned %i", statuscode))
	}

	return img.FromBlob(body)
}

// Image from blob
func (img Image) FromBlob(blob ImageBlob) error {
	return img.ReadImageBlob(blob)
}

// Resize image
func (img Image) Resize(w, h int) error {
	return img.ResizeImage(uint(w), uint(h), imagick.FILTER_LANCZOS, 1)
}
