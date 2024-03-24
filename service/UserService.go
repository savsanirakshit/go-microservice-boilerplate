package service

import (
	"fmt"
	"golang-microservice-boilerplate/common"
	"golang-microservice-boilerplate/logger"
	"golang-microservice-boilerplate/model"
	"golang-microservice-boilerplate/repository"
	"golang-microservice-boilerplate/rest"
	"net/http"
)

type UserService struct {
	Repository *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		Repository: repository.NewUserRepository(),
	}
}

func (service UserService) convertToModel(restModel rest.UserRest) *model.User {
	return &model.User{
		BaseEntityModel: ConvertToBaseEntityModel(restModel.BaseEntityRest),
		Password:        restModel.Password,
		Email:           restModel.Email,
	}
}

func (service UserService) convertToRest(domainModel model.User) rest.UserRest {
	return rest.UserRest{
		BaseEntityRest: ConvertToBaseEntityRest(domainModel.BaseEntityModel),
		Password:       domainModel.Password,
		Email:          domainModel.Email,
	}
}

func (service UserService) convertListToRest(users []model.User) []rest.UserRest {
	var userRestList []rest.UserRest
	if len(users) != 0 {
		for _, user := range users {
			userRest := service.convertToRest(user)
			userRestList = append(userRestList, userRest)
		}
	}
	return userRestList
}

func (service UserService) GetUserById(id int64, includeArchive bool) (rest.UserRest, error) {
	logger.ServiceLogger.Info(fmt.Sprintf("Process started to get user for id %v", id))
	var userRest rest.UserRest
	usr, err := service.Repository.GetById(id, includeArchive)
	if err != nil {
		logger.ServiceLogger.Error(fmt.Sprintf("Error while get user for id - %v, Error : %s ", id, err.Error()))
		return userRest, err
	}
	logger.ServiceLogger.Info(fmt.Sprintf("Process Completed to get user for id %v", id))
	return service.convertToRest(usr), nil
}

func (service UserService) DeletePackage(id int64, permanentDelete bool) (bool, error) {
	logger.ServiceLogger.Info("Process started to delete user for id - ", id)
	pkg, err := service.Repository.GetById(id, permanentDelete)
	if err != nil {
		logger.ServiceLogger.Error(fmt.Sprintf("Error while get user to delete for id - %v, Error : %s ", id, err.Error()))
		return false, err
	}

	success := false
	if permanentDelete {
		success, err = service.Repository.PermanentDeleteById(id)
	} else {
		success, err = service.Repository.DeleteById(id)
	}

	if err != nil {
		logger.ServiceLogger.Error(fmt.Sprintf("Error while delete user for id - %v, Error : %s ", id, err.Error()))
		return success, err
	}

	go service.AfterDelete(pkg, permanentDelete)

	logger.ServiceLogger.Info("Process Completed to delete user for id - ", id)
	return true, nil
}

func (service UserService) AfterDelete(usr model.User, permanentDelete bool) {
	defer func() {
		if err := recover(); err != nil {
			logger.ServiceLogger.Error("Error while performing after delete process for user id : ", usr.Id, ", permanentDelete :", permanentDelete)
		}
	}()

	//TODO : Handle after delete
}

func (service UserService) Create(rest rest.UserRest) (int64, common.CustomError) {
	logger.ServiceLogger.Info(fmt.Sprintf("Process started to create package"))
	customErr := service.BeforeCreate(rest)
	if customErr.Message != "" {
		logger.ServiceLogger.Error(fmt.Sprintf("Error while before creating package ,Error : %s ", customErr.Message))
		return 0, customErr
	}

	user := service.convertToModel(rest)
	id, err := service.Repository.Create(user)
	if err != nil {
		logger.ServiceLogger.Error(fmt.Sprintf("Error while creating user ,Error : %s ", err.Error()))
		return 0, common.CustomError{Message: err.Error()}
	}

	go service.AfterCreate(user)

	logger.ServiceLogger.Info(fmt.Sprintf("Create user process completed successfully"))
	return id, customErr
}

