package rest

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang-microservice-boilerplate/common"
	"net/http"
	"strconv"
	"strings"
)

type SearchFilter struct {
	Offset        int             `json:"offset"`
	Size          int             `json:"size"`
	Archived      bool            `json:"archived"`
	SortBy        string          `json:"sortBy"`
	Qualification []Qualification `json:"qualification"`
}

type Qualification struct {
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

func PrepareQueryFromSearchFilter(filter SearchFilter, tableName string, isCountQual bool) string {
	columnSelection := "*"
	if isCountQual {
		columnSelection = "count(*)"
	}
	q := fmt.Sprintf("select %s from %s ", columnSelection, tableName)
	q += fmt.Sprintf(" where removed = %s ", strconv.FormatBool(filter.Archived))
	quals := filter.Qualification
	for i, qual := range quals {
		if i == 0 {
			q += " AND "
		} else {
			q += " OR "
		}
		column := qual.Column
		column = common.ToSnakeCase(column)
		operator := qual.Operator
		value := qual.Value

		if strings.Contains("contains", strings.ToLower(operator)) {
			q += "\"" + column + "\" ILIKE '%" + value + "%' "
		} else if strings.Contains("equals", strings.ToLower(operator)) {
			q += "\"" + column + "\" = '" + value + "' "
		}
	}

	if !isCountQual {
		if filter.SortBy != "" {
			by := filter.SortBy
			direction := "DESC"
			if strings.Contains(by, "-") {
				direction = "ASC"
				by = strings.ReplaceAll(by, "-", "")
			}
			by = common.ToSnakeCase(by)
			q += fmt.Sprintf("order by \"%s\" %s ", by, direction)
		} else {
			q += fmt.Sprintf("order by \"%s\" %s ", "id", "DESC")
		}

		if filter.Size > 0 {
			q += fmt.Sprintf("offset %v limit %v ", filter.Offset, filter.Size)
		}
	}
	q += ";"
	return q
}

func ConvertJsonToSearchFilter(w http.ResponseWriter, r *http.Request, searchFilter SearchFilter) (SearchFilter, error) {
	body := common.GetRequestBody(r)
	err := json.Unmarshal(body, &searchFilter)
	v := validator.New()
	err = v.Struct(searchFilter)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			jsonData, _ := common.RestToJson(w, common.Error(fmt.Sprintf("Validation error on field %s, Expected values : %s", err.Field(), err.Param()), http.StatusInternalServerError))
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, string(jsonData))
			break
		}
	}
	return searchFilter, err
}
