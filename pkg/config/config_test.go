package config

import (
	"fmt"
	"os"
	"testing"
)

func TestPrintConfigFormat(t *testing.T) {
	dir, _ := os.Getwd()
	fmt.Printf("%s\n", dir)

	data, err := PrintConfigFormat()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s\n", data)
}

func TestNew(t *testing.T) {
	dir, _ := os.Getwd()
	conf, err := New(dir + "/conf.yaml")
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	fmt.Printf("%#v\n", conf)
}
