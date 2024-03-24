package service

import (
	"encoding/json"
	"fmt"
	"golang-microservice-boilerplate/common"
	"golang-microservice-boilerplate/model"
	"golang-microservice-boilerplate/rest"
	"net/http"
)

func PerformPartialUpdateForBase(entityModel *model.BaseEntityModel, entityRest rest.BaseEntityRest) map[string]map[string]interface{} {
	diffMap := make(map[string]map[string]interface{})

	if entityRest.PatchMap["name"] != nil && entityModel.Name != entityRest.Name {
		diffMap = common.AddInDiffMap("name", entityModel.Name, entityRest.Name)
		entityModel.Name = entityRest.Name
	}

	if entityRest.PatchMap["displayName"] != nil && entityModel.DisplayName != entityRest.DisplayName {
		diffMap = common.AddInDiffMap("display_name", entityModel.DisplayName, entityRest.DisplayName)
		entityModel.DisplayName = entityRest.DisplayName
	}

	if entityRest.PatchMap["updatedById"] != nil && entityModel.UpdatedById != entityRest.UpdatedById {
		diffMap = common.AddInDiffMap("updated_by_id", entityModel.UpdatedById, entityRest.UpdatedById)
		entityModel.UpdatedById = entityRest.UpdatedById
	}

	if entityRest.PatchMap["updatedTime"] != nil && entityModel.UpdatedTime != entityRest.UpdatedTime {
		diffMap = common.AddInDiffMap("updated_time", entityModel.UpdatedTime, entityRest.UpdatedTime)
		entityModel.UpdatedTime = entityRest.UpdatedTime
	}

	return diffMap
}

func PerformPartialUpdateForBaseRef(entityRefModel *model.BaseEntityRefModel, refModelRest rest.BaseEntityRefModelRest) map[string]map[string]interface{} {
	diffMap := PerformPartialUpdateForBase(&entityRefModel.BaseEntityModel, refModelRest.BaseEntityRest)

	if refModelRest.PatchMap["refId"] != nil && entityRefModel.RefId != refModelRest.RefId {
		diffMap = common.AddInDiffMap("ref_id", entityRefModel.RefId, refModelRest.RefId)
		entityRefModel.RefId = refModelRest.RefId
	}

	if refModelRest.PatchMap["refModel"] != nil && entityRefModel.RefModel != refModelRest.RefModel {
		diffMap = common.AddInDiffMap("ref_model", entityRefModel.RefModel, refModelRest.RefModel)
		entityRefModel.RefModel = refModelRest.RefModel
	}

	return diffMap
}

func ConvertToBaseEntityModel(rest rest.BaseEntityRest) model.BaseEntityModel {
	return model.BaseEntityModel{
		Id:          rest.Id,
		Name:        rest.Name,
		DisplayName: rest.DisplayName,
		CreatedById: rest.CreatedById,
		CreatedTime: rest.CreatedTime,
		UpdatedById: rest.UpdatedById,
		UpdatedTime: rest.UpdatedTime,
		Removed:     rest.Removed,
	}
}

func ConvertToBaseEntityRefModel(rest rest.BaseEntityRefModelRest) model.BaseEntityRefModel {
	baseModel := ConvertToBaseEntityModel(rest.BaseEntityRest)
	return model.BaseEntityRefModel{
		BaseEntityModel: baseModel,
		RefId:           rest.RefId,
		RefModel:        rest.RefModel,
	}
}

func ConvertToBaseEntityRest(entityModel model.BaseEntityModel) rest.BaseEntityRest {
	return rest.BaseEntityRest{
		Id:          entityModel.Id,
		Name:        entityModel.Name,
		DisplayName: entityModel.DisplayName,
		CreatedById: entityModel.CreatedById,
		CreatedTime: entityModel.CreatedTime,
		UpdatedById: entityModel.UpdatedById,
		UpdatedTime: entityModel.UpdatedTime,
		Removed:     entityModel.Removed,
	}
}

func ConvertToBaseEntityRefModelRest(entityModel model.BaseEntityRefModel) rest.BaseEntityRefModelRest {
	baseRest := ConvertToBaseEntityRest(entityModel.BaseEntityModel)
	return rest.BaseEntityRefModelRest{
		BaseEntityRest: baseRest,
		RefId:          entityModel.RefId,
		RefModel:       entityModel.RefModel,
	}
}

func ConvertJsonToRequestByIdsRest(w http.ResponseWriter, r *http.Request, requestByIds rest.RequestByIdsAndModel) (rest.RequestByIdsAndModel, error) {
	body := common.GetRequestBody(r)
	err := json.Unmarshal(body, &requestByIds)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error(fmt.Sprintf("Error : %s", err.Error()), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(jsonData))
		return requestByIds, err
	}
	return requestByIds, err
}

func ConvertJsonToBulkUpdateRequest(w http.ResponseWriter, r *http.Request, bulkUpdateRequest rest.BulkUpdateRequest) (rest.BulkUpdateRequest, error) {
	body := common.GetRequestBody(r)
	err := json.Unmarshal(body, &bulkUpdateRequest)
	if err != nil {
		jsonData, _ := common.RestToJson(w, common.Error(fmt.Sprintf("Error : %s", err.Error()), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(jsonData))
		return bulkUpdateRequest, err
	}
	return bulkUpdateRequest, err
}
