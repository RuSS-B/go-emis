package emis

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Emis struct {
	login    string
	password string
	host     string
}

func New(login string, password string, host string) *Emis {
	return &Emis{
		login:    login,
		password: password,
		host:     host,
	}
}

func (c *Emis) Request(method string, uri string, data []byte) (*http.Request, error) {
	body := bytes.NewReader(data)

	req, err := http.NewRequest(method, c.host+uri, body)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("X-Ekonerg-Login", c.login)
	req.Header.Set("X-Ekonerg-MAC", c.buildMac(method, uri, data))
	if err != nil {
		return req, fmt.Errorf("NewRequest: %w", err)
	}

	return req, nil
}

func (c *Emis) buildMac(method string, uri string, data []byte) string {
	str := fmt.Sprintf("%s\n%s\n%s\n%s%s", method, uri, c.login, c.password, data)

	h := sha256.New()
	h.Write([]byte(str))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func handleRequest(r *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("handleRequest: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("handleRequest: %w", err)
	}

	if res.StatusCode > 299 {
		return nil, errors.New(string(body))
	}

	return body, nil
}

func handleResponse[T interface{}](r *http.Request, res *T) ([]byte, error) {
	body, err := handleRequest(r)
	if err != nil {
		return nil, fmt.Errorf("handleResponse: %w", err)
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return body, fmt.Errorf("handleResponse: %w", err)
	}

	return body, nil
}

func (c *Emis) GetSensorTypes() ([]byte, error) {
	req, err := c.Request(http.MethodGet, "/em-remote-service/sensors/json/types", nil)
	if err != nil {
		return nil, fmt.Errorf("GetSensorTypes: %w", err)
	}

	body, err := handleRequest(req)
	if err != nil {
		return nil, fmt.Errorf("GetSensorTypes: %w", err)
	}

	return body, nil
}

type SensorsListResponse struct {
	EmisResponse
	Sensors []Sensor `json:"sensors"`
}

func (c *Emis) Sensors() (SensorsListResponse, error) {
	var res SensorsListResponse
	req, err := c.Request(http.MethodGet, "/em-remote-service/sensors/json/sensors", nil)
	if err != nil {
		return res, fmt.Errorf("Sensors: %w", err)
	}

	_, err = handleResponse(req, &res)
	if err != nil {
		return res, fmt.Errorf("Sensors: %w", err)
	}

	return res, nil
}

type SensorReadingsListResponse struct {
	EmisResponse
	SensorReadings []SensorReading `json:"sensorReadings"`
}

func (c *Emis) SensorReadings(ID int, year int) (SensorReadingsListResponse, error) {
	var res SensorReadingsListResponse

	req, err := c.Request(http.MethodGet, fmt.Sprintf("/em-remote-service/sensors/json/readings/%d/%d", ID, year), nil)
	if err != nil {
		return res, fmt.Errorf("SensorReadings: %w", err)
	}

	_, err = handleResponse(req, &res)

	if err != nil {
		return res, fmt.Errorf("SensorReadings: %w", err)
	}

	return res, nil
}

type SendSensorReadingsResponse struct {
	EmisResponse
	SuccessfulEntries int `json:"successfulEntries"`
	FailedEntries     int `json:"failedEntries"`
}

func (c *Emis) SendSensorReadings(reading SensorReading) (SendSensorReadingsResponse, error) {
	var res SendSensorReadingsResponse

	data, _ := json.Marshal(reading)
	req, err := c.Request(http.MethodPost, "/em-remote-service/sensors/json/readings", data)
	if err != nil {
		return res, fmt.Errorf("SensorReadings: %w", err)
	}
	_, err = handleResponse(req, &res)

	if err != nil {
		return res, fmt.Errorf("SensorReadings: %w", err)
	}

	return res, nil
}
