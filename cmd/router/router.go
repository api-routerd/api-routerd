// SPDX-License-Identifier: Apache-2.0

package router

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/api-routerd/api-routerd/cmd/container"
	"github.com/api-routerd/api-routerd/cmd/network"
	"github.com/api-routerd/api-routerd/cmd/proc"
	"github.com/api-routerd/api-routerd/cmd/share"
	"github.com/api-routerd/api-routerd/cmd/system"
	"github.com/api-routerd/api-routerd/cmd/systemd"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// StartRouter Init and start Gorilla mux router
func StartRouter(ip string, port string, tlsCertPath string, tlsKeyPath string) error {
	var srv http.Server

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	// Register services
	container.RegisterRouterContainer(s)
	network.RegisterRouterNetwork(s)
	proc.RegisterRouterProc(s)

	systemd.InitSystemd()
	systemd.RegisterRouterSystemd(s)

	system.RegisterRouterSystem(s)

	// Authenticate users
	amw, err := InitAuthMiddleware()
	if err != nil {
		log.Fatalf("Faild to init auth DB existing: %s", err)
		return fmt.Errorf("Failed to init Auth DB: %s", err)
	}

	r.Use(amw.AuthMiddleware)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop

		log.Printf("Received signal: %+v", sig)
		log.Println("Shutting down api-routerd ...")

		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Errorf("Failed to shutdown server gracefuly: %s", err)
		}

		os.Exit(0)
	}()

	// socket activation
	listeners, err := activation.Listeners()
	if err != nil {
		log.Infof("Failed to retrieve listeners: %s", err)
	}

	if share.PathExists(tlsCertPath) && share.PathExists(tlsKeyPath) {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: false,
		}
		srv = http.Server{
			Addr:         ip + ":" + port,
			Handler:      r,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		log.Info("Starting api-routerd in TLS mode")

		if len(listeners) <= 0 {
			log.Fatal(srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath))
		} else {
			log.Fatal(srv.ServeTLS(listeners[0], tlsCertPath, tlsKeyPath))
		}
	} else {
		srv = http.Server{
			Addr:    ip + ":" + port,
			Handler: r,
		}
		log.Info("Starting api-routerd in plain text mode")

		if len(listeners) <= 0 {
			log.Fatal(srv.ListenAndServe())
		} else {
			log.Fatal(srv.Serve(listeners[0]))
		}
	}

	return nil
}
