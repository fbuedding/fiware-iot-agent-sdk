package iotagentsdk

import (
	"reflect"
	"testing"
)

func TestDevice_MarshalJSON(t *testing.T) {
	type fields struct {
		Id                 DeciveId
		Service            string
		ServicePath        string
		EntityName         string
		EntityType         string
		Timezone           string
		Timestamp          *bool
		Apikey             Apikey
		Endpoint           string
		Protocol           string
		Transport          string
		Attributes         []Attribute
		Commands           []Command
		Lazy               []LazyAttribute
		StaticAttributes   []StaticAttribute
		InternalAttributes []interface{}
		ExplicitAttrs      any
		NgsiVersion        string
		PayloadType        string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test true as string",
			fields: fields{
				Id:                 "",
				Service:            "",
				ServicePath:        "",
				EntityName:         "",
				EntityType:         "",
				Timezone:           "",
				Timestamp:          nil,
				Apikey:             "",
				Endpoint:           "",
				Protocol:           "",
				Transport:          "",
				Attributes:         []Attribute{},
				Commands:           []Command{},
				Lazy:               []LazyAttribute{},
				StaticAttributes:   []StaticAttribute{},
				InternalAttributes: []interface{}{},
				ExplicitAttrs:      "true",
				NgsiVersion:        "",
				PayloadType:        "",
			},
			want:    []byte(`{"explicitAttrs":true}`),
			wantErr: false,
		},
		{
			name: "Test false as string",
			fields: fields{
				Id:                 "",
				Service:            "",
				ServicePath:        "",
				EntityName:         "",
				EntityType:         "",
				Timezone:           "",
				Timestamp:          nil,
				Apikey:             "",
				Endpoint:           "",
				Protocol:           "",
				Transport:          "",
				Attributes:         []Attribute{},
				Commands:           []Command{},
				Lazy:               []LazyAttribute{},
				StaticAttributes:   []StaticAttribute{},
				InternalAttributes: []interface{}{},
				ExplicitAttrs:      "false",
				NgsiVersion:        "",
				PayloadType:        "",
			},
			want:    []byte(`{"explicitAttrs":false}`),
			wantErr: false,
		},
		{
			name: "Test True with spaces and tabs as string",
			fields: fields{
				Id:                 "",
				Service:            "",
				ServicePath:        "",
				EntityName:         "",
				EntityType:         "",
				Timezone:           "",
				Timestamp:          nil,
				Apikey:             "",
				Endpoint:           "",
				Protocol:           "",
				Transport:          "",
				Attributes:         []Attribute{},
				Commands:           []Command{},
				Lazy:               []LazyAttribute{},
				StaticAttributes:   []StaticAttribute{},
				InternalAttributes: []interface{}{},
				ExplicitAttrs:      " True	 ",
				NgsiVersion:        "",
				PayloadType:        "",
			},
			want:    []byte(`{"explicitAttrs":true}`),
			wantErr: false,
		},
		{
			name: "test bool",
			fields: fields{
				Id:                 "",
				Service:            "",
				ServicePath:        "",
				EntityName:         "",
				EntityType:         "",
				Timezone:           "",
				Timestamp:          nil,
				Apikey:             "",
				Endpoint:           "",
				Protocol:           "",
				Transport:          "",
				Attributes:         []Attribute{},
				Commands:           []Command{},
				Lazy:               []LazyAttribute{},
				StaticAttributes:   []StaticAttribute{},
				InternalAttributes: []interface{}{},
				ExplicitAttrs:      true,
				NgsiVersion:        "",
				PayloadType:        "",
			},
			want:    []byte(`{"explicitAttrs":true}`),
			wantErr: false,
		},
		{
			name: "test a string which is not a bool",
			fields: fields{
				Id:                 "",
				Service:            "",
				ServicePath:        "",
				EntityName:         "",
				EntityType:         "",
				Timezone:           "",
				Timestamp:          nil,
				Apikey:             "",
				Endpoint:           "",
				Protocol:           "",
				Transport:          "",
				Attributes:         []Attribute{},
				Commands:           []Command{},
				Lazy:               []LazyAttribute{},
				StaticAttributes:   []StaticAttribute{},
				InternalAttributes: []interface{}{},
				ExplicitAttrs:      "test",
				NgsiVersion:        "",
				PayloadType:        "",
			},
			want:    []byte(`{"explicitAttrs":"test"}`),
			wantErr: false,
		},
		{
			name: "Test empty string",
			fields: fields{
				Id:                 "",
				Service:            "",
				ServicePath:        "",
				EntityName:         "",
				EntityType:         "",
				Timezone:           "",
				Timestamp:          nil,
				Apikey:             "",
				Endpoint:           "",
				Protocol:           "",
				Transport:          "",
				Attributes:         []Attribute{},
				Commands:           []Command{},
				Lazy:               []LazyAttribute{},
				StaticAttributes:   []StaticAttribute{},
				InternalAttributes: []interface{}{},
				ExplicitAttrs:      "",
				NgsiVersion:        "",
				PayloadType:        "",
			},
			want:    []byte(`{}`),
			wantErr: false,
		},
		{
			name: "Test wrong type",
			fields: fields{
				Id:                 "",
				Service:            "",
				ServicePath:        "",
				EntityName:         "",
				EntityType:         "",
				Timezone:           "",
				Timestamp:          nil,
				Apikey:             "",
				Endpoint:           "",
				Protocol:           "",
				Transport:          "",
				Attributes:         []Attribute{},
				Commands:           []Command{},
				Lazy:               []LazyAttribute{},
				StaticAttributes:   []StaticAttribute{},
				InternalAttributes: []interface{}{},
				ExplicitAttrs:      1,
				NgsiVersion:        "",
				PayloadType:        "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				Id:                 tt.fields.Id,
				Service:            tt.fields.Service,
				ServicePath:        tt.fields.ServicePath,
				EntityName:         tt.fields.EntityName,
				EntityType:         tt.fields.EntityType,
				Timezone:           tt.fields.Timezone,
				Timestamp:          tt.fields.Timestamp,
				Apikey:             tt.fields.Apikey,
				Endpoint:           tt.fields.Endpoint,
				Protocol:           tt.fields.Protocol,
				Transport:          tt.fields.Transport,
				Attributes:         tt.fields.Attributes,
				Commands:           tt.fields.Commands,
				Lazy:               tt.fields.Lazy,
				StaticAttributes:   tt.fields.StaticAttributes,
				InternalAttributes: tt.fields.InternalAttributes,
				ExplicitAttrs:      tt.fields.ExplicitAttrs,
				NgsiVersion:        tt.fields.NgsiVersion,
				PayloadType:        tt.fields.PayloadType,
			}
			got, err := d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Device.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Device.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
