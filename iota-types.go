package iotagentsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/niemeyer/golang/src/pkg/container/vector"
)

// IoTA represents an IoT Agent instance.
type IoTA struct {
	Host       string
	Port       int
	timeout_ms time.Duration
	client     *http.Client
}

// FiwareService represents a Fiware service and its associated service path.
type FiwareService struct {
	Service     string
	ServicePath string
}

// RespHealthcheck represents the response of a health check request.
type RespHealthcheck struct {
	LibVersion string `json:"libVersion"`
	Port       string `json:"port"`
	BaseRoot   string `json:"baseRoot"`
	Version    string `json:"version"`
}

// ApiError represents an error in an API call.
type ApiError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// Attribute represents an attribute in the data model.
type Attribute struct {
	ObjectID   string              `json:"object_id,omitempty" form:"object_id"`
	Name       string              `json:"name" form:"name"`
	Type       string              `json:"type" form:"type"`
	Expression string              `json:"expression,omitempty" form:"expression"`
	SkipValue  string              `json:"skipValue,omitempty" form:"skipValue"`
	EntityName string              `json:"entity_name,omitempty" form:"entity_name"`
	EntityType string              `json:"entity_type,omitempty" form:"entity_type"`
	Metadata   map[string]Metadata `json:"metadata,omitempty" form:"metadata"`
}

// LazyAttribute represents a lazy attribute in the data model.
type LazyAttribute struct {
	ObjectID string              `json:"object_id,omitempty" form:"object_id"`
	Name     string              `json:"name" form:"name"`
	Type     string              `json:"type" form:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty" form:"metadata"`
}

// StaticAttribute represents a static attribute in the data model.
type StaticAttribute struct {
	ObjectID string              `json:"object_id,omitempty" form:"object_id"`
	Name     string              `json:"name" form:"name"`
	Type     string              `json:"type" form:"type"`
	Value    any                 `json:"value" form:"value"`
	Metadata map[string]Metadata `json:"metadata,omitempty" form:"metadata"`
}

// Command represents a command in the data model.
type Command struct {
	ObjectID    string              `json:"object_id,omitempty" form:"object_id"`
	Name        string              `json:"name" form:"name"`
	Type        string              `json:"type" form:"type"`
	Expression  string              `json:"expression,omitempty" form:"expression"`
	PayloadType string              `json:"payloadType,omitempty" form:"payloadType"`
	ContentType string              `json:"contentType,omitempty" form:"contentType"`
	Metadata    map[string]Metadata `json:"metadata,omitempty" form:"metadata"`
}

// Metadata represents metadata for attributes and commands.
type Metadata struct {
	Type  string `json:"type" form:"type"`
	Value string `json:"value" form:"value"`
}

// Apikey represents an API key.
type Apikey string

// Resource represents a resource.
type Resource string

// ConfigGroup represents a configuration group.
// See datamodel [Config Group]: https://iotagent-node-lib.readthedocs.io/en/latest/api.html#service-group-datamodel
type ConfigGroup struct {
	Service                      string            `json:"service,omitempty" form:"service"`
	ServicePath                  string            `json:"subservice,omitempty" form:"subservice"`
	Resource                     Resource          `json:"resource" form:"resource"`
	Apikey                       Apikey            `json:"apikey" form:"apikey"`
	Timestamp                    *bool             `json:"timestamp,omitempty" form:"timestamp"`
	EntityType                   string            `json:"entity_type,omitempty" form:"entity_type"`
	Trust                        string            `json:"trust,omitempty" form:"trust"`
	CbHost                       string            `json:"cbHost,omitempty" form:"cbHost"`
	Lazy                         []LazyAttribute   `json:"lazy,omitempty" form:"lazy"`
	Commands                     []Command         `json:"commands,omitempty" form:"commands"`
	Attributes                   []Attribute       `json:"attributes,omitempty" form:"attributes"`
	StaticAttributes             []StaticAttribute `json:"static_attributes,omitempty" form:"static_attributes"`
	InternalAttributes           []interface{}     `json:"internal_attributes,omitempty" form:"internal_attributes"`
	ExplicitAttrs                string            `json:"explicitAttrs,omitempty" form:"explicitAttrs"`
	EntityNameExp                string            `json:"entityNameExp,omitempty" form:"entityNameExp"`
	NgsiVersion                  string            `json:"ngsiVersion,omitempty" form:"ngsiVersion"`
	DefaultEntityNameConjunction string            `json:"defaultEntityNameConjunction,omitempty" form:"defaultEntityNameConjunction"`
	Autoprovision                bool              `json:"autoprovision,omitempty" form:"autoprovision"`
	PayloadType                  string            `json:"payloadType,omitempty" form:"payloadType"`
	Transport                    string            `json:"transport,omitempty" form:"transport"`
	Endpoint                     string            `json:"endpoint,omitempty" form:"endpoint"`
}

