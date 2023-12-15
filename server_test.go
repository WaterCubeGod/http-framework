package main

import (
	"fmt"
	"testing"
)

func TestHTTP_Start(t *testing.T) {
	h := NewHTTP(WithHTTPServerStop(nil))
	go func() {
		err := h.Start(":8080")
		if err != nil {
			fmt.Println("HTTP.Start err:", err)
			t.Fail()
		}
	}()
	err := h.Stop()
	if err != nil {
		t.Fail()
	}
}
