package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Host     string
	UserID   string
	DeviceID string
}

type Client struct {
	host              string
	userID            string
	deviceID          string
	payload           []byte
	stopHeartbeatChan chan struct{}
}

func New(config Config) *Client {
	return &Client{
		host:     config.Host,
		userID:   config.UserID,
		deviceID: config.DeviceID,
	}
}

func (c *Client) Connect() {
	code, err := c.sendRequest("/connect")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	if code == http.StatusOK {
		fmt.Println("Connected successfully")
		c.initHeartbeat()
	} else {
		fmt.Println("Failed to connect")
		return
	}
}

func (c *Client) Disconnect() {
	code, err := c.sendRequest("/disconnect")
	if err != nil {
		fmt.Println("Error disconnecting:", err)
		return
	}
	if code == http.StatusOK {
		fmt.Println("Disconnected successfully")
		c.stopHeartbeat()
	} else {
		fmt.Println("Failed to disconnect")
		return
	}
}

func (c *Client) heartbeat() {
	code, err := c.sendRequest("/heartbeat")
	if err != nil {
		fmt.Println("Error sending heartbeat:", err)
		return
	}
	if code == http.StatusOK {
		fmt.Println("Heartbeat sent successfully")
	} else {
		fmt.Println("Failed to send heartbeat")
		return
	}
}

func (c *Client) initHeartbeat() {
	c.stopHeartbeatChan = make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.heartbeat()
			case <-c.stopHeartbeatChan:
				return
			}
		}
	}()
}

func (c *Client) stopHeartbeat() {
	if c.stopHeartbeatChan != nil {
		close(c.stopHeartbeatChan)
	}
}

func (c *Client) getPayload() ([]byte, error) {
	if c.payload != nil {
		return c.payload, nil
	}

	var err error

	payloadMap := map[string]string{
		"user_id":   c.userID,
		"device_id": c.deviceID,
	}
	// Convert the payload to JSON
	c.payload, err = json.Marshal(payloadMap)

	return c.payload, err
}

func (c *Client) getRequest(path string) (*http.Request, error) {
	payload, err := c.getPayload()
	if err != nil {
		// log error
		fmt.Println("Error getting payload:", err)
		return nil, err
	}

	url := c.host + path

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) sendRequest(path string) (int, error) {
	req, err := c.getRequest(path)
	if err != nil {
		return 500, err
	}

	// Create an HTTP client and send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return 500, err
	}
	defer resp.Body.Close()

	fmt.Printf("Request sent to %s\n", path)
	fmt.Printf("Request method: %s\n", req.Method)

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response body: %s\n", string(body))

	// Handle the response
	fmt.Printf("Response code: %d\n", resp.StatusCode)
	return resp.StatusCode, nil
}
