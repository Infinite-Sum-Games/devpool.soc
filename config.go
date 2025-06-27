package main

type AppConfig struct {
	AppID          int64
	PrivateKeyPath string
	Environment    string
	Repos          []RepoConfig
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

type RepoConfig struct {
	InstallationID int64
	Owner          string
	Repo           string
}

func NewRepoConfig() *RepoConfig {
	return &RepoConfig{}
}
