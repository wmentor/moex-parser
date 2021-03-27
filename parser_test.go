package parser

import (
	"testing"
)

func TestParser(t *testing.T) {

	_, _, err := Get()
	if err != nil {
		t.Fatal(err)
	}
}
