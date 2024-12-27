package lib

import (
	"io"
	"log"
	"testing"
)

func TestRecover(t *testing.T) {
	log.SetOutput(io.Discard)
	go func() {
		defer Recover()
		panic("error")
	}()
}
