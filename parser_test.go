package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {

	shares, trades, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(shares)
	fmt.Println(trades)

	for _, v := range trades {
		data, _ := json.Marshal(v)
		fmt.Println(string(data))
	}
}
