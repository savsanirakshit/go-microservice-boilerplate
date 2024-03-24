package rest

type ListResponseRest struct {
	ObjectList interface{} `json:"result"`
	TotalCount int         `json:"totalCount"`
}
