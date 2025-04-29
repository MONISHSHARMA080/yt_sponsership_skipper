package helperfunc1_test

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/razorpay/razorpay-go"
)

func StartLocalTunnel() error {
	localPort := "8080"

	// Run the lt command to create a tunnel
	cmd := exec.Command("lt", "-p", localPort, "-s", "yt-sponsor-skipper")

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting localtunnel: %v\n", err)
		return err
	}

	fmt.Println("Started localtunnel, waiting for URL...")

	// Create a scanner to read the output
	scanner := bufio.NewScanner(stdout)

	// Scan the output for the URL
	var tunnelURL string
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		// Check if the line contains the URL indicator
		if strings.Contains(line, "your url is:") {
			// Extract the URL part
			parts := strings.Split(line, "your url is:")
			if len(parts) > 1 {
				tunnelURL = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if tunnelURL != "" {
		fmt.Printf("Tunnel established! URL: %s\n", tunnelURL)
		fmt.Println("Press Ctrl+C to stop the tunnel")

		// Wait for the command to finish (will block until the tunnel is closed)
		// if err := cmd.Wait(); err != nil {
		// 	fmt.Printf("Tunnel closed with error: %v\n", err)
		// 	return err
		// }
	} else {
		fmt.Println("Could not find tunnel URL in the output")
		// Kill the process if we couldn't find the URL
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("Failed to kill process: %v\n", err)
			return err
		}
	}
	if !strings.Contains(tunnelURL, "yt-sponsor-skipper") {
		return fmt.Errorf("the tunnel url does not contain the  yt-sponsor-skipper, in the url , your webhook will not work for this one")
	}
	return nil
}

// note the url to receive webhook on should be a valid url and also contain the path /webHookIntegrationForPaymentCapture , else we will give an error
func CreateRazorpayHook(urlToRecieveWebHookOn string) error {
	err := validateTheUrl(urlToRecieveWebHookOn)
	if err != nil {
		return err
	}
	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET_ID"))

	webhookData := map[string]interface{}{
		"url":    urlToRecieveWebHookOn,
		"active": true,
		"events": []string{
			"payment.captured",
		},
	}
	// println("the razorpay account id form env is ->", os.Getenv("RAZORPAY_ACCOUNT_ID"), " and the razorpay key id is ->", os.Getenv("RAZORPAY_KEY_ID"), " and the razorpay secret id is ->", os.Getenv("RAZORPAY_SECRET_ID"))
	// os.Getenv("RAZORPAY_ACCOUNT_ID")
	webHook, err := client.Webhook.Create("", webhookData, nil)
	if err != nil {
		println("there is a error in creating the webhook and it is ->", err.Error())
		return err
	}
	fmt.Printf("webhook created successfully and the webhook id is ->%+v", webHook)
	return nil
}

func deleteRazorpayHook(razorpayWebhookID string) error {
	return nil
}

func validateTheUrl(urlStringToValidate string) error {
	// Validate the URL format
	_, err := url.Parse(urlStringToValidate)
	if err != nil {
		return err
	}
	if !strings.Contains(urlStringToValidate, "webHookIntegrationForPaymentCapture") {
		return fmt.Errorf("the url:" + urlStringToValidate + " does not contain the path /webHookIntegrationForPaymentCapture")
	}
	return nil
}
