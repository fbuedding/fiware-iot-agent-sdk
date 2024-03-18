package iotagentsdk_test

import (
	"os"
	"testing"

	i "github.com/fbuedding/fiware-iot-agent-sdk"
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
	host := "localhost"
	if os.Getenv("TEST_HOST") != "" {
		host = os.Getenv("TEST_HOST")
	}
	log.Info().Msgf("Starting test with iot-agent host: %s", host)
	iota = i.IoTA{Host: host, Port: 4061}
	fs = i.FiwareService{Service: service, ServicePath: servicePath}
	d = i.Device{Id: deviceId, EntityName: entityName}
	sg = i.ConfigGroup{
		Service:       service,
		ServicePath:   servicePath,
		Resource:      resource,
		Apikey:        apiKey,
		EntityType:    "Test",
		Autoprovision: false,
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
