package model

type Application struct {
	AppName      string   `json:"AppName"`
	Replicas     int32    `json:"Replicas"`
	ImageAddress string   `json:"ImageAddress"`
	ImageTag     string   `json:"ImageTag"`
	ServicePort  int      `json:"ServicePort"`
	Resources    Resource `json:"Resources"`
	Envs         []Env    `json:"Envs"`
}
