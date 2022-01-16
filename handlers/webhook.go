package handlers

import (
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/google/shlex"
)

type webhookHandler struct {
	logger log.Channels
	path string
	endpoint config.WebhookEndpointStruct
}
func WebhookHandler(logger log.Channels, path string, endpoint config.WebhookEndpointStruct) http.Handler {
	return &webhookHandler{logger, path, endpoint}
}

func (h *webhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	if !assertPath(h.path, w, r, h.logger.Err) {
		return
	}

	if h.endpoint.Method != "" && strings.ToLower(r.Method) != strings.ToLower(h.endpoint.Method) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for header, value := range h.endpoint.Headers {
		if r.Header[header][0] == value {
			continue
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	cmdline, _ := shlex.Split(h.endpoint.Exec)
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	args := strings.Join(cmdline[1:], " ")
	if cmd == nil {
		h.logger.Err.Printf("Cannot create process '%s' with arguments '%s'", cmdline[0], args)
		return
	}

	cmd.Stdout = h.logger.Info.Writer()
	cmd.Stderr = h.logger.Err.Writer()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		h.logger.Err.Printf("Cannot get stdin of spawned process: %v", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		h.logger.Err.Printf("Cannot spawn process: %v", err)
		return
	}
	h.logger.Info.Printf("Spawned process '%s' with arguments '%s'", cmdline[0], args)

	_, err = io.Copy(stdin, r.Body)
	if err != nil {
		h.logger.Err.Printf("Cannot send request body to process stdin: %v", err)
		return
	}
	stdin.Close()
}
