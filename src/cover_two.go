package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os/exec"
    "strings"
    "time"
)

// Constants for the server URL and file paths
const (
    serverURL     = "http://localhost:8080/blocked_ips" // Replace with your server URL
    blockedFile   = "/etc/pf.blocked.example" // Replace with your pf.blocked file
    pfConfFile    = "/etc/pf.conf" // Replace with your pf.conf file
    checkInterval = 1 * time.Minute // Check every 1 minute
    httpTimeout   = 10 * time.Second // Timeout for HTTP requests
)

// FetchBlockedIPs fetches the blocked IPs from the server
func FetchBlockedIPs() (string, error) {
    client := &http.Client{Timeout: httpTimeout}
    resp, err := client.Get(serverURL)
    if err != nil {
        return "", fmt.Errorf("error fetching blocked IPs: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %v", err)
    }

    return string(body), nil
}

// UpdateBlockedFile updates the blocked IPs file
func UpdateBlockedFile(data string) error {
    return ioutil.WriteFile(blockedFile, []byte(data), 0644)
}

// RunPFScript checks if PF is enabled, enables it if not, and applies the rules
func RunPFScript() error {
    // Check if PF is already enabled
    cmd := exec.Command("pfctl", "-s", "info")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("error checking pf status: %v", err)
    }

    if !contains(output, "Status: Enabled") {
        err := exec.Command("sudo", "pfctl", "-e").Run()
        if err != nil {
            return fmt.Errorf("error enabling pf: %v", err)
        }
    }

    commands := []string{
        fmt.Sprintf("pfctl -f %s", pfConfFile),
        fmt.Sprintf("pfctl -t blocked -T replace -f %s", blockedFile),
    }

    for _, cmd := range commands {
        err := exec.Command("sudo", "sh", "-c", cmd).Run()
        if err != nil {
            return fmt.Errorf("error running command '%s': %v", cmd, err)
        }
    }

    return nil
}

// Helper function to check if the PF status output contains a specific substring
func contains(output []byte, substring string) bool {
    return strings.Contains(string(output), substring)
}

func logWithTimestamp(message string) {
    log.Printf("[%s] %s\n", time.Now().UTC().Format(time.RFC3339), message)
}

func main() {
    for {
        logWithTimestamp("Starting check-in...")

        // Fetch the blocked IPs from the server
        ips, err := FetchBlockedIPs()
        if err != nil {
            logWithTimestamp(fmt.Sprintf("Error: %v", err))
            time.Sleep(checkInterval)
            continue
        }
        logWithTimestamp("Fetched blocked IPs successfully.")

        // Update the blocked IPs file
        if err := UpdateBlockedFile(ips); err != nil {
            logWithTimestamp(fmt.Sprintf("Error updating blocked file: %v", err))
            time.Sleep(checkInterval)
            continue
        }
        logWithTimestamp("Updated blocked file successfully.")

        // Run the pf script directly from Go
        if err := RunPFScript(); err != nil {
            logWithTimestamp(fmt.Sprintf("Error running pf script: %v", err))
            time.Sleep(checkInterval)
            continue
        }
        logWithTimestamp("Ran check-in  successfully.")

        // Wait for the next check-in with a countdown timer
        logWithTimestamp("Check-in completed.")
        for i := int(checkInterval.Seconds()); i > 0; i-- {
            if i == 1 {
                fmt.Printf("\rNext check-in in %d second", i)
            } else {
                fmt.Printf("\rNext check-in in %d seconds", i)
            }
            time.Sleep(1 * time.Second)
        }
        fmt.Println() // Move to the next line after the countdown
    }
}

