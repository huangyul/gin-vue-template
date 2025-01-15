package dto

type ListResp struct {
	Data  any   `json:"data"`
	Total int64 `json:"total"`
}

type QuerySelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
