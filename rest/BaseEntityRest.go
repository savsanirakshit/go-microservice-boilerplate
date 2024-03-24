package rest

type BaseEntityRest struct {
	Id          int64                  `json:"id"`
	Name        string                 `json:"name"`
	DisplayName string                 `json:"displayName"`
	CreatedById int64                  `json:"createdById"`
	CreatedTime int64                  `json:"createdTime"`
	UpdatedById int64                  `json:"updatedById"`
	UpdatedTime int64                  `json:"updatedTime"`
	Removed     bool                   `json:"removed"`
	PatchMap    map[string]interface{} `json:"patchMap,omitempty"`
}

type BaseEntityRefModelRest struct {
	BaseEntityRest
	RefId    int64  `json:"refId"`
	RefModel string `json:"refModel"`
}
