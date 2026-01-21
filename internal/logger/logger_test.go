package logger

import "testing"

func TestSetLevelBeforeInit(t *testing.T) {
	l := NewNop()
	l.SetLevel("debug")
}

func TestSetNopLogger(t *testing.T) {
	l := NewNop()

	if l == nil {
		t.Fatal("expected logger to be created")
	}
	if err := l.Sync(); err != nil {
		t.Fatalf("unexpected sync error: %v", err)
	}
}

func TestInitLogger(t *testing.T) {
	l, err := New("debug", true)
	if err != nil {
		t.Fatalf("init logger error: %v", err)
	}
	if l == nil {
		t.Fatal("expected logger to be set")
	}
}
