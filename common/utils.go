package common

import (
	"encoding/json"
	"fmt"
	"golang-microservice-boilerplate/logger"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetRequestBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.ServiceLogger.Error("Error reading response body:", err)
		return nil
	}
	return body
}

func RestToJson(w http.ResponseWriter, rest any) (string, error) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, error := json.Marshal(&rest)
	if error != nil {
		return "", error
	}
	logger.ServiceLogger.Debug(string(jsonData))
	return string(jsonData), nil
}

func UploadFileToFileDB(handler *multipart.FileHeader, err error, file multipart.File) (bool, error) {
	uploadDirectory := FileDirectoryPath()
	if _, err := os.Stat(uploadDirectory); os.IsNotExist(err) {
		os.Mkdir(uploadDirectory, os.ModePerm)
	}

	filePath := filepath.Join(uploadDirectory, handler.Filename)
	newFile, err := os.Create(filePath)
	if err != nil {
		logger.ServiceLogger.Error(fmt.Sprintf("Error creating the file"))
		return false, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		logger.ServiceLogger.Error(fmt.Sprintf("Error cpoying the file"))
		return false, err
	}
	defer newFile.Close()
	return true, nil
}

func FileDirectoryPath() string {
	currentDir, _ := os.Getwd()
	return GetEnv("FILE_DB_PATH", currentDir+"/filedb")
}

func CurrentWorkingDir() string {
	WorkingDir, _ := os.Getwd()
	profile := os.Getenv("PROFILE")
	if "prod" == profile {
		WorkingDir = "/opt/service"
	}
	return WorkingDir
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func AddInDiffMap(dbKey string, oldVal any, newVal any) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		dbKey: map[string]interface{}{
			"oldvalue": oldVal,
			"newvalue": newVal,
		},
	}
}
