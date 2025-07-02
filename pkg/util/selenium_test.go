package util

// This test requires the following dependencies to be installed:
// 1. github.com/tebeka/selenium: Run `go get github.com/tebeka/selenium`
// 2. ChromeDriver: Download from https://chromedriver.chromium.org/downloads and place in PATH
// 3. Chrome browser installed
//
// Before running this test, make sure to set the following environment variables:
// - DEMO_APP_GOOGLE_CLIENT_ID
// - DEMO_APP_GOOGLE_CLIENT_SECRET

import (
	"fmt"
	"testing"
	"time"

	"github.com/tebeka/selenium"
	_ "github.com/tebeka/selenium/chrome"
)

func TestSeleniumBasic(t *testing.T) {
	// Skip this test by default as it requires external dependencies
	// Remove this line when you want to run the test manually
	//t.Skip("Skip TestSeleniumOAuth - requires ChromeDriver and Google credentials")

	const (
		seleniumPath     = "selenium-server.jar"
		chromeDriverPath = "chromedriver"
		port             = 4444
	)

	var opts []selenium.ServiceOption
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		t.Fatalf("Error starting ChromeDriver service: %v\nProbably the chromedriver is not installed. Try: \"sudo apt install chromium-chromedriver\"\n", err)
	}
	defer service.Stop()

	// Configure Chrome options for CI environment
	chromeOpts := []string{
		"--headless",                   // Run in headless mode
		"--no-sandbox",                 // Required for running as root in containers
		"--disable-dev-shm-usage",      // Overcome limited resource problems
		"--disable-gpu",                // Disable GPU acceleration
		"--remote-debugging-port=9222", // Enable remote debugging
		"--disable-web-security",       // Disable web security (use with caution)
		"--disable-features=VizDisplayCompositor",
		"--user-data-dir=/tmp/chrome-test-profile", // Specify unique user data directory
	}

	// Connect to the WebDriver instance running on localhost
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"goog:chromeOptions": map[string]interface{}{
			"args": chromeOpts,
		},
	}
	// Fix: Connect to the local WebDriver server, not to example.com
	webDriverURL := fmt.Sprintf("http://localhost:%d/wd/hub", port)
	wd, err := selenium.NewRemote(caps, webDriverURL)
	if err != nil {
		t.Fatalf("Error connecting to WebDriver: %v", err)
	}
	defer wd.Quit()

	// Navigate to the app server
	if err := wd.Get("https://example.com"); err != nil {
		t.Fatalf("Error navigating to app server: %v", err)
	}

	// Wait for the page to load
	time.Sleep(2 * time.Second)

	// Find and check the header
	header, err := wd.FindElement(selenium.ByTagName, "h1")
	if err != nil {
		t.Fatalf("Error finding header: %v", err)
	}

	if hdrTxt, err := header.Text(); err != nil || hdrTxt != "Example Domain" {
		t.Fatalf("got: %+v", header)
	}
}
