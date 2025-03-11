package main

import (
	"main.go/cmd/seed/seed_data"
)

func main() {
	setFakeUsers := false
	setFakeProducts := false
	setFakeReviews := true
	setFakeCategories := false

	if setFakeUsers {
		seed_data.SeedUsers()
	}

	if setFakeCategories {
		seed_data.SeedCategories()
	}

	if setFakeProducts {
		seed_data.SeedProducts()
	}

	if setFakeReviews {
		seed_data.SeedReviews()
	}
}
