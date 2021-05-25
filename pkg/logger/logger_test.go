package logger

import "testing"

func TestInit(t *testing.T) {
	Init(
		"test.log",
		"test",
		128,
		30,
		30,
		-1,
		true,
		false,
	)
	Debugf("test: %s", "abc")
	Infof("test2: %s", "def")
}
