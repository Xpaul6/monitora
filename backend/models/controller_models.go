package models

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

type GetAllServersResponse struct {
	Count   int64    `json:"count"`
	Servers []Server `json:"servers"`
}

type AddServerRequest struct {
	Name string `json:"name" binding:"required"`
	IP   string `json:"ip" binding:"required"`
}

type DeleteServerRequest struct {
	ID uint `json:"id" binding:"required"`
}
