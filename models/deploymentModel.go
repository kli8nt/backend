package models

import "gorm.io/gorm"

type Deployment struct {
	gorm.Model
	RepoID               uint
	Repo                 Repo     `gorm:"foreignKey:RepoID;"`
	Technology           string   `json:"technology"`
	Version              string   `json:"version"`
	RepositoryURL        string   `json:"repository_url"`
	GithubToken          string   `json:"github_token"`
	ApplicationName      string   `json:"application_name"`
	RunCommand           string   `json:"run_command"`
	BuildCommand         string   `json:"build_command"`
	InstallCommand       string   `json:"install_command"`
	DependenciesFiles    []string `json:"dependencies_files"`
	IsStatic             bool     `json:"is_static"`
	OutputDirectory      string   `json:"output_directory"`
	EnvironmentVariables string   `json:"environment_variables"`
	Port                 string   `json:"port"`
}
