package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"crypto/rand"

	"github.com/Sirupsen/logrus"
	"path/filepath"
)

var logger *logrus.Logger

func run(args []string) int {
	bindAddress := flag.String("ip", "0.0.0.0", "IP address to bind")
	listenPort := flag.Int("port", 25478, "port number to listen on")
	// 5,242,880 bytes == 5 MiB
	maxUploadSize := flag.Int64("upload_limit", 5242880, "max size of uploaded file (byte)")
	tokenFlag := flag.String("token", "", "specify the security token (it is automatically generated if empty)")
	getTokenFlag := flag.String("gettoken", "", "specify the security token for accessing the resource only (it is automatically generated if empty)")
	logLevelFlag := flag.String("loglevel", "info", "logging level")
	flag.Parse()
	serverRoot := flag.Arg(0)
	if len(serverRoot) == 0 {
		flag.Usage()
		return 2
	}
	if logLevel, err := logrus.ParseLevel(*logLevelFlag); err != nil {
		logrus.WithError(err).Error("failed to parse logging level, so set to default")
	} else {
		logger.Level = logLevel
	}
	token := *tokenFlag
	getToken := *getTokenFlag
	if token == "" {
		count := 10
		b := make([]byte, count)
		if _, err := rand.Read(b); err != nil {
			logger.WithError(err).Fatal("could not generate token")
			return 1
		}
		token = fmt.Sprintf("%x", b)
		logger.WithField("token", token).Warn("token generated")
	}
	if getToken == "" {
		count := 10
		b := make([]byte, count)
		if _, err := rand.Read(b); err != nil {
			logger.WithError(err).Fatal("could not generate get-token")
			return 1
		}
		getToken = fmt.Sprintf("%x", b)
		logger.WithField("get-token", getToken).Warn("get-token generated")
	}
	logger.WithFields(logrus.Fields{
		"ip":           *bindAddress,
		"port":         *listenPort,
		"token":        token,
		"get-token":    getToken,
		"upload_limit": *maxUploadSize,
		"root":         serverRoot,
	}).Info("start listening")
	fullPath, _ := filepath.Abs(serverRoot)
	logger.Infof("Current Storage Path: %s", fullPath)
	server := NewServer(serverRoot, *maxUploadSize, token, getToken)
	http.Handle("/upload", server)
	http.Handle("/files/", server)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *bindAddress, *listenPort), nil)
	return 0
}

func main() {
	logger = logrus.New()
	logger.Info("starting up simple-upload-server")

	result := run(os.Args)
	os.Exit(result)
}
