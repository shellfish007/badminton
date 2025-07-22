package main

import (
    "fmt"
    "strings"
    "strconv"
	"regexp"
    "time"
)

// Slot represents a single court slot entry from the response
type Slot struct {
    Id                        Id          `json:"Id"` //e.g. "Badminton07/23/2025 01:30:00"
    Title                     *string     `json:"Title"`
    Description               *string     `json:"Description"`
    IsAllDay                  bool        `json:"IsAllDay"`
    Start                     Point       `json:"Start"` // e.g. "\/Date(1753232400000)\/"
    End                       Point       `json:"End"` // e.g. "\/Date(1753234200000)\/"
    StartTimezone             *string     `json:"StartTimezone"`
    EndTimezone               *string     `json:"EndTimezone"`
    RecurrenceRule            *string     `json:"RecurrenceRule"`
    RecurrenceException       *string     `json:"RecurrenceException"`
    CourtType                 string      `json:"CourtType"`
    AvailableCourts           int         `json:"AvailableCourts"`
    IsClosed                  bool        `json:"IsClosed"`
    IsInPast                  bool        `json:"IsInPast"`
    IsWaitListSlot            bool        `json:"IsWaitListSlot"`
    ShowWaitList              bool        `json:"ShowWaitList"`
    WaitListCount             *int        `json:"WaitListCount"`
    IsAvailableTemplate       bool        `json:"IsAvailableTemplate"`
    IsCourtAssignmentHiddenOnPortal bool  `json:"IsCourtAssignmentHiddenOnPortal"`
    MemberIds                 []int       `json:"MemberIds"`
    AvailableCourtIds         []int       `json:"AvailableCourtIds"`
    QueuedMembers             []int       `json:"QueuedMembers"`
    ShowCourtWaitlistOrderNumber bool     `json:"ShowCourtWaitlistOrderNumber"`
}

type Id string

func (id Id) Time() time.Time {
    // id is in the format "Badminton07/24/2025 03:30:00"
    // Parse the string to extract the time
    // Remove the "Badminton" prefix
    t, err := time.Parse("01/02/2006 15:04:05", strings.TrimPrefix(string(id), "Badminton"))
    if err != nil {
        panic("Error parsing time: " + err.Error())
        return time.Time{} // Return zero time on error
    }
    // Return the parsed time
    return t
}

type Point string

func (p Point) Time() time.Time {
    // Assuming Point is a string representation of time in the format "\/Date(1753232400000)\/"
    re := regexp.MustCompile(`\d+`)
	matches := re.FindStringSubmatch(string(p))

	if len(matches) == 0 {
		panic("No matches found in Point string")
	}
	millisStr := matches[0]
	millis, err := strconv.ParseInt(millisStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(millis/1000, (millis%1000)*1000000)
}

func (s Slot) String() string {
    loc, err := time.LoadLocation("America/New_York")
    if err != nil {
        panic(err)
    }
    newYorkTime := s.Id.Time().In(loc)
    // Format without year but with day in a week
    return fmt.Sprintf("%s %s-%s, AvailableCourtIds: %v", newYorkTime.Format("Mon 02"), s.Start.Time().In(loc).Format("15:04"), s.End.Time().In(loc).Format("15:04"), s.AvailableCourtIds)
}

func (s Slot) Print() { 
    fmt.Println(s.String())
}

// Response represents the top-level response structure
type Response struct {
    Data             []Slot      `json:"Data"`
    Total            int         `json:"Total"`
    AggregateResults interface{} `json:"AggregateResults"`
    Errors           interface{} `json:"Errors"`
}
