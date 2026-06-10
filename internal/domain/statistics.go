package domain

type DailyCount struct {
	Date  string `json:"date" db:"day"`
	Count int64  `json:"count" db:"count"`
}

type Statistics struct {
	Total     int64        `json:"total"`
	Today     int64        `json:"today"`
	Last7Days []DailyCount `json:"last_7_days"`
}
