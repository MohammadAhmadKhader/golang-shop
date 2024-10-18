package main

import (
	"main.go/cmd/seed/seed_data"
)

func main() {
	setFakeUsers := false
	setFakeProducts := false
	setFakeReviews := false

	if setFakeUsers {
		seed_data.SeedUsers()
	}

	if setFakeProducts {
		seed_data.SeedProducts()
	}

	if setFakeReviews {
		seed_data.SeedReviews()
	}

}
