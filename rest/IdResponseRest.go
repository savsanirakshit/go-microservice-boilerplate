package rest

type IdResponseRest struct {
	Id      int64 `json:"result"`
	Success bool
}

type BulkUpdateRequest struct {
	Ids      []int64                `json:"ids"`
	PatchMap map[string]interface{} `json:"payload"`
}

type BulkActionResponse struct {
	SuccessIds  []int64          `json:"successIds"`
	FailedIdMap map[int64]string `json:"failedIdMap"`
}

type RequestByIdsAndModel struct {
	Ids   []int64 `json:"ids"`
	Model string  `json:"model"`
}
