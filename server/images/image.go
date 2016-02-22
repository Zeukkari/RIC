package images

import (
	"errors"
	"github.com/gographics/imagick/imagick"
	"strings"
)

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
