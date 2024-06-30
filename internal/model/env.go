package model

type Env struct {
	Key      string `json:"Key"`
	Value    string `json:"Value"`
	IsSecret bool   `json:"IsSecret"`
}
