package models

// "type": "text",
// "contents": "Anonymous hiker",
// "name": "Credits",
// "image": ""

type Media struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"       binding:"required"`
	Contents string `json:"contents"       binding:"required"`
	Type     string `json:"type"       binding:"required"`
	ImageURL string `json:"image"       binding:"required"`
}
