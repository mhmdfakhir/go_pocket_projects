package pocketlog_test

import (
	"learngo/logger/pocketlog"
	"testing"
)

func ExampleLogger_Debugf() {
	debugLogger := pocketlog.New(pocketlog.LevelDebug)
	debugLogger.Debugf("Hello, %s", "world")
	// Output: [DEBUG] Hello, world
}

func ExampleLogger_Infof() {
	debugLogger := pocketlog.New(pocketlog.LevelInfo)
	debugLogger.Infof("Hello, %s", "world")
	// Output: [INFO] Hello, world
}

func ExampleLogger_Errorf() {
	debugLogger := pocketlog.New(pocketlog.LevelError)
	debugLogger.Errorf("Hello, %s", "world")
	// Output: [ERROR] Hello, world
}

// testWriter is a struct that implements io.Writer.
// We use it to validate that we can write to a specific output.
type testWriter struct {
	contents string
}

// Write implements the io.Writer interface.
func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.contents = tw.contents + string(p)
	return len(p), nil
}

const (
	debugMessage = "Why write I still all one, ever the same,"
	infoMessage  = "And keep invention in a noted weed,"
	errorMessage = "That every word doth almost tell my name,"
)

func TestLogger_DebugfInfofErrorf(t *testing.T) {
	type testCase struct {
		level    pocketlog.Level
		expected string
	}
	tt := map[string]testCase{
		"debug": {
			level:    pocketlog.LevelDebug,
			expected: "[DEBUG] Why write I still all one, ever the same,\n[INFO] And keep invention in a noted weed,\n[ERROR] That every word doth almost tell my name,\n",
		},
		"info": {
			level:    pocketlog.LevelInfo,
			expected: "[INFO] And keep invention in a noted weed,\n[ERROR] That every word doth almost tell my name,\n",
		},
		"error": {
			level:    pocketlog.LevelError,
			expected: "[ERROR] That every word doth almost tell my name,\n",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tw := &testWriter{}
			testedLogger := pocketlog.New(tc.level, pocketlog.WithOutput(tw))
			testedLogger.Debugf(debugMessage)
			testedLogger.Infof(infoMessage)
			testedLogger.Errorf(errorMessage)
			if tw.contents != tc.expected {
				t.Errorf("invalid contents, expected %q, got %q",
					tc.expected, tw.contents)
			}
		})
	}
}