func (service UserService) BeforeCreate(userRest rest.UserRest) common.CustomError {
	logger.ServiceLogger.Debug("After create for user : ", userRest.Name)
	//TODO Handle Before Create
	return common.CustomError{}
}

func (service UserService) AfterCreate(user *model.User) {
	defer func() {
		if err := recover(); err != nil {
			logger.ServiceLogger.Error("Error while performing after create process for user id : ", user.Id)
		}
	}()

	//TODO : Handle after create
}

func (service UserService) Update(id int64, restModel rest.UserRest) (bool, common.CustomError) {
	logger.ServiceLogger.Info(fmt.Sprintf("Process started to update user with id - %v", id))
	user, err := service.Repository.GetById(id, false)
	if err != nil {
		logger.ServiceLogger.Error(fmt.Sprintf("Error while updating user for id - %v ,Error : %s ", id, err.Error()))
		return false, common.CustomError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	customErr := service.BeforeUpdate(user, restModel)
	if customErr.Message != "" {
		logger.ServiceLogger.Error(fmt.Sprintf("Error while updating package for id - %v ,Error : %s ", id, customErr.Message))
		return false, customErr
	}

	diffMap, isUpdatable := service.performPartialUpdate(&user, restModel)
	if isUpdatable {
		_, err := service.Repository.Update(&user)
		if err != nil {
			logger.ServiceLogger.Error(fmt.Sprintf("Error while updating user for id - %v ,Error : %s ", id, err.Error()))
			return false, common.CustomError{Message: err.Error(), Code: http.StatusInternalServerError}
		}

		service.AfterUpdate(diffMap, user, restModel)
		logger.ServiceLogger.Info(fmt.Sprintf("Process Completed to update user with id - %v", id))
		return true, common.CustomError{}
	} else {
		logger.ServiceLogger.Info(fmt.Sprintf("No fields need to updated"))
		logger.ServiceLogger.Info(fmt.Sprintf("Process Completed to update user with id - %v", id))
		return isUpdatable, customErr
	}
}

func (service UserService) BeforeUpdate(domainModel model.User, restModel rest.UserRest) common.CustomError {
	logger.ServiceLogger.Debug("After create for user name : ", restModel.Name, ", id : ", domainModel.Id)
	//TODO: handle before update
	return common.CustomError{}
}

// Create diffMap to handle after update process to specific field
func (service UserService) performPartialUpdate(model *model.User, restModel rest.UserRest) (map[string]map[string]interface{}, bool) {
	diffMap := PerformPartialUpdateForBase(&model.BaseEntityModel, restModel.BaseEntityRest)

	if restModel.PatchMap["password"] != nil && model.Password != restModel.Password {
		diffMap = common.AddInDiffMap("password", model.Password, restModel.Password)
		model.Password = restModel.Password
	}

	if restModel.PatchMap["email"] != nil && model.Email != restModel.Email {
		diffMap = common.AddInDiffMap("email", model.Email, restModel.Email)
		model.Email = restModel.Email
	}

	return diffMap, len(diffMap) != 0
}

func (service UserService) AfterUpdate(diffMap map[string]map[string]interface{}, user model.User, restModel rest.UserRest) {
	defer func() {
		if err := recover(); err != nil {
			logger.ServiceLogger.Error("Error while performing after update process for user id : ", user.Id, ", updated field :", len(diffMap), ", patch map size : ", len(restModel.PatchMap))
		}
	}()

	//Todo: handle after update
}

func (service UserService) GetAllUser(filter rest.SearchFilter) (rest.ListResponseRest, error) {
	countQuery := rest.PrepareQueryFromSearchFilter(filter, "users", true)
	var responsePage rest.ListResponseRest
	var packagePageList []model.User
	var err error
	count := service.Repository.Count(countQuery)
	if count > 0 {
		searchQuery := rest.PrepareQueryFromSearchFilter(filter, "users", false)
		packagePageList, err = service.Repository.GetAllUser(searchQuery)
		if err != nil {
			return responsePage, err
		}
		responsePage.ObjectList = service.convertListToRest(packagePageList)
	} else {
		responsePage.ObjectList = make([]interface{}, 0)
	}
	responsePage.TotalCount = count
	return responsePage, nil
}
