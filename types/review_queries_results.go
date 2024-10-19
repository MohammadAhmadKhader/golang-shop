package types

import "time"

type GetAllReviewsRow struct {
	Id        uint
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Rate      uint8
	UserId    uint
	Email     string
	Name      string
	Avatar    string
}

type RespAllRevs struct {
	Id        uint      `json:"id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Rate      uint8     `json:"rate"`
	User      RowReviewUser   `json:"user"`
}

type RowReviewUser struct {
	UserId uint   `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}