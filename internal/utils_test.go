package internal

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault_Bool_Invalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The function did not panic")
		}
	}()

	_ = GetEnvOrDefault("INVALID", "invalid", "bool")
	// If we reach this point, the function did not panic
	t.Error("Expected panic, function returned value")
}

func TestGetEnvOrDefault_Bool_Valid(t *testing.T) {
	err := os.Setenv("TEST_BOOL", "false")
	if err != nil {
		t.Error("Failed to set environment variable")
	}
	test := GetEnvOrDefault("TEST_BOOL", "true", "bool")
	if test != false {
		t.Error("Expected false, got true")
	}
}

func TestGetEnvOrDefault_Bool_Default(t *testing.T) {
	test := GetEnvOrDefault("INVALID", "true", "bool")
	if test != true {
		t.Error("Expected true, got false")
	}
}

func TestGetEnvOrDefault_Int_Invalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The function did not panic")
		}
	}()

	_ = GetEnvOrDefault("INVALID", "invalid", "int")
	// If we reach this point, the function did not panic
	t.Error("Expected panic, function returned value")
}

func TestGetEnvOrDefault_Int_Valid(t *testing.T) {
	err := os.Setenv("TEST_INT", "1")
	if err != nil {
		t.Error("Failed to set environment variable")
	}
	test := GetEnvOrDefault("TEST_INT", "0", "int")
	if test != 1 {
		t.Error("Expected 1, got 0")
	}
}

func TestGetEnvOrDefault_Int_Default(t *testing.T) {
	test := GetEnvOrDefault("INVALID", "1", "int")
	if test != 1 {
		t.Error("Expected 1, got 0")
	}
}

func TestGetEnvOrDefault_String_Valid(t *testing.T) {
	err := os.Setenv("TEST_STRING", "test")
	if err != nil {
		t.Error("Failed to set environment variable")
	}
	test := GetEnvOrDefault("TEST_STRING", "default", "string")
	if test != "test" {
		t.Error("Expected test, got default")
	}
}

func TestGetEnvOrDefault_String_Default(t *testing.T) {
	test := GetEnvOrDefault("INVALID", "default", "string")
	if test != "default" {
		t.Error("Expected default, got invalid")
	}
}

func TestGetEnvOrDefault_Slice_Valid(t *testing.T) {
	err := os.Setenv("TEST_SLICE", "test1,test2,test3")
	if err != nil {
		t.Error("Failed to set environment variable")
	}
	test := GetEnvOrDefault("TEST_SLICE", "", "[]string")
	if len(test.([]string)) != 3 {
		t.Error("Expected 3, got", len(test.([]string)))
	}
	if test.([]string)[0] != "test1" {
		t.Error("Expected test1, got", test.([]string)[0])
	}
	if test.([]string)[1] != "test2" {
		t.Error("Expected test2, got", test.([]string)[1])
	}
	if test.([]string)[2] != "test3" {
		t.Error("Expected test3, got", test.([]string)[2])
	}
}

func TestGetEnvOrDefault_Slice_Default(t *testing.T) {
	test := GetEnvOrDefault("INVALID", "", "[]string")
	if len(test.([]string)) != 0 {
		t.Error("Expected 0, got", len(test.([]string)))
	}
}
