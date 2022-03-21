package container

import (
	"password-caddy/api/src/lib/sesclient"
	"reflect"
	"testing"
)

func TestContainerGetConfig(t *testing.T) {
	actual := TestGetConfig()
	expected := "FOOBAR"

	if actual != expected {
		t.Errorf("FAILED - TestContainerGetConfig | Actual: %s | Expected: %s", actual, expected)
	}
}

func TestLoadAwsConfig(t *testing.T) {
	actual := LoadAwsConfig()

	if actual.Region != "us-east-2" {
		t.Errorf("FAILED - TestLoadAwsConfig | Actual: %s | Expected us-east-2", actual.Region)
	}
}

func TestSesClient(t *testing.T) {
	client := SesClient()
	actual := reflect.TypeOf(*client).Kind()
	expected := reflect.TypeOf(sesclient.SesClient{}).Kind()

	if actual != expected {
		t.Errorf("FAILED - TestSesClient | Actual: %+v | Expected: %+v", actual, expected)
	}
}
