package dto

type RegisterRequestBody struct {
	PhoneNumber string `json:"phone_number" binding:"required,phone_number"`
	Password    string `json:"password" binding:"required"`
}

type LoginRequestBody = RegisterRequestBody

type LoginResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenRequestBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse = LoginResponseBody
