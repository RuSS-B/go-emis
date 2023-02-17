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
}

func TestEmis_GetSensorTypes(t *testing.T) {
	res, err := c.GetSensorTypes()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(res))
}

func TestEmis_Sensors(t *testing.T) {
	res, err := c.Sensors()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestEmis_SensorReadings(t *testing.T) {
	res, err := c.SensorReadings(13214, 2023)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
