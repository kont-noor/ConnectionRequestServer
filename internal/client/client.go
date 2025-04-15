package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	Host     string
	UserID   string
	DeviceID string
	Log      *zap.Logger
}

type Client struct {
	host              string
	userID            string
	deviceID          string
	payload           []byte
	stopHeartbeatChan chan struct{}
	log               *zap.Logger
}

func New(config Config) *Client {
	return &Client{
		host:     config.Host,
		userID:   config.UserID,
		deviceID: config.DeviceID,
		log:      config.Log,
	}
}

func (c *Client) Connect() error {
	code, err := c.sendRequest("/connect")
	if err != nil {
		c.log.Sugar().Errorf("Error connecting: %v\n", err)
		return errors.New("failed to connect")
	}
	switch code {
	case http.StatusOK:
		c.log.Info("Connected successfully")
		c.initHeartbeat()
		return nil
	case http.StatusConflict:
		c.log.Warn("Connection exists at another device")
		return errors.New("failed to connect")
	default:
		c.log.Error("Failed to connect")
		return errors.New("failed to connect")
	}
}

func (c *Client) Disconnect() error {
	code, err := c.sendRequest("/disconnect")
	if err != nil {
		c.log.Sugar().Errorf("Error disconnecting: %v\n", err)
		return errors.New("failed to disconnect")
	}
	if code == http.StatusOK {
		c.log.Info("Disconnected successfully")
		c.stopHeartbeat()
		return nil
	} else {
		c.log.Error("Failed to disconnect")
		return errors.New("failed to disconnect")
	}
}

func (c *Client) heartbeat() {
	code, err := c.sendRequest("/heartbeat")
	if err != nil {
		c.log.Sugar().Errorf("Error sending heartbeat: %v\n", err)
		return
	}
	if code == http.StatusOK {
		c.log.Info("Heartbeat sent successfully")
	} else {
		c.log.Error("Failed to send heartbeat")
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
		c.log.Sugar().Errorf("Error getting payload: %v\n", err)
		return nil, err
	}

	url := c.host + path

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		c.log.Sugar().Errorf("Error creating request: %v\n", err)
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
		c.log.Sugar().Errorf("Error sending request: %v\n", err)
		return 500, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	c.log.Debug("Request sent",
		zap.String("path", path),
		zap.String("method", req.Method),
		zap.String("response_body", string(body)),
		zap.Int("status_code", resp.StatusCode),
	)
	return resp.StatusCode, nil
}
