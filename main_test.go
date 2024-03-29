package iotagentsdk_test

import (
	"os"
	"testing"

	i "github.com/fbuedding/fiware-iot-agent-sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	iota i.IoTA
	fs   i.FiwareService
	d    i.Device
	sg   i.ConfigGroup
)

const (
	deviceId          = i.DeciveId("test_device")
	entityName        = "TestEntityName"
	updatedEntityName = "TestEntityNameUpdated"
	service           = "testing"
	servicePath       = "/"
	resource          = i.Resource("/iot/d")
	apiKey            = "testKey"
)

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	host := "localhost"
	if os.Getenv("TEST_HOST") != "" {
		host = os.Getenv("TEST_HOST")
	}
	log.Info().Msgf("Starting test with iot-agent host: %s", host)
	iota = *i.NewIoTAgent(host, 4061, 1000)
	fs = i.FiwareService{Service: service, ServicePath: servicePath}
	d = i.Device{
		Id:                 deviceId,
		Service:            service,
		ServicePath:        servicePath,
		EntityName:         entityName,
		EntityType:         "",
		Timezone:           "",
		Timestamp:          new(bool),
		Apikey:             apiKey,
		Endpoint:           "",
		Protocol:           "",
		Transport:          "",
		Attributes:         []i.Attribute{},
		Commands:           []i.Command{},
		Lazy:               []i.LazyAttribute{},
		StaticAttributes:   []i.StaticAttribute{},
		InternalAttributes: []interface{}{},
		ExplicitAttrs:      true,
		NgsiVersion:        "",
		PayloadType:        "",
	}
	sg = i.ConfigGroup{
		Service:                      service,
		ServicePath:                  servicePath,
		Resource:                     resource,
		Apikey:                       apiKey,
		Timestamp:                    new(bool),
		EntityType:                   "Test",
		Trust:                        "",
		CbHost:                       host,
		Lazy:                         []i.LazyAttribute{},
		Commands:                     []i.Command{},
		Attributes:                   []i.Attribute{},
		StaticAttributes:             []i.StaticAttribute{{Name: "test", Type: "Number", Value: 6, Metadata: map[string]i.Metadata{}}},
		InternalAttributes:           []interface{}{},
		ExplicitAttrs:                "",
		EntityNameExp:                entityName,
		NgsiVersion:                  "",
		DefaultEntityNameConjunction: entityName,
		Autoprovision:                false,
		PayloadType:                  "",
		Transport:                    "",
		Endpoint:                     "",
	}
	iota.DeleteDevice(fs, d.Id)
	err := iota.CreateDevice(fs, d)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create device for tests")
	}
	iota.DeleteConfigGroup(fs, resource, apiKey)
	err = iota.CreateConfigGroup(fs, sg)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not delete config group for tests")
	}
	m.Run()
	teardown()
}

func teardown() {
	err := iota.DeleteDevice(fs, d.Id)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not delete device for teardown")
	}

	err = iota.DeleteConfigGroup(fs, resource, apiKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not delete device for teardown")
	}
}
