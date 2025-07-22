# badminton

Filter courts and minimum duration for a time range for kotofit badminton reservation

## Sample input

{
    "startDate":"",
    "days": 14,
    "filterDesiredCourts": true,
    "filterMinDuration": true,
    "orgId":"8848",
    "TimeZone":"America/New_York",
    "UiCulture":"en-US",
    "CostTypeId":"98331",
    "CustomSchedulerId":"",
    "ReservationMinInterval":60
}

## Sample output

Wed 23 22:00-23:30, CourtId: 28263
Wed 23 22:00-23:30, CourtId: 28430
Wed 23 22:00-23:30, CourtId: 28262
Fri 25 21:30-23:30, CourtId: 28263
Fri 25 22:00-23:30, CourtId: 28430
Mon 28 21:00-22:00, CourtId: 28430
Wed 30 22:00-23:30, CourtId: 28430
Wed 30 22:00-23:30, CourtId: 28262
Wed 30 22:00-23:30, CourtId: 28263
Fri 01 21:00-23:30, CourtId: 28262
Fri 01 22:00-23:30, CourtId: 28430
Fri 01 22:00-23:30, CourtId: 28263
Mon 04 19:00-20:00, CourtId: 28430
Mon 04 20:30-22:00, CourtId: 28263
