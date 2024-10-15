package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

var cldinary *cloudinary.Cloudinary

func init() {
	cld , err := cloudinary.NewFromParams(Envs.CLOUDINARY_NAME,Envs.CLOUDINARY_APIKEY,Envs.CLOUDINARY_SECRET)
	if err != nil {
		log.Fatal(err)
	}

	cldinary = cld
}

func GetCloudinary() *cloudinary.Cloudinary {
	return cldinary
}