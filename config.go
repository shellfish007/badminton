package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	StartDate              string     `json:"startDate"`
	Days                   int        `json:"days"`
	OrgId                  string     `json:"orgId"`
	TimeZone               string     `json:"TimeZone"`
	UiCulture              string     `json:"UiCulture"`
	CostTypeId             string     `json:"CostTypeId"`
	CustomSchedulerId      string     `json:"CustomSchedulerId"`
	ReservationMinInterval int        `json:"ReservationMinInterval"`
	DesiredCourtIds	       []int      `json:"DesiredCourtIds"`
	FilterDesiredCourts    bool       `json:"filterDesiredCourts"`
	FilterMinDuration      bool       `json:"filterMinDuration"`
}

// Payload represents the request body structure for the reservation API
type Payload struct {
	StartDate              string     `json:"startDate"`
	OrgId                  string     `json:"orgId"`
	TimeZone               string     `json:"TimeZone"`
	Date                   string     `json:"Date"`
	KendoDate              KendoDate  `json:"KendoDate"`
	UiCulture              string     `json:"UiCulture"`
	CostTypeId             string     `json:"CostTypeId"`
	CustomSchedulerId      string     `json:"CustomSchedulerId"`
	ReservationMinInterval string     `json:"ReservationMinInterval"`
}

type KendoDate struct {
	Year  int `json:"Year"`
	Month int `json:"Month"`
	Day   int `json:"Day"`
}

func GenerateConfig(file string) Config {
	// Read the JSON file and unmarshal it into Config struct
	jsonBytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		panic(err)
	}

	defaultConfig := DefaultConfig()
	if config.StartDate == "" {
		config.StartDate = defaultConfig.StartDate
	} else {
		startDate, err := time.Parse("2006-01-02T15:04:05-04:00", config.StartDate)
		if err != nil {
			panic(err)
		}
		config.StartDate = startDate.Format("2006-01-02T15:04:05.000Z")
	}
	if config.Days <= 0 {
		config.Days = defaultConfig.Days
	}
	if config.OrgId == "" {
		config.OrgId = defaultConfig.OrgId
	}
	if config.TimeZone == "" {
		config.TimeZone = defaultConfig.TimeZone
	}
	if config.UiCulture == "" {
		config.UiCulture = defaultConfig.UiCulture
	}
	if config.CostTypeId == "" {
		config.CostTypeId = defaultConfig.CostTypeId
	}
	if config.CustomSchedulerId == "" {
		config.CustomSchedulerId = defaultConfig.CustomSchedulerId
	}
	if config.ReservationMinInterval == 0 {
		config.ReservationMinInterval = defaultConfig.ReservationMinInterval
	}
	if len(config.DesiredCourtIds) == 0 {
		config.DesiredCourtIds = defaultConfig.DesiredCourtIds
	}
	return config
}

func DefaultConfig() Config {
	return Config{
		StartDate:              time.Now().Format("2006-01-02T15:04:05.000Z"),
		Days:                   7,
		OrgId:                  "8848",
		TimeZone:               "America/New_York",
		UiCulture:              "en-US",
		CostTypeId:             "98331",
		CustomSchedulerId:      "",
		ReservationMinInterval: 60,
		DesiredCourtIds:        []int{28262, 28263, 28430}, // Closer courts
	}
}

func (c Config) ToPayload() []Payload {
	payloads := make([]Payload, c.Days)
	payload := Payload{
		OrgId:                  c.OrgId,
		TimeZone:               c.TimeZone,
		UiCulture:              c.UiCulture,
		CostTypeId:             c.CostTypeId,
		CustomSchedulerId:      c.CustomSchedulerId,
		ReservationMinInterval: fmt.Sprintf("%d", c.ReservationMinInterval),
	}
	startDate, err := time.Parse("2006-01-02T15:04:05.000Z", c.StartDate)
	if err != nil {
		panic(fmt.Errorf("error parsing start date: %v", err))
	}
	for i := 0; i < c.Days; i++ {
		d := startDate.AddDate(0, 0, i)
		payload.StartDate = d.Format("2006-01-02T15:04:05.000Z")
		payload.Date = d.Format("Mon, 02 Jan 2006 15:04:05 GMT")
		payload.KendoDate = KendoDate{Year: d.Year(), Month: int(d.Month()), Day: d.Day()}
		payloads[i] = payload
	}
	return payloads
}