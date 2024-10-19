package review

import (
	"fmt"

	"main.go/types"
)

var reviewsSelectCols = `reviews.id as id, reviews.comment as comment, reviews.created_at as created_at, 
				reviews.updated_at as updated_at,reviews.rate as rate,
				users.id as user_id, users.email as email, users.name as name ,users.avatar as avatar
			`
var reviewsJoin = "LEFT JOIN users on reviews.user_id = users.id"

func convertRowsToResp(rows []types.GetAllReviewsRow) *[]types.RespAllRevs {
	fmt.Println(len(rows))
	resp := make([]types.RespAllRevs, 0)
	for _, row := range rows {
		
			var newResp = new(types.RespAllRevs)
			
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