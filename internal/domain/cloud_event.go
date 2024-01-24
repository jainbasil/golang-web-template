package domain

type CloudEvent struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Data   any    `json:"data"`
}
