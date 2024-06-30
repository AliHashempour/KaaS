package model

type ApplicationRequest struct {
	AppName      string   `json:"AppName"`
	Replicas     int      `json:"Replicas"`
	ImageAddress string   `json:"ImageAddress"`
	ImageTag     string   `json:"ImageTag"`
	ServicePort  int      `json:"ServicePort"`
	Resources    Resource `json:"Resources"`
	Envs         []EnvVar `json:"Envs"`
}
