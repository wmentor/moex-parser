package parser

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {

	rh, err := os.Open("securities.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer rh.Close()

	parse(rh)
}
