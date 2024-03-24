package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang-microservice-boilerplate/common"
	"golang-microservice-boilerplate/controller"
	"golang-microservice-boilerplate/db"
	"golang-microservice-boilerplate/handler"
	"golang-microservice-boilerplate/logger"
	"net/http"
)

func main() {
	WorkingDir := common.CurrentWorkingDir()
	//Load app.config file
	err := godotenv.Load(WorkingDir + "/app.config")
	if err != nil {
		fmt.Printf("Error loading app.config file %s", err.Error())
	}

	fmt.Printf("Deployment Server started...")
	logDir := common.GetEnv("LOG_DIR", WorkingDir+"/logs")
	logLevel := common.GetEnv("LOG_LEVEL", "info")
	logger.ConfigLogger(logDir, logLevel)

	logger.ServiceLogger.Info("Init DB Connection...")
	_, err = db.Connect()
	if err != nil {
		panic(err)
	}
	logger.ServiceLogger.Info("DB successfully connected")

	logger.ServiceLogger.Info("Start Processing of system on board")
	handler.SystemOnBoardService()
	logger.ServiceLogger.Info("End Processing of system on board")

	router := mux.NewRouter()
	controller.Routes(router)

	logger.ServiceLogger.Info("Server started on port 8088")

	http.ListenAndServe(":8088", router)

	// TODO : Use this for http2
	//
	//server := &http.Server{Addr: ":8088", Handler: router}
	//
	//http2.ConfigureServer(server, &http2.Server{})
	//
	//server.ListenAndServeTLS(WorkingDir+"/server.pem", WorkingDir+"/server.key")

	defer db.Conn.Close()
}
