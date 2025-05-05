package helperfunc1_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
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
	logChan <-chan string,
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
	// go func() {
	// 	for {
	// 		conn, err := net.Dial("tcp", net.JoinHostPort("localhost", port))
	// 		if err == nil {
	// 			conn.Close() // successful connect means server is ready
	// 			return       // exit the goroutine
	// 		}
	// 		time.Sleep(50 * time.Millisecond) // wait a bit before retrying
	// 	}
	// }()

	serverReady := make(chan struct{}) // Signal channel for server readiness
	var wg sync.WaitGroup              // WaitGroup to ensure logging goroutine finishes

	// --- Goroutine to poll for server readiness ---
	go func() {
		defer close(serverReady) // Close channel when done polling (success or context cancel)
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done(): // If the main context is cancelled before server starts
				fmt.Fprintf(os.Stderr, "Server readiness check cancelled: %v\n", ctx.Err())
				return // Exit the goroutine
			case <-ticker.C:
				conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", port), 2*time.Second) // Added timeout
				if err == nil {
					conn.Close()
					// successful connect means server is ready
					println("Server connection successful on port", port)
					return // exit the goroutine
				}
				// Optional: Log the connection error for debugging
				// fmt.Fprintf(os.Stderr, "Polling server on port %s failed: %v\n", port, err)
			}
		}
	}()

	// --- Goroutine to listen on logChan and write to logFile ---
	wg.Add(1)

	go func() {
		defer wg.Done() // Signal that this goroutine has finished
		println("Log writing goroutine started for server.log")
		for logMsg := range logChan { // Loop continues until logChan is closed
			// formattedMsg := fmt.Sprintf("[TEST LOG] %s\n", logMsg)
			_, err := logFile.WriteString(logMsg)
			if err != nil {
				// Log error to stderr, as writing to logFile failed
				fmt.Fprintf(os.Stderr, "ERROR writing to server log file %s: %v\n", logPath, err)
				// Decide if you want to break or continue trying
			}
		}
		println("Log writing goroutine finished for server.log (channel closed)")
	}()

	// Wait until the server is ready before returning
	select {
	case <-serverReady:
		// Server is ready, continue
		println("Server reported as ready.")
	case <-time.After(30 * time.Second): // Add a timeout for server readiness
		cancelCtx()     // Cancel the context
		cmd.Wait()      // Wait for the command to exit
		logFile.Close() // Close the log file
		wg.Wait()       // Ensure logging goroutine finishes *after* file is potentially closed by failed start
		return nil, fmt.Errorf("server failed to become ready on port %s within 30 seconds", port)
	case <-ctx.Done(): // Handle case where context is cancelled externally during startup wait
		cmd.Wait()
		logFile.Close()
		wg.Wait()
		return nil, fmt.Errorf("server startup cancelled: %w", ctx.Err())
	}

	// Return a cancel func that kills the server process and closes the log
	return func() {
		println("the server is shutting down as the stopping func is called (form defer hopefully)")
		cancelCtx()       // send kill signal to `go run`
		err := cmd.Wait() // wait for process to exit and release resources
		if err != nil {
			println("there is a error in closing the server and the error is ", err.Error())
		}
		// Wait for the logging goroutine to finish processing any remaining messages
		// This happens *after* the server process has exited but *before* closing the file.
		// The logging goroutine finishes when logChan is closed (which happens in TestMain's defer).
		println("Waiting for log writing goroutine to finish...")
		wg.Wait()
		println("Log writing goroutine finished.")

		err = logFile.Close() // flush and close the log file
		if err != nil {
			println("there is a error in closing the server.log file", err.Error())
		}

		println("Server shutdown complete.")
	}, nil
}