// DeciveId represents a device ID.
type DeciveId string

// Device represents a device.
// See datamodel [Device]: https://iotagent-node-lib.readthedocs.io/en/3.3.0/api.html#device-datamodel
type Device struct {
	Id                 DeciveId          `json:"device_id,omitempty" form:"device_id"`
	Service            string            `json:"service,omitempty" form:"service"`
	ServicePath        string            `json:"service_path,omitempty" form:"service_path"`
	EntityName         string            `json:"entity_name,omitempty" form:"entity_name"`
	EntityType         string            `json:"entity_type,omitempty" form:"entity_type"`
	Timezone           string            `json:"timezon,omitempty" form:"timezone"`
	Timestamp          *bool             `json:"timestamp,omitempty" form:"timestamp"`
	Apikey             Apikey            `json:"apikey,omitempty" form:"apikey"`
	Endpoint           string            `json:"endpoint,omitempty" form:"endpoint"`
	Protocol           string            `json:"protocol,omitempty" form:"protocol"`
	Transport          string            `json:"transport,omitempty" form:"transport"`
	Attributes         []Attribute       `json:"attributes,omitempty" form:"attributes"`
	Commands           []Command         `json:"commands,omitempty" form:"commands"`
	Lazy               []LazyAttribute   `json:"lazy,omitempty" form:"lazy"`
	StaticAttributes   []StaticAttribute `json:"static_attributes,omitempty" form:"static_attributes"`
	InternalAttributes []interface{}     `json:"internal_attributes,omitempty" form:"internal_attributes"`
	ExplicitAttrs      any               `json:"explicitAttrs,omitempty" form:"explicitAttrs"`
	NgsiVersion        string            `json:"ngsiVersion,omitempty" form:"ngsiVersion"`
	PayloadType        string            `json:"payloadType,omitempty" form:"payloadType"`
}

func (d *Device) MarshalJSON() ([]byte, error) {
	type Alias Device
	switch v := d.ExplicitAttrs.(type) {
	case string:
		tmp := strings.ToLower(v)
		tmp = strings.TrimSpace(tmp)
		tmp = strings.Trim(tmp, "\t")
		if tmp == "true" || tmp == "false" {
			return json.Marshal(&struct {
				ExplicitAttrs bool `json:"explicitAttrs" form:"explicitAttrs"`
				*Alias
			}{
				ExplicitAttrs: tmp == "true",
				Alias:         (*Alias)(d),
			})
		}

		return json.Marshal(&struct {
			ExplicitAttrs string `json:"explicitAttrs,omitempty" form:"explicitAttrs"`
			*Alias
		}{
			ExplicitAttrs: v,
			Alias:         (*Alias)(d),
		})
	case bool:
		return json.Marshal(&struct {
			ExplicitAttrs bool `json:"explicitAttrs" form:"explicitAttrs"`
			*Alias
		}{
			ExplicitAttrs: v,
			Alias:         (*Alias)(d),
		})

	default:
		return nil, fmt.Errorf("ExplicitAttrs must be a string or a bool")
	}
}

type MissingFields struct {
	Fields  vector.StringVector
	Message string
}
