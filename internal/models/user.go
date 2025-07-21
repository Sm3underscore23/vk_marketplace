package models

type UserVerify struct {
	UserID       int
	PasswordHash string
}

type UserData struct {
	Login    string
	Password string
}
