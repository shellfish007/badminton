package main

import (
	"fmt"
	"time"
	"sort"
)

type unit struct {
	StartTime time.Time
	EndTime   time.Time
	CourtId   int
}

func (u unit) String() string {
    loc, err := time.LoadLocation("America/New_York")
    if err != nil {
        panic(err)
    }
    newYorkTime := u.StartTime.In(loc)
    // Format without year but with day in a week
    return fmt.Sprintf("%s %s-%s, CourtId: %d", newYorkTime.Format("Mon 02"), u.StartTime.In(loc).Format("15:04"), u.EndTime.In(loc).Format("15:04"), u.CourtId)
}

func calculateAvailableSlots(availableSlots []Slot, minDuration time.Duration, filterMinDuration bool) []unit {
	var units []unit
	previousUnit := map[int]unit{} // Track the end time of the last slot for each court
	for _, slot := range availableSlots {
		if len(slot.AvailableCourtIds) == 0 {
			continue // Skip slots with no available courts
		}
		startTime := slot.Start.Time()
		endTime := slot.End.Time() 
		// endTime - startTime = 30min
		for _, courtId := range slot.AvailableCourtIds {
			if previous, exists := previousUnit[courtId]; exists {
				if previous.EndTime.Equal(startTime) {
					// If the previous slot ends exactly when this one starts, merge them
					previous.EndTime = endTime
					previousUnit[courtId] = previous // Update the previous unit
					continue
				} else {
					units = append(units, previous) // Append the previous unit if it exists
				}
			}
			// If not merging, create a new unit
			previousUnit[courtId] = unit{
				StartTime: startTime,
				EndTime:   endTime,
				CourtId:   courtId,
			}
		}
	}
	// Append the last units for each court
	for _, previous := range previousUnit {
		units = append(units, previous)
	}
	// Filter units based on minimum duration if required
	if filterMinDuration {
		var filteredUnits []unit
		for _, u := range units {
			if u.EndTime.Sub(u.StartTime) >= minDuration {
				filteredUnits = append(filteredUnits, u)
			}
		}
		units = filteredUnits
	}
	// Sort units by start time
	sort.Slice(units, func(i, j int) bool {
		return units[i].StartTime.Before(units[j].StartTime)
	})
	return units
}

func filterDesiredCourts(slot Slot, desiredCourtIds map[int]struct{}, filter bool) *Slot {
	if filter {
		desiredCourts := []int{}
		for _, courtId := range slot.AvailableCourtIds {
			if _, ok := desiredCourtIds[courtId]; ok {
				desiredCourts = append(desiredCourts, courtId)
			}
		}
		if len(desiredCourts) > 0 {
			slot.AvailableCourtIds = desiredCourts
			return &slot // Return the slot with filtered courts
		}
	} else {
		return &slot // Return the slot as is if filtering is not required
	}
	return nil
}