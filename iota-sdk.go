// Package iotagentsdk provides an SDK for communicating with Fiware IoT Agents.
package iotagentsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

// Constants for the IoT Agent URL and Healthcheck URL.
const (
	urlBase        = "http://%v:%d"
	urlHealthcheck = urlBase + "/iot/about"
)

// Error returns the error as a formatted string.
func (e ApiError) Error() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Message)
}

// init initializes the logging level based on the "LOG_LEVEL" environment variable.
// By default, "panic" is used if no environment variable is set.
func init() {
	logLvl := os.Getenv("LOG_LEVEL")
	if logLvl == "" {
		logLvl = "panic"
	}
	SetLogLevel(logLvl)
}

// NewIoTAgent creates a new instance of the IoT Agent.
func NewIoTAgent(host string, port int, timeout_ms int) *IoTA {
	iota := IoTA{
		Host:       host,
		Port:       port,
		timeout_ms: time.Duration(timeout_ms) * time.Millisecond,
		client:     &http.Client{Timeout: time.Duration(timeout_ms) * time.Millisecond},
	}
	return &iota
}

// SetLogLevel sets the global logging level based on the given value.
// Possible values are: "trace", "debug", "info", "warning", "error", "fatal", "panic".
func SetLogLevel(ll string) {
	ll = strings.ToLower(ll)
	switch ll {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		log.Fatal().Msg("Log level need to be one of this: [TRACE DEBUG INFO WARNING ERROR FATAL PANIC]")
	}
}

// Healthcheck performs a health check of the IoT Agent and returns the result.
func (i IoTA) Healthcheck() (*RespHealthcheck, error) {
	response, err := http.Get(fmt.Sprintf(urlHealthcheck, i.Host, i.Port))
	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}
	var respHealth RespHealthcheck
	json.Unmarshal(responseData, &respHealth)
	if respHealth.LibVersion == "" {
		return nil, fmt.Errorf("Error healtchecking IoT-Agent, host: %s", i.Host)
	}
	log.Debug().Str("Response healthcheck", string(responseData)).Any("Healthcheck", respHealth).Send()
	return &respHealth, nil
}

// GetAllServicePathsForService returns all service paths for the specified service.
func (i IoTA) GetAllServicePathsForService(service string) ([]string, error) {
	cgs, err := i.ListConfigGroups(FiwareService{service, "/*"})
	if err != nil {
		return nil, err
	}

	if cgs.Count == 0 {
		return nil, nil
	}

	servicePaths := []string{}
	for _, cg := range cgs.Services {
		if !slices.Contains(servicePaths, cg.ServicePath) {
			servicePaths = append(servicePaths, cg.ServicePath)
		}
	}

	return servicePaths, nil
}

// Client returns the HTTP client used for communication with the IoT Agent.
// If no client is present, a new one is created.
func (i IoTA) Client() *http.Client {
	if i.client == nil {
		log.Debug().Msg("Creating http client")
		i.client = &http.Client{Timeout: i.timeout_ms}
	}
	return i.client
}
