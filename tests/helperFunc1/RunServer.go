package helperfunc1_test

import (
	"context"
	"net"
	"os"
	"os/exec"
	"time"
)

// StartTestServer launches your server by running `go run .` in the given `dir`,
// writing all its output to logPath (fresh each run), waits until it accepts
// connections on port, and returns a cancel func to shut it down when your test finishes.
func StartTestServer(
	dir string,
	args []string, // additional args passed to `go run .` (e.g. flags)
	logPath string,
	port string,
) (cancel func(), err error) {
	println(" we are starting the server :)")
	// Create a cancellable context
	ctx, cancelCtx := context.WithCancel(context.Background())

	// Prepare the `go run .` command in the target directory
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", "."}, args...)...)
	cmd.Dir = dir

	// Create (or truncate) the log file
	logFile, err := os.Create(logPath)
	if err != nil {
		cancelCtx()
		return nil, err
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Start the server process
	if err := cmd.Start(); err != nil {
		logFile.Close()
		cancelCtx()
		return nil, err
	}

	// Launch a goroutine that polls TCP until the server is listening
	go func() {
		for {
			conn, err := net.Dial("tcp", net.JoinHostPort("localhost", port))
			if err == nil {
				conn.Close() // successful connect means server is ready
				return       // exit the goroutine
			}
			time.Sleep(50 * time.Millisecond) // wait a bit before retrying
		}
	}()

	// Return a cancel func that kills the server process and closes the log
	return func() {
		cancelCtx()     // send kill signal to `go run`
		cmd.Wait()      // wait for process to exit and release resources
		logFile.Close() // flush and close the log file
	}, nil
}
