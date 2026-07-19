package data

import "testing"

func TestNewDataRejectsNilConfig(t *testing.T) {
	initializedData, err := NewData(nil)
	if err == nil {
		t.Fatal("NewData(nil) error = nil, want error")
	}
	if initializedData != nil {
		t.Fatalf("NewData(nil) data = %#v, want nil", initializedData)
	}
}

func TestNewRedisRejectsNilConfig(t *testing.T) {
	client, err := NewRedis(nil)
	if err == nil {
		t.Fatal("NewRedis(nil) error = nil, want error")
	}
	if client != nil {
		t.Fatalf("NewRedis(nil) client = %#v, want nil", client)
	}
}
