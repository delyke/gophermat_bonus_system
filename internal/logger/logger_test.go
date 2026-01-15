package logger

import "testing"

func TestSetLevelBeforeInit(t *testing.T) {
	SetLevel("debug")
}

func TestSetNopLogger(t *testing.T) {
	SetNopLogger()

	if Logger() == nil {
		t.Fatal("expected logger to be set")
	}
	if err := Sync(); err != nil {
		t.Fatalf("unexpected sync error: %v", err)
	}
}

func TestInitLogger(t *testing.T) {
	if err := Init("debug", true); err != nil {
		t.Fatalf("init logger error: %v", err)
	}
	if Logger() == nil {
		t.Fatal("expected logger to be set")
	}
}
