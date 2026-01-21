package luhn

import (
	"errors"
	"testing"
)

func TestValidateValid(t *testing.T) {
	valid, err := Validate([]byte("20000006"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valid {
		t.Fatal("expected valid luhn number")
	}
}

func TestValidateInvalidNumber(t *testing.T) {
	valid, err := Validate([]byte("495599"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid {
		t.Fatal("expected invalid luhn number")
	}
}

func TestValidateInvalidChar(t *testing.T) {
	_, err := Validate([]byte("g495599"))
	if err == nil {
		t.Fatal("expected error for invalid character")
	}
	if !errors.Is(err, ErrInvalidChar) {
		t.Fatalf("expected ErrInvalidChar, got %v", err)
	}
}

func TestValidateNoDigits(t *testing.T) {
	_, err := Validate([]byte("--"))
	if err == nil {
		t.Fatal("expected error for no digits")
	}
	if !errors.Is(err, ErrNoDigits) {
		t.Fatalf("expected ErrNoDigits, got %v", err)
	}
}

func TestSumLuhnNoDigits(t *testing.T) {
	_, _, err := sumLuhn([]byte("   "), true)
	if err == nil {
		t.Fatal("expected error for no digits")
	}
	if !errors.Is(err, ErrNoDigits) {
		t.Fatalf("expected ErrNoDigits, got %v", err)
	}
}
