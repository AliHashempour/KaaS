package model

type PostgresSpec struct {
	AppName   string   `json:"AppName"`
	Resources Resource `json:"Resources"`
	External  bool     `json:"External"`
}
