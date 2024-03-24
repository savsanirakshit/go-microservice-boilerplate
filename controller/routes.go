package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang-microservice-boilerplate/common"
	"golang-microservice-boilerplate/handler"
	"golang-microservice-boilerplate/middleware"
	"net/http"
)

func Routes(route *mux.Router) {
	route.Use(middleware.AuthMiddleware, middleware.CommonMiddleware)

	route.Methods(http.MethodGet).Path("/home").Handler(http.HandlerFunc(HomeHandler))
	route.Methods(http.MethodGet).Path("/login/{user}/{pass}").Handler(http.HandlerFunc(LoginHandler))

	//To Use Middleware to specific api
	//
	//apiRoute.Methods(http.MethodGet).Path("/protected").Handler(AuthMiddleware(http.HandlerFunc(DashboardHandler)))
	//apiRoute.Methods(http.MethodGet).Path("/protected").HandlerFunc(DashboardHandler)

	apiRoute := route.PathPrefix("/api").Subrouter()

	//User API
	userHandler := handler.NewUserHandler()
	apiRoute.Methods(http.MethodGet).Path("/user/{id}").HandlerFunc(userHandler.GetUserHandler)
	apiRoute.Methods(http.MethodPost).Path("/user").HandlerFunc(userHandler.CreateUserHandler)
	apiRoute.Methods(http.MethodPut).Path("/user/{id}").HandlerFunc(userHandler.UpdateUserHandler)
	apiRoute.Methods(http.MethodDelete).Path("/user/{id}").HandlerFunc(userHandler.DeleteUserHandler)
	apiRoute.Methods(http.MethodPost).Path("/user/search").HandlerFunc(userHandler.GetAllUserHandler)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, string("Welcome to Go-microservice"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user, _ := vars["user"]
	if user != "admin" {
		jsonData, _ := common.RestToJson(w, common.Error("Enter Valid username and password", http.StatusUnauthorized))
		fmt.Fprintf(w, string(jsonData))
		return
	}

	pass, _ := vars["pass"]
	if pass != "admin" {
		jsonData, _ := common.RestToJson(w, common.Error("Enter Valid username and password", http.StatusUnauthorized))
		fmt.Fprintf(w, string(jsonData))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	generateJWT, _ := middleware.GenerateJWT()
	m := map[string]interface{}{
		"token": generateJWT,
	}
	jsonData, _ := json.Marshal(&m)
	fmt.Fprintf(w, string(jsonData))
	return
}
