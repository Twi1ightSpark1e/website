package config

import (
	"io/ioutil"
	"net"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/Twi1ightSpark1e/website/log"
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
	TemplatesPath string `yaml:"templates_path"`
	Handlers struct {
		FileIndex FileindexHandlerStruct `yaml:"fileindex,omitempty"`
		Graphviz GraphvizStruct `yaml:"graphviz,omitempty"`
	} `yaml:"handlers,omitempty"`
	RootContent []CardStruct `yaml:"root_content"`
}

var config Config
var logger = log.New("ConfigParser")

func Initialize(path string) {
	logger.Info.Printf("Using configuration file %s", path)

	confFile, err := os.Open(path)
	if err != nil {
		logger.Err.Fatalf("Cannot open configuration file: %v", err)
	}

	confRaw, err := ioutil.ReadAll(confFile)
	_ = confRaw
	if err != nil {
		logger.Err.Fatalf("Cannot read configuration file: %v", err)
	}

	err = yaml.Unmarshal(confRaw, &config)
	if err != nil {
		logger.Err.Fatalf("Invalid configuration file: %v", err)
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
