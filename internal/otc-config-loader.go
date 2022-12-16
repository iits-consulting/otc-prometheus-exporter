package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	Clouds []Cloud
}

type Domain struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Token struct {
	Secret    string `json:"secret"`
	IssuedAt  string `json:"issued_at"`
	ExpiresAt string `json:"expires_at"`
}

type Project struct {
	Name        string `json:"name"`
	Id          string `json:"id"`
	ScopedToken Token
}

type Cluster struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Cloud struct {
	Domain        Domain
	UnscopedToken Token
	Projects      []Project
	Clusters      []Cluster
}

func expandUserHome(path string) string {

	usr, _ := user.Current()
	dir := usr.HomeDir

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	return path
}

func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(expandUserHome(path))
	if err != nil {
		return nil, err
	}

	var config Config

	err = json.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetScopedToken(c Config, projectname string) (*Token, error) {
	for _, cloud := range c.Clouds {
		for _, project := range cloud.Projects {
			if project.Name == projectname {
				return &project.ScopedToken, nil
			}
		}

	}

	return nil, fmt.Errorf("no such Project \"%v\"", projectname)
}
