package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Runtime configuration
type Runtime struct {
	LogDir           string
	TempDir          string
	CurrentDirectory string
}

// SmsGateway configuration
type SmsGateway struct {
	Addr       string
	Port       int
	SystemID   string
	Password   string
	SystemType string
	// Alias default alias
	Alias                         string
	MaxSmsLength                  int
	SendSMSIntervalInMiliSec      int
	EnquiryIntervalInSec          int
	SeqMod                        uint
	SeqSeed                       uint
	IsDebugMode                   bool
	ShowHeartBeat                 bool
	CacheDurForProcessStatusInSec int
	NumConcurrentHandler          int
	SendSMSExpiredInMinute        int
}

// PDUAddr ...
type PDUAddr struct {
	Ton byte
	Npi byte
}

// PDUConfigs pdu configurations
type PDUConfigs struct {
	SourceAddr           *PDUAddr
	DestinationAddr      *PDUAddr
	AddRange             byte
	SubmitRespRadix      byte
	DeliverReportRadix   byte
	ProtocolID           byte
	RegisteredDelivery   byte
	ReplaceIfPresentFlag byte
	EsmClass             byte
}

// Database ...
type Database struct {
	Driver                  string
	DSN                     string
	MaxIdleConn             int
	MaxOpenConn             int
	ConnMaxLifetimeInMinute int
}

// Secure config
type Secure struct {
	// EnableTLS or not
	EnableTLS bool

	// CertFile where server certificate file is
	CertFile string

	// KeyFile where server key file is
	KeyFile string

	// SipHashSum0
	SipHashSum0 uint64

	// SipHashSum1
	SipHashSum1 uint64
}

// WebServer hold configurations for WebServer
type WebServer struct {
	// BodyLimit The body limit is determined based on both Content-Length request header and actual content read, which makes it super secure.
	// Limit can be specified as 4x or 4xB, where x is one of the multiple from K, M, G, T or P. Example: 2M = 2 Megabyte
	BodyLimit string

	// Port Which port WebServer would listen
	Port int

	// Secure
	Secure Secure
}

// SMSSenderAuth basic auth of sms sender
type SMSSenderAuth struct {
	Username string
	Password string
}

// Configuration store all configuration
type Configuration struct {
	Runtime       *Runtime
	SmsGateway    *SmsGateway
	PDUConfigs    *PDUConfigs
	Database      *Database
	WebServer     *WebServer
	SMSSenderAuth *SMSSenderAuth
}

// globalConfig hold skeleton/static configuration for whole application
var globalConfig *Configuration
var lock sync.Mutex

// GetConfigurations get current system configuration
func GetConfigurations() *Configuration {
	return globalConfig
}

// LoadConfigurations load configurations from json file
func LoadConfigurations(jsonFilePath string) error {
	// Try to open file config
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode json config
	decoder := json.NewDecoder(file)
	configure := Configuration{}
	if err := decoder.Decode(&configure); err != nil {
		return err
	}

	lock.Lock()
	defer lock.Unlock()

	// Set global config
	globalConfig = &configure

	// Try to create logs directory
	if err = os.MkdirAll(configure.Runtime.LogDir, 0700); err != nil {
		return err
	}

	// Initialize logger
	// ...

	return nil
}
