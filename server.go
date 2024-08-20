package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"
)

const (
    blockedFile = "/etc/pf.blocked"
    port        = "0.0.0.0:8080" // Force IPv4
)

func blockedIPsHandler(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadFile(blockedFile)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to read blocked file: %v", err), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    logConnection(r, "Disconnected")
}

func logConnection(r *http.Request, action string) {
    ip := r.RemoteAddr
    timestamp := time.Now().UTC().Format(time.RFC3339)
    fmt.Printf("[%s] %s: %s\n", timestamp, action, ip)
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
    logConnection(r, "Connected")
    blockedIPsHandler(w, r)
}

func main() {
    http.HandleFunc("/blocked_ips", connectionHandler)
    fmt.Printf("Starting server at port %s\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
