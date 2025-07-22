package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	urlStr := "https://app.courtreserve.com/Online/Reservations/ReadConsolidated/8848"
	desiredCourtIds := map[int]struct{}{28262: {}, 28263: {}, 28430: {}} // Closer courts
	config := GenerateConfig("config.json")
	
	availableSlots := fetchAndFilterBadmintonCourts(urlStr, desiredCourtIds, config)
	units := calculateAvailableSlots(availableSlots, time.Duration(config.ReservationMinInterval)*time.Minute, config.FilterMinDuration)
	for _, unit := range units {
		fmt.Println(unit)
	}
}	

func fetchAndFilterBadmintonCourts(urlStr string, desiredCourtIds map[int]struct{}, config Config) []Slot {
	availableSlots := []Slot{}
	for _, payload := range config.ToPayload() {
		response, err := makeCall(urlStr, payload)
		if err != nil {
			fmt.Printf("Error fetching response: %v\n", err)
			return nil
		}
		for _, slot := range response.Data {
			// find only badminton courts
			if slot.CourtType == "Badminton" && slot.AvailableCourts > 0 && !slot.IsClosed {
				filteredSlot := filterDesiredCourts(slot, desiredCourtIds, config.FilterDesiredCourts)
				if filteredSlot != nil {
					availableSlots = append(availableSlots, *filteredSlot)
				}
			}
		}
	}
	return availableSlots
}

func makeCall(urlStr string, payload Payload) (*Response, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req := genenerateRequest(urlStr, string(jsonPayload))
	response, err := getResponse(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	return response, nil
}

func genenerateRequest(urlStr string, jsonPayload string) *http.Request {
	// Wrap it inside form field `jsonData`
	form := url.Values{}
	form.Set("jsonData", jsonPayload)

	// Prepare the request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://app.courtreserve.com")
	req.Header.Set("Referer", "https://app.courtreserve.com/Online/Reservations/Index/8848")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	return req
}

func getResponse(req *http.Request) (*Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsedResp Response
	err = json.Unmarshal(body, &parsedResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &parsedResp, nil
}