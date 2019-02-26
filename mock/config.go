package mock

import (
	"log"
	"net"
	"net/http"
	"path/filepath"

	"vuvuzela.io/alpenhorn/config"
)

// LaunchConfigServer creates a local config server, a config HTTP server, and a
// config client.
func LaunchConfigServer(dir string) (*config.Server, *config.Client) {
	configServer, err := config.CreateServer(filepath.Join(dir, "config-server-state"))
	if err != nil {
		log.Panicf("config.CreateServer: %s", err)
	}
	configListener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Panic(err)
	}
	configHTTPServer := &http.Server{
		Handler: configServer,
	}
	go func() {
		err := configHTTPServer.Serve(configListener)
		if err != http.ErrServerClosed {
			log.Fatalf("http.Serve: %s", err)
		}
	}()
	configClient := &config.Client{
		ConfigServerURL: "http://" + configListener.Addr().String(),
	}
	return configServer, configClient
}
