package repository

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v10"
	"golang-microservice-boilerplate/db"
	"golang-microservice-boilerplate/model"
)

type UserRepository struct {
	dbConnection *pg.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		dbConnection: db.Conn,
	}
}

func (repo UserRepository) Create(usr *model.User) (int64, error) {
	_, err := repo.dbConnection.Model(usr).Returning("id").Insert()
	if err != nil {
		return 0, err
	}
	return usr.Id, nil
}

func (repo UserRepository) Update(usr *model.User) (int64, error) {
	_, err := repo.dbConnection.Model(usr).WherePK().Update()
	if err != nil {
		return 0, err
	}
	return usr.Id, nil
}

func (repo UserRepository) GetById(usrId int64, includeArchive bool) (model.User, error) {
	var usr model.User
	var err error
	if includeArchive {
		err = repo.dbConnection.Model(&usr).Where("id = ?", usrId).Select()
	} else {
		err = repo.dbConnection.Model(&usr).Where("id = ?", usrId).Where("removed = ?", false).Select()
	}
	if err != nil {
		return usr, err
	}
	return usr, nil
}

func (repo UserRepository) DeleteById(usrId int64) (bool, error) {
	usr, err := repo.GetById(usrId, false)
	if err != nil {
		return false, err
	} else if usr.Id == 0 {
		return false, errors.New(fmt.Sprintf("User not found with id : %v", usrId))
	}
	usr.Removed = true
	_, err = repo.dbConnection.Model(&usr).WherePK().Update()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo UserRepository) PermanentDeleteById(pkgId int64) (bool, error) {
	pkg := new(model.User)
	pkg.Id = pkgId
	_, err := repo.dbConnection.Model(pkg).WherePK().Delete()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo UserRepository) GetAllUser(query string) ([]model.User, error) {
	var usr []model.User
	_, err := repo.dbConnection.Query(&usr, query)
	if err != nil {
		return usr, err
	}
	return usr, nil
}

func (repo UserRepository) Count(query string) int {
	var count int
	_, err := repo.dbConnection.QueryOne(pg.Scan(&count), query)
	if err != nil {
		return 0
	}
	return count
}
