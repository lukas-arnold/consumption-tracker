package utils

import (
	"testing"
	"time"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestConvertConsumptionStorageToBytes(t *testing.T) {
	bytes, err := ConvertConsumptionStorageToBytes(models.ConsumptionStorage{})
	if err != nil {
		t.Fatal(err)
	}
	if len(bytes) == 0 {
		t.Fatal("expected json bytes")
	}
}

func TestConvertBytesToConsumptionStorage(t *testing.T) {
	_, err := ConvertBytesToConsumptionStorage([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestConvertBytesToConsumptionStorageInvalid(t *testing.T) {
	_, err := ConvertBytesToConsumptionStorage([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConvertId(t *testing.T) {
	got, err := ConvertId("123")
	if err != nil || got != 123 {
		t.Fatalf("got %d, want 123 (err: %v)", got, err)
	}
}

func TestConvertIdInvalid(t *testing.T) {
	if _, err := ConvertId("abc"); err == nil {
		t.Fatal("expected error")
	}
}

func TestConvertFloat(t *testing.T) {
	got, err := ConvertFloat("42.5")
	if err != nil || got != 42.5 {
		t.Fatalf("got %f, want 42.5 (err: %v)", got, err)
	}
}

func TestConvertFloatInvalid(t *testing.T) {
	if _, err := ConvertFloat("abc"); err == nil {
		t.Fatal("expected error")
	}
}

func TestConvertInt(t *testing.T) {
	got, err := ConvertInt("123")
	if err != nil || got != 123 {
		t.Fatalf("got %d, want 123 (err: %v)", got, err)
	}
}

func TestConvertIntInvalid(t *testing.T) {
	if _, err := ConvertInt("abc"); err == nil {
		t.Fatal("expected error")
	}
}

func TestConvertTime(t *testing.T) {
	got, err := ConvertTime("2024-01-15")
	if err != nil {
		t.Fatal(err)
	}

	expected := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Fatalf("got %v, want %v", got, expected)
	}
}

func TestConvertTimeInvalid(t *testing.T) {
	if _, err := ConvertTime("15-01-2024"); err == nil {
		t.Fatal("expected error")
	}
}
