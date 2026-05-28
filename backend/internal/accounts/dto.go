// Package accounts has Data Transfer Objects
package accounts

type CreateAccountRequest struct {
	Name string `json:"name"`
}

type RenameAccountRequest struct {
	Name string `json:"name"`
}

type AccountResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"is_archived"`
	CreatedAt  string `json:"created_at"`
}

type CreateAccountResponse struct {
	ID int64 `json:"id"`
}
