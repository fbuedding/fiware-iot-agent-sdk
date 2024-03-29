package iotagentsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	u "net/url"

	"github.com/niemeyer/golang/src/pkg/container/vector"
	log "github.com/rs/zerolog/log"
)

// Constants
const (
	urlDevice = urlBase + "/iot/devices"
)

// Request struct for creating a device
type reqCreateDevice struct {
	Devices []Device `json:"devices"`
}

// Response struct for listing devices
type respListDevices struct {
	Count   int      `json:"count"`
	Devices []Device `json:"devices"`
}

// Function to validate a Device
func (d Device) Validate() error {
	mF := &MissingFields{make(vector.StringVector, 0), "Missing fields"}
	if d.Id == "" {
		mF.Fields.Push("Id")
	}

	if mF.Fields.Len() == 0 {
		return nil
	} else {
		return mF
	}
}

// Method to read a device
func (i IoTA) ReadDevice(fs FiwareService, id DeciveId) (*Device, error) {
	url, err := u.JoinPath(fmt.Sprintf(urlDevice, i.Host, i.Port), u.PathEscape(string(id)))
	if err != nil {
		return nil, err
	}
	method := "GET"

	client := i.Client()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		err = json.Unmarshal(resData, &apiError)
		if err != nil {
			log.Panic().Err(err).Msg("Could not Marshal struct")
		}
		return nil, apiError
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}

	var device Device
	json.Unmarshal(responseData, &device)
	if err != nil {
		log.Panic().Err(err).Msg("Could not Marshal struct")
	}
	return &device, nil
}

// Method to check if a device exists
func (i IoTA) DeviceExists(fs FiwareService, id DeciveId) bool {
	_, err := i.ReadDevice(fs, id)
	if err != nil {
		return false
	}
	return true
}

// Method to list devices
func (i IoTA) ListDevices(fs FiwareService) (*respListDevices, error) {
	url := fmt.Sprintf(urlDevice, i.Host, i.Port)

	method := "GET"

	client := i.Client()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		err = json.Unmarshal(resData, &apiError)
		if err != nil {
			log.Panic().Err(err).Msg("Could not Marshal struct")
		}
		return nil, apiError
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}

	var respDevices respListDevices
	json.Unmarshal(responseData, &respDevices)
	return &respDevices, nil
}

// Method to create a device
func (i IoTA) CreateDevices(fs FiwareService, ds []Device) error {
	for _, sg := range ds {
		err := sg.Validate()
		if err != nil {
			return err
		}
	}
	rcd := reqCreateDevice{}
	rcd.Devices = ds[:]
	method := "POST"

	payload, err := json.Marshal(rcd)
	if err != nil {
		log.Panic().Err(err).Msg("Could not Marshal struct")
	}
	client := i.Client()
	req, err := http.NewRequest(method, fmt.Sprintf(urlDevice, i.Host, i.Port), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("Error while creating Request %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while requesting resource %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return apiError
	}

	return nil
}

// Method to create a device
func (i IoTA) CreateDevice(fs FiwareService, d Device) error {
	ds := [1]Device{d}
	return i.CreateDevices(fs, ds[:])
}

// Method to update a device
func (i IoTA) UpdateDevice(fs FiwareService, d Device) error {
	err := d.Validate()
	if err != nil {
		return err
	}

	url, err := u.JoinPath(fmt.Sprintf(urlDevice, i.Host, i.Port), u.PathEscape(string(d.Id)))

	// Ensure these fields are not set
	d.Id = ""
	d.Transport = ""
	d.Service = ""
	d.ServicePath = ""

	method := "PUT"

	payload, err := json.Marshal(d)
	if err != nil {
		log.Panic().Err(err).Msg("Could not Marshal struct")
	}
	if string(payload) == "{}" {
		return nil
	}
	client := i.Client()
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("Error while creating Request %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while requesting resource %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return apiError
	}
	return nil
}

// Method to delete a device
func (i IoTA) DeleteDevice(fs FiwareService, id DeciveId) error {
	url, err := u.JoinPath(fmt.Sprintf(urlDevice, i.Host, i.Port), u.PathEscape(string(id)))
	if err != nil {
		return err
	}
	method := "DELETE"

	client := i.Client()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return fmt.Errorf("Error while getting service: %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while getting service: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return apiError
	}
	return nil
}

// Method to upsert a device
func (i IoTA) UpsertDevice(fs FiwareService, d Device) error {
	exists := i.DeviceExists(fs, d.Id)
	if !exists {
		log.Debug().Msg("Creating device...")
		err := i.CreateDevice(fs, d)
		if err != nil {
			return err
		}
	} else {
		log.Debug().Msg("Update device...")
		dTmp, err := i.ReadDevice(fs, d.Id)
		if err != nil {
			return err
		}

		if dTmp.EntityName == "" {
			return errors.New("Error before getting updating device: No entity_name")
		}

		d.Transport = ""
		d.EntityName = dTmp.EntityName
		err = i.UpdateDevice(fs, d)
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates a device an updates the given Device
func (i IoTA) CreateDeviceWSE(fs FiwareService, d *Device) error {
	if d == nil {
		return errors.New("Device reference cannot be nil")
	}
	err := i.CreateDevice(fs, *d)
	if err != nil {
		return err
	}
	dTmp, err := i.ReadDevice(fs, d.Id)
	if err != nil {
		return err
	}
	*d = *dTmp
	return nil
}
