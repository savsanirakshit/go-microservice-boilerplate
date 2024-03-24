package handler

import (
	"golang-microservice-boilerplate/logger"
	"golang-microservice-boilerplate/rest"
	"golang-microservice-boilerplate/service"
)

func SystemOnBoardService() {
	createAdminUser()
}

func createAdminUser() {
	logger.ServiceLogger.Info("Process started to create admin user")
	userRest := rest.UserRest{
		BaseEntityRest: rest.BaseEntityRest{Name: "admin"},
		Password:       "admin",
		Email:          "admin",
	}

	_, err := service.NewUserService().Create(userRest)
	if err.Message != "" {
		logger.ServiceLogger.Warn("admin user not created : ", err.Message)
	} else {
		logger.ServiceLogger.Debug("admin user created successfully.")
	}
	logger.ServiceLogger.Info("Process Completed to create admin user")
}
