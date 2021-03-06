// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package logs

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (LogConfig) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func (LogConfig) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in LogConfig and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg LogConfig) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("LogConfig", pflag.ExitOnError)
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "cloudwatch-enabled"), *new(bool), "Enable Cloudwatch Logging")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "cloudwatch-region"), *new(string), "AWS region in which Cloudwatch logs are stored.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "cloudwatch-log-group"), *new(string), "Log group to which streams are associated.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "cloudwatch-template-uri"), *new(string), "Template Uri to use when building cloudwatch log links")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "kubernetes-enabled"), *new(bool), "Enable Kubernetes Logging")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "kubernetes-url"), *new(string), "Console URL for Kubernetes logs")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "kubernetes-template-uri"), *new(string), "Template Uri to use when building kubernetes log links")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "stackdriver-enabled"), *new(bool), "Enable Log-links to stackdriver")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "gcp-project"), *new(string), "Name of the project in GCP")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "stackdriver-logresourcename"), *new(string), "Name of the logresource in stackdriver")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "stackdriver-template-uri"), *new(string), "Template Uri to use when building stackdriver log links")
	return cmdFlags
}
