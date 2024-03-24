package model

type User struct {
	BaseEntityModel
	Password string `pg:",notnull"`
	Email    string `pg:",notnull,unique"`
}
