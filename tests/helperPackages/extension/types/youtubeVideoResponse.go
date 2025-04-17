package types

type YouTubeVideoResponse struct {
	Status                 int     `json:"status"`
	Message                string  `json:"message"`
	StartTime              float64 `json:"startTime"`
	EndTime                float64 `json:"endTime"`
	ContainSponserSubtitle bool    `json:"containSponserSubtitle"`
}
