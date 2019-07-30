package utils

import (
	"github.com/corona10/goimagehash"

	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func ImageDistance(image1 string, image2 string)(res int, err error){

	img1, err := GetImageFromFile(image1)
	if err != nil {
		return -1, err
	}
	img2, err := GetImageFromFile(image2)
	if err != nil {
		return -1, err
	}

	if img1 != nil && img2 != nil {

		width, height := 8, 8
		hash1, _ := goimagehash.ExtAverageHash(img1, width, height)
		hash2, _ := goimagehash.ExtAverageHash(img2, width, height)

		res, err = hash1.Distance(hash2)
	}

	return res, err
}

func GetImageFromFile(imagePath string)(image.Image, error){

	var img image.Image
	var err error

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	ext := filepath.Ext(imagePath)

	if ext == ".jpg" || ext == ".jpeg" {
		img, err = jpeg.Decode(file)
	}

	if ext == ".png" {
		img, err = png.Decode(file)
	}

	return img, err

}