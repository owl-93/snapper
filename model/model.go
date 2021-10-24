package model

type MetaTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SnapperRequest struct {
	Page string `json:"page"`
}
