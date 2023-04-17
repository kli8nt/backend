package models

import "gorm.io/gorm"

type Repo struct {
	gorm.Model
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	Branch       string `json:"branch"`
	HeadCommit   string `json:"head_commit"`
	Stack        string `json:"stack"`
	RunCommand   string `json:"run_command"`
	BuildCommand string `json:"build_command"`
	NginxPath    string `json:"nginx_path"`
	Subdomain    string `json:"subdomain"`
	K8sIP        string `json:"k8s_ip"` // or whatever
}
