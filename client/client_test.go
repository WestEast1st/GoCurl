package client

import (
	"fmt"
	"testing"
)

var url = "http://example.com"

func TestClientGet(t *testing.T) {
	client := New(url)
	response, err := client.Get()
	if err != nil {
		t.Errorf("Could not access %s .", url)
	}
	fmt.Println(response)
}
