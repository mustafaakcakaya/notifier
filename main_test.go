package main

import (
	"testing"
	"time"
)

func TestNewBuySellNotifier(t *testing.T) {
	_, err := NewBuySellNotifier()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestGetFormattedNow(t *testing.T) {
	bsn := &BuySellNotifier{}
	expectedFormat := "2006-01-02 15:04:05" // The layout you use in the GetFormattedNow method

	now := time.Now()
	formattedNow := bsn.GetFormattedNow()

	expectedFormattedNow := now.Format(expectedFormat)

	if formattedNow != expectedFormattedNow {
		t.Errorf("Expected formatted time %s, but got %s", expectedFormattedNow, formattedNow)
	}
}
