package model

type BaseEntityModel struct {
	Id          int64  `pg:"type:serial,pk"`
	Name        string `pg:",notnull"`
	DisplayName string
	CreatedById int64
	CreatedTime int64
	UpdatedById int64
	UpdatedTime int64
	Removed     bool `pg:",use_zero"`
}

type BaseEntityRefModel struct {
	BaseEntityModel
	RefId    int64
	RefModel string
}
