package emis

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var c *Emis

func TestMain(t *testing.M) {
	envFile := "./.env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Println("Unable to load file", envFile)
	}
	c = New(os.Getenv("EMIS_USERNAME"), os.Getenv("EMIS_PASSWORD"), os.Getenv("EMIS_HOST"))

	os.Exit(t.Run())
}

func TestEmis_SensorTypes(t *testing.T) {
	_, err := c.SensorTypes()
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmis_Sensors(t *testing.T) {
	_, err := c.Sensors()
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmis_SensorReadings(t *testing.T) {
	res, err := c.SensorReadings(13214, 2022)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res.SensorReadings)
}

func TestEmis_Meters(t *testing.T) {
	res, err := c.Meters()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res.Meters)
}

func TestEmis_MeterReadings(t *testing.T) {
	res, err := c.MeterReadings(1553797269, 2022, 11)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(res.RawBody))
}
