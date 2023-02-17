package emis

type Sensor struct {
	ID                           int    `json:"id"`
	InternalID                   string `json:"dsInternalId"`
	SensorTypeId                 string `json:"sensorTypeId"`
	ShortDescription             string `json:"shortDescription"`
	SerialNumber                 string `json:"serialNumber"`
	DataSupplierId               string `json:"dataSupplierId"`
	ObjectEmisCode               string `json:"objectEmisCode"`
	PermissionReadSensor         bool   `json:"permissionReadSensor"`
	PermissionWriteSensor        bool   `json:"permissionWriteSensor"`
	PermissionReadSensorReading  bool   `json:"permissionReadSensorReading"`
	PermissionWriteSensorReading bool   `json:"permissionWriteSensorReading"`
}

type SensorReading struct {
	SensorID       int     `json:"sensorId"`
	Date           string  `json:"date"`
	Value          float64 `json:"value"`
	DataSupplierID string  `json:"dataSupplierId,omitempty"`
}
