package rest

type UserRest struct {
	BaseEntityRest
	Password string `json:"password"`
	Email    string `json:"email"`
}
