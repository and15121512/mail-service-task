package models

const (
	OkAuthRespStatus        = iota
	RefreshedAuthRespStatus = iota
	RefusedAuthRespStatus   = iota
)

type AuthResult struct {
	Status       int
	AccessToken  string
	RefreshToken string
	Login        string
}
