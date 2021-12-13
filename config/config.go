package config

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type FileindexHandlerEndpointStruct struct {
	View string `yaml:"view,omitempty"`
}
type FileindexHandlerStruct struct {
	BasePath string `yaml:"base_path"`
	Hide []string `yaml:"hide,omitempty"`
	Endpoints map[string]FileindexHandlerEndpointStruct `yaml:"endpoints"`
}

type GraphvizEndpointStruct struct {
	View string `yaml:"view,omitempty"`
	Edit string `yaml:"edit,omitempty"`
}
type GraphvizStruct struct {
	Endpoints map[string]GraphvizEndpointStruct `yaml:"endpoints"`
}

type CardStruct struct {
	Title string `yaml:"title"`
	Description string `yaml:"description"`
	View string `yaml:"view,omitempty"`
	Links []struct {
		Title string `yaml:"title"`
		Address string `yaml:"address"`
	} `yaml:"links"`
}

type Config struct {
	ACL map[string][]string `yaml:"acl,omitempty"`
	Port int `yaml:"port"`
	Handlers struct {
		FileIndex FileindexHandlerStruct `yaml:"fileindex,omitempty"`
		Graphviz GraphvizStruct `yaml:"graphviz,omitempty"`
	} `yaml:"handlers,omitempty"`
	RootContent []CardStruct `yaml:"root_content"`
}

var config Config

func Initialize(path string) {
	confFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open configuration file: %v", err)
	}

	confRaw, err := ioutil.ReadAll(confFile)
	_ = confRaw
	if err != nil {
		log.Fatalf("Cannot read configuration file: %v", err)
	}

	err = yaml.Unmarshal(confRaw, &config)
	if err != nil {
		log.Fatalf("Invalid configuration file: %v", err)
	}
}

func Get() Config {
	return config
}

func IsAllowedByACL(addr net.IP, aclName string) bool {
	validNetStrs, ok := config.ACL[aclName]
	if !ok || len(validNetStrs) == 0 {
		return false
	}

	for _, netStr := range validNetStrs {
		_, validNet, err := net.ParseCIDR(netStr)
		if err != nil {
			validAddr := net.ParseIP(netStr)
			return validAddr.Equal(addr)
		}

		return validNet.Contains(addr)
	}

	return false
}
