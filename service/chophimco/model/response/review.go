package response

import "time"

type ReviewResponse struct {
	ID          int       `json:"id"`
	UserName    string    `json:"user_name"`
	ProductID   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}
