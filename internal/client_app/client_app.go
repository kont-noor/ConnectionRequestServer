package clientapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Run() {
	fmt.Println("Init client app")
	time.Sleep(4 * time.Second)
	host := "http://localhost:3000" // Replace with your server URL

	payload := map[string]string{
		"user_id":   "1001",
		"device_id": "2002",
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	sendRequest(host+"/connect", jsonPayload)

	for {
		sendRequest(host+"/heartbeat", jsonPayload)
		time.Sleep(4 * time.Second)
	}
}

func sendRequest(url string, payload []byte) {
	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{Timeout: 50 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request successful")
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}
}
