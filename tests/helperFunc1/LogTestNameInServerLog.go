package helperfunc1_test

import (
	"fmt"
	"os"
	"testing"
)

func LogTestNameInTheServerLogFile(t *testing.T) {
	logPath := "server.log"
	// Open (or create) the log file in append-only mode
	f, err := os.OpenFile(logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // :contentReference[oaicite:0]{index=0}
	if err != nil {
		t.Logf("unable to open log file %s: %v", logPath, err)
		return
	}
	// Write the "start" banner
	fmt.Fprintf(f, "\n\n======================== %s ====================\n\n", t.Name()) // :contentReference[oaicite:1]{index=1}
	f.Close()

	// Register cleanup to write the "end" banner
	t.Cleanup(func() { // :contentReference[oaicite:2]{index=2}
		f2, err := os.OpenFile(logPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			t.Logf("unable to open log file %s: %v", logPath, err)
			return
		}
		fmt.Fprintf(f2, "---- %s ends ----\n", t.Name())
		f2.Close()
	})
}
