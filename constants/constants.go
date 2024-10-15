package constants

import "main.go/types"

const (
	Prefix       = "/api/v1"
	TokenPayload = types.TokenKey("tokenPayload")
	UserKey = types.UserKey("UserKey")
	ResourceKey = types.AuthorizedResource("AuthorizedResource")
)

// This contains the columns that can be changed by the user, it's used for create and update processes
//
// It's set as var because golang do not allow slices as constant
var (
	ProductCols = []string{"Name","Quantity","Description","CategoryID","Price"}
	CategoryCols = []string{"Name"}
	ImageCols = []string{"ProductID","ImageUrl","IsMain","ImagePublicId"}
	IdUrlPathKey = "id"
	CommentCols = []string{"Comment","Rate"}
	UserCreateCols = []string{"Name", "Email", "Password"}
	UserUpdateCols = []string{"Name", "Email", "MobileNumber"}
	RoleCols = []string{"Role"}
	AddressCols = []string{"FullName","Country","StreetAddress","City","ZipCode","State"}
)

 
