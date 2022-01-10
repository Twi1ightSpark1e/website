package config

import (
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"

	"github.com/Twi1ightSpark1e/website/log"
)

type FileindexHandlerEndpointStruct struct {
	View string `yaml:"view,omitempty"`
	Auth []string `yaml:"auth,omitempty"`
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

type WebhookEndpointStruct struct {
	View string `yaml:"view,omitempty"`
	Method string `yaml:"method,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
	Exec string `yaml:"exec"`
}
type WebhookStruct struct {
	Endpoints map[string]WebhookEndpointStruct `yaml:"endpoints"`
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
	Auth map[string]string `yaml:"auth,omitempty"`
	ACL map[string][]string `yaml:"acl,omitempty"`
	Port int `yaml:"port"`
	TemplatesPath string `yaml:"templates_path"`
	Handlers struct {
		FileIndex FileindexHandlerStruct `yaml:"fileindex,omitempty"`
		Graphviz GraphvizStruct `yaml:"graphviz,omitempty"`
		Webhook WebhookStruct `yaml:"webhook,omitempty"`
	} `yaml:"handlers,omitempty"`
	RootContent []CardStruct `yaml:"root_content"`
}

var config Config
var logger = log.New("ConfigParser")

func Initialize(path string) {
	logger.Info.Printf("Using configuration file %s", path)

	confRaw, err := os.ReadFile(path)
	if err != nil {
		logger.Err.Fatalf("Cannot read configuration file: %v", err)
	}

	err = yaml.Unmarshal(confRaw, &config)
	if err != nil {
		logger.Err.Fatalf("Invalid configuration file: %v", err)
	}

	updatePasswords(path)
}

func updatePasswords(path string) {
	var oldUsers []string
	var newUsers []string

	for user, pass := range(config.Auth) {
		_, err := bcrypt.Cost([]byte(pass))
		if err == nil {
			continue
		}

		newPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			logger.Err.Fatalf("Cannot hash password of user '%s': %v", user, err)
		}
		config.Auth[user] = string(newPass)

		oldUsers = append(oldUsers, fmt.Sprintf("%s:\\s+%s", user, pass))
		newUsers = append(newUsers, fmt.Sprintf("%s: %s", user, newPass))

		logger.Info.Printf("Updated password of user '%s'", user)
	}

	if len(oldUsers) == 0 {
		return
	}

	confRaw, err := os.ReadFile(path)
	if err != nil {
		logger.Err.Fatalf("Cannot read configuration file: %v", err)
	}

	for idx := range(oldUsers) {
		old, err := regexp.Compile(oldUsers[idx])
		if err != nil {
			logger.Err.Fatalf("Cannot compile replacement regular expression: %v", err)
		}
 		confRaw = old.ReplaceAllLiteral(confRaw, []byte(newUsers[idx]))
	}

	file, err := os.OpenFile(path, os.O_WRONLY | os.O_TRUNC, 0644)
	if err != nil {
		logger.Err.Fatalf("Cannot open configuration file to write: %v", err)
	}

	file.Write(confRaw)
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
		if err == nil {
			if validNet.Contains(addr) {
				return true
			}
			continue
		}

		validAddr := net.ParseIP(netStr)
		if validAddr.Equal(addr) {
			return true
		}
	}

	return false
}

func Authenticate(r *http.Request, allowedUsers []string) bool {
	if len(allowedUsers) == 0 {
		return true
	}

	user, pass, ok := r.BasicAuth()
	if !ok {
		return false
	}

	for _, allowedUser:= range(allowedUsers) {
		if subtle.ConstantTimeCompare([]byte(user), []byte(allowedUser)) == 1 {
			hashPass := config.Auth[user]
			return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass)) == nil
		}
	}

	return false
}
