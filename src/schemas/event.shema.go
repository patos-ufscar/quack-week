package schemas

type CreateEvent struct {
	EventName   string  `json:"eventName"`
	CoverBase64 *string `json:"coverBase64"`
}
