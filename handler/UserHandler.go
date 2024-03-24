package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"golang-microservice-boilerplate/common"
	"golang-microservice-boilerplate/rest"
	"golang-microservice-boilerplate/service"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{Service: service.NewUserService()}
}

func (handler UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString, _ := vars["id"]
	if idString == "0" {
		jsonData, _ := common.RestToJson(w, common.Error("Error while getting package with id : "+idString, http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(jsonData))
		return
	}

	intValue, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error("Invalid package id : "+idString, http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	var includeArchive bool = false
	if isVisibleParam, err := strconv.ParseBool(r.URL.Query().Get("includeArchive")); err == nil {
		includeArchive = isVisibleParam
	}

	pkgRest, err := handler.Service.GetUserById(intValue, includeArchive)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error(err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	jsonData, _ := common.RestToJson(w, pkgRest)
	fmt.Fprintf(w, jsonData)
}

func (handler UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userRest rest.UserRest
	userRest, err := convertJsonToUserRest(w, r, userRest)
	if err != nil {
		return
	}

	id, customErr := handler.Service.Create(userRest)
	if customErr.Message != "" {
		jsonData, _ := common.RestToJson(w, common.Error("Error while creating user : "+customErr.Message, http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	jsonData, _ := common.RestToJson(w, rest.IdResponseRest{
		Id:      id,
		Success: true,
	})
	fmt.Fprintf(w, jsonData)
}

func (handler UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userRest rest.UserRest
	userRest, err := convertJsonToUserRest(w, r, userRest)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	idString, _ := vars["id"]
	if idString == "0" {
		jsonData, _ := common.RestToJson(w, common.Error("Error while getting user with id : "+idString, http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error("Error : "+err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	isUpdated, customErr := handler.Service.Update(id, userRest)
	if customErr.Message != "" {
		jsonData, _ := common.RestToJson(w, customErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	jsonData, _ := common.RestToJson(w, rest.IdResponseRest{
		Id:      id,
		Success: isUpdated,
	})
	fmt.Fprintf(w, jsonData)
}

func (handler UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString, _ := vars["id"]
	if idString == "0" {
		jsonData, _ := common.RestToJson(w, common.Error("Error while deleting user with id : "+idString, http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error("Invalid user id : "+idString, http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	var permanentDelete bool = false
	if isVisibleParam, err := strconv.ParseBool(r.URL.Query().Get("permanentDelete")); err == nil {
		permanentDelete = isVisibleParam
	}

	success, err := handler.Service.DeletePackage(id, permanentDelete)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error(err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	jsonData, _ := common.RestToJson(w, rest.IdResponseRest{
		Id:      id,
		Success: success,
	})
	fmt.Fprintf(w, jsonData)
}

func (handler UserHandler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	var searchFilter rest.SearchFilter
	searchFilter, err := rest.ConvertJsonToSearchFilter(w, r, searchFilter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pkgRest, err := handler.Service.GetAllUser(searchFilter)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error(err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonData)
		return
	}

	jsonData, _ := common.RestToJson(w, pkgRest)
	fmt.Fprintf(w, jsonData)
}

func convertJsonToUserRest(w http.ResponseWriter, r *http.Request, pkgRest rest.UserRest) (rest.UserRest, error) {
	body := common.GetRequestBody(r)
	err := json.Unmarshal(body, &pkgRest)
	v := validator.New()
	err = v.Struct(pkgRest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			jsonData, _ := common.RestToJson(w, common.Error(fmt.Sprintf("Validation error on field %s, Expected values : %s", err.Field(), err.Param()), http.StatusInternalServerError))
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, string(jsonData))
			break
		}
	}
	var patchMap map[string]interface{}
	json.Unmarshal([]byte(body), &patchMap)
	pkgRest.PatchMap = patchMap
	return pkgRest, err
}
