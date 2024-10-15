package review

import (
	"fmt"
	"time"
)

var reviewsSelectCols = `reviews.id as id, reviews.comment as comment, reviews.created_at as created_at, 
				reviews.updated_at as updated_at,reviews.rate as rate,
				users.id as user_id, users.email as email, users.name as name ,users.avatar as avatar
			`
var reviewsJoin = "LEFT JOIN users on reviews.user_id = users.id"

type GetAllReviewsRow struct {
	Id            uint 
	Comment       string 
	CreatedAt     time.Time 
	UpdatedAt     time.Time 
	Rate uint8
	UserId       uint
	Email      string
	Name 		string
	Avatar     string
}

type respAllRevs struct {
	Id            uint `json:"id"`
	Comment       string `json:"comment"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Rate uint8 `json:"rate"`
	User rowUser `json:"user"`
}

type rowUser struct {
	UserId uint `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
}

func convertRowsToResp(rows []GetAllReviewsRow) *[]respAllRevs {
	fmt.Println(len(rows))
	resp := make([]respAllRevs, 0)
	for _, row := range rows {
		
			var newResp = new(respAllRevs)
			
			newResp.Id = row.Id
			newResp.Comment = row.Comment
			newResp.Rate = row.Rate
			newResp.CreatedAt = row.CreatedAt
			newResp.UpdatedAt = row.UpdatedAt
			newResp.User.Name = row.Name
			newResp.User.Email = row.Email
			newResp.User.Avatar = row.Avatar
			newResp.User.UserId = row.UserId
			resp = append(resp, *newResp)
		
	}

	return &resp
}