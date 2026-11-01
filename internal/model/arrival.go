package model

type Arrival struct {
	FetchedAt int64 `json:"fetchedAt"` // unix seconds, set by us when the data was fetched
	Data      struct {
		Stops []struct {
			GtfsID   string `json:"gtfsId"`
			Name     string `json:"name"`
			Arrivals []struct {
				ServiceDay      int64  `json:"serviceDay"`
				RealtimeArrival int    `json:"realtimeArrival"`
				Realtime        bool   `json:"realtime"`
				Headsign        string `json:"headsign"`
			} `json:"stoptimesWithoutPatterns"`
		} `json:"stops"`
	} `json:"data"`
}
