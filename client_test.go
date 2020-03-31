package api

import (
	"os"
	"testing"
)

func TestCheck(t *testing.T) {
	c, err := NewClient(os.Getenv("baseurl"), os.Getenv("token"))

	if err != nil {
		t.Errorf("Error reached : err %v", err)
	}

	resp, err := c.Check(&CheckRequest{
		Lon: 107.5743943,
		Lat: -6.8781377,
	})

	if err != nil {
		t.Errorf("Error when checking covid data %v", err)
	}

	if resp.Kelurahan != "Sarijadi" {
		t.Errorf("Kelurahan mismatch , got %s want: Sarijadi", resp.Kelurahan)
	}
}
