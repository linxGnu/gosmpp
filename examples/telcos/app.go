package main

import (
	"github.com/linxGnu/gosmpp/examples/telcos/config"
	"github.com/linxGnu/gosmpp/examples/telcos/controller/SmppSession"
	"github.com/linxGnu/gosmpp/examples/telcos/controller/SmsSender"
	"github.com/linxGnu/gosmpp/examples/telcos/daemons"
	"github.com/linxGnu/gosmpp/examples/telcos/dao"

	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// initializeRouter initialize router and http server
func initializeRouter(config *config.Configuration) (e *echo.Echo, server *http.Server) {
	// Initialize router
	e = echo.New()

	// Initialize http server
	server = &http.Server{
		Addr: ":" + strconv.Itoa(config.WebServer.Port),
	}

	// Check TLS enable
	if config.WebServer.Secure.EnableTLS {
		// HTTPS redirect middleware redirects http requests to https. For more example, please refer to https://echo.labstack.com/
		e.Pre(mw.HTTPSRedirect())

		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		// Load certificate files
		if cer, err := tls.LoadX509KeyPair(config.WebServer.Secure.CertFile, config.WebServer.Secure.KeyFile); err != nil {
			panic(err)
		} else {
			tlsConfig.Certificates = []tls.Certificate{cer}
		}

		// Enable TLS
		server.TLSConfig = tlsConfig
	}

	// Recover middleware recovers from panics anywhere in the chain, prints stack trace and handles the control to the centralized HTTPErrorHandler.
	// Default stack size for trace is 4KB. For more example, please refer to https://echo.labstack.com/
	e.Use(mw.Recover())

	// Remove trailing slash middleware removes a trailing slash from the request URI.
	e.Pre(mw.RemoveTrailingSlash())

	// Set BodyLimit Middleware. It will panic if fail. For more example, please refer to https://echo.labstack.com/
	e.Use(mw.BodyLimit(config.WebServer.BodyLimit))

	// Secure middleware provides protection against cross-site scripting (XSS) attack, content type sniffing, clickjacking, insecure connection and other code injection attacks.
	// For more example, please refer to https://echo.labstack.com/
	e.Use(mw.Secure())

	grSender := e.Group("/sender")

	// BasicAuth for napas
	grSender.Use(mw.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == config.SMSSenderAuth.Username && password == config.SMSSenderAuth.Password {
			return true, nil
		}

		return false, nil
	}))

	grSender.POST("/sendSMS", SmsSender.SendSMS)

	return
}

// changeWorkingDir change current working directory to directory holding
// current executing file
func changeWorkingDir() (currentDir string, err error) {
	currentDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	err = os.Chdir(currentDir)
	return
}

func main() {
	// change to current directory and load configurations
	var configs *config.Configuration
	if curDir, err := changeWorkingDir(); err != nil {
		panic(err)
	} else {
		if err = config.LoadConfigurations("configs.json"); err != nil {
			panic(err)
		}
		configs = config.GetConfigurations()
		configs.Runtime.CurrentDirectory = curDir
	}

	// try to connect to database first
	if err := dao.BindDAO(configs.Database); err != nil {
		panic(fmt.Errorf("Connect to database error: %v", err.Error()))
	}

	// initialize smpp session
	SmppSession.InitializeSmppSessionInstance(configs.SmsGateway, configs.PDUConfigs)
	if err := SmppSession.Instance.Bind(); err != nil {
		panic(fmt.Errorf("Binding to SMSC error: %v", err.Error()))
	}

	// Now run daemons
	daemons.RunDaemons()

	// Now initialize router and server
	router, server := initializeRouter(configs)

	// Try to start server, if something wrong, it would panic to stderr
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := router.StartServer(server); err != nil {
			panic(err)
		}
	}()

	// Receive signal interrupt
	<-quit

	// Stop daemons
	daemons.StopDaemons()

	// Do unbind
	SmppSession.Instance.Destroy()

	// wait no longer than 5 seconds to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := router.Shutdown(ctx); err != nil {
		panic(err)
	}
}
