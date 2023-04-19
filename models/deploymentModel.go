package models

import "gorm.io/gorm"

type Deployment struct {
	gorm.Model
	RepoID       uint
	Repo         Repo   `gorm:"foreignKey:RepoID;"`
	Stack        string `json:"stack"`
	RunCommand   string `json:"run_command"`
	BuildCommand string `json:"build_command"`
	NginxPath    string `json:"nginx_path"`
	Subdomain    string `json:"subdomain"`
	K8sIP        string `json:"k8s_ip"` // or whatever
}
