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
    blockedFile   = "/etc/pf.blocked"
    pfScript      = "/usr/local/bin/reload_pf.sh"
    checkInterval = 1 * time.Minute // Check every 1 minute
    httpTimeout   = 10 * time.Second // Timeout for HTTP requests
)

// FetchBlockedIPs fetches the blocked IPs from the server
func FetchBlockedIPs() (string, error) {
    client := &http.Client{Timeout: httpTimeout}
"pf_updater.go" 99L, 2944B
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
    blockedFile   = "/etc/pf.blocked"
    pfScript      = "/usr/local/bin/reload_pf.sh"
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
    blockedFile   = "/etc/pf.blocked"
    pfScript      = "/usr/local/bin/reload_pf.sh"
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

// RunPFScript runs the pf scriptpackage main

import (
    "fmt"
    "net/http"
    "os/exec"
    "time"
)

// Constants for the server URL and file paths
const (
    blockedFile   = "/etc/pf.blocked"
    pfScript      = "/usr/local/bin/reload_pf.sh"
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
    if err != nil {
        return "", fmt.Errorf("error reading response body: %v", err)
    }

    return string(body), nil
}

// UpdateBlockedFile updates the blocked IPs file
func UpdateBlockedFile(data string) error {
    return ioutil.WriteFile(blockedFile, []byte(data), 0644)
}

// RunPFScript runs the pf script
func RunPFScript() error {
    cmd := exec.Command("sudo", pfScript)
    output, err := cmd.CombinedOutput()
    if err != nil {
    }
    return nil
}

func logWithTimestamp(message string) {
    log.Printf("[%s] %s\n", time.Now().UTC().Format(time.RFC3339), message)
}

package main

import (
    "fmt"
    "net/http"
    "os/exec"
    "time"
)

// Constants for the server URL and file paths
const (
    blockedFile   = "/etc/pf.blocked"
    pfScript      = "/usr/local/bin/reload_pf.sh"
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
    if err != nil {
        return "", fmt.Errorf("error reading response body: %v", err)
    }

    return string(body), nil
}

// UpdateBlockedFile updates the blocked IPs file
func UpdateBlockedFile(data string) error {
    return ioutil.WriteFile(blockedFile, []byte(data), 0644)
}

// RunPFScript runs the pf script
func RunPFScript() error {
    cmd := exec.Command("sudo", pfScript)
    output, err := cmd.CombinedOutput()
    if err != nil {
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

        // Run the pf script
        if err := RunPFScript(); err != nil {
            logWithTimestamp(fmt.Sprintf("Error running pf script: %v", err))
            time.Sleep(checkInterval)
            continue
        }
        logWithTimestamp("Ran pf script successfully.")

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
