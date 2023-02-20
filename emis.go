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

func (c *Emis) SensorTypes() ([]byte, error) {
	req, err := c.Request(http.MethodGet, "/em-remote-service/sensors/json/types", nil)
	if err != nil {
		return nil, fmt.Errorf("SensorTypes: %w", err)
	}

	body, err := handleRequest(req)
	if err != nil {
		return nil, fmt.Errorf("SensorTypes: %w", err)
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
	SensorReadings []SensorReading `json:"readings"`
}

func (c *Emis) SensorReadings(ID int, year int) (SensorReadingsListResponse, error) {
	var res SensorReadingsListResponse

	req, err := c.Request(http.MethodGet, fmt.Sprintf("/em-remote-service/sensors/json/readings/%d/%d", ID, year), nil)
	if err != nil {
		return res, fmt.Errorf("SensorReadings: %w", err)
	}

	body, err := handleResponse(req, &res)
	res.RawBody = body

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

type Meter struct {
	MeterID               int64     `json:"meterId"`
	MeterSerialNumber     string    `json:"meterSerialNumber"`
	MeterShortDescription string    `json:"meterShortDescription"`
	MeterDescription      string    `json:"meterDescription"`
	ObjectName            string    `json:"objectName"`
	IsgeCode              string    `json:"isgeCode"`
	Energent              string    `json:"energent"`
	EnergentId            string    `json:"energentId"`
	Counters              []Counter `json:"counters"`
}

type Counter struct {
	Counter int    `json:"counter"`
	Name    string `json:"name"`
}

type MeterListResponse struct {
	EmisResponse
	Meters []Meter `json:"meters"`
}

func (c *Emis) Meters() (MeterListResponse, error) {
	var res MeterListResponse
	req, err := c.Request(http.MethodGet, "/em-remote-service/query/json/meters", nil)
	if err != nil {
		return res, fmt.Errorf("Meters: %w", err)
	}

	_, err = handleResponse(req, &res)
	if err != nil {
		return res, fmt.Errorf("Meters: %w", err)
	}

	return res, nil
}

type MeterReading struct {
	ServiceUrl string  `json:"serviceUrl,omitempty"`
	ID         string  `json:"id,omitempty"`
	MeterID    string  `json:"meterId"`
	Date       string  `json:"date"`
	C1         float64 `json:"c1"`
	C2         float64 `json:"c2,omitempty"`
	C3         float64 `json:"c3,omitempty"`
	C4         float64 `json:"c4,omitempty"`
	C5         float64 `json:"c5,omitempty"`
	InternalID string  `json:"internalId"`
}

type MeterReadingsListResponse struct {
	EmisResponse
	Readings []MeterReading `json:"readings"`
}

func (c *Emis) MeterReadings(ID int, year int, month int) (MeterReadingsListResponse, error) {
	var res MeterReadingsListResponse

	req, err := c.Request(http.MethodGet, fmt.Sprintf("/em-remote-service/query/json/meter/%d/readings/%d/%d", ID, year, month), nil)
	if err != nil {
		return res, fmt.Errorf("MeterReadings: %w", err)
	}

	body, err := handleResponse(req, &res)
	res.RawBody = body

	if err != nil {
		return res, fmt.Errorf("MeterReadings: %w", err)
	}

	return res, nil
}

type InsertReadingRequest[T SuccessReading | FailedReading | MeterReading] struct {
	Insert Reading[T] `json:"insert"`
}

type Reading[T SuccessReading | FailedReading | MeterReading] struct {
	Readings []T `json:"readings"`
}

type SuccessReading struct {
	TS         string `json:"ts"`
	ID         string `json:"id"`
	InternalID string `json:"internalId"`
}

type FailedReading struct {
	TS      string  `json:"ts"`
	Err     string  `json:"err"`
	ErrDesc string  `json:"errDesc"`
	MeterID string  `json:"meterId"`
	Date    string  `json:"date"`
	C1      float64 `json:"c1"`
}

type SendMeterReadingsResponse struct {
	EmisResponse
	Succeeded InsertReadingRequest[SuccessReading] `json:"succeeded"`
	Failed    InsertReadingRequest[FailedReading]  `json:"failed"`
}

func (c *Emis) SendMeterReadings(r []MeterReading) (SendMeterReadingsResponse, error) {
	reqData := InsertReadingRequest[MeterReading]{
		Insert: Reading[MeterReading]{
			Readings: r,
		},
	}

	var res SendMeterReadingsResponse

	data, _ := json.Marshal(reqData)
	req, err := c.Request(http.MethodPost, "/em-remote-service/batch/json/send", data)
	if err != nil {
		return res, fmt.Errorf("SendMeterReadings: %w", err)
	}
	_, err = handleResponse(req, &res)

	if err != nil {
		return res, fmt.Errorf("SendMeterReadings: %w", err)
	}

	return res, nil
}
