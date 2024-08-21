package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os/exec"
    "time"
)

// Constants for the server URL and file paths
const (
    serverURL     = "http://localhost:8080/blocked_ips" // Replace with your server URL
    blockedFile   = "/etc/pf.blocked.example" // Replace with your actual pf.blocked file
    pfConfFile    = "/etc/pf.conf.example" // Replace with your actual pf.conf file
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

// RunPFScript runs the equivalent of the original bash script
func RunPFScript() error {
    commands := []string{
        "pfctl -e",
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

        // Run the pf commands directly from Go
        if err := RunPFScript(); err != nil {
            logWithTimestamp(fmt.Sprintf("Error running pf script: %v", err))
            time.Sleep(checkInterval)
            continue
        }
        logWithTimestamp("Ran check-in successfully.")

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

