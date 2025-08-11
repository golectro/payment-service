package model

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Message map[string]string

type WebResponse[T any] struct {
	Status           Status        `json:"status"`
	StatusCode       int           `json:"statusCode"`
	Code             string        `json:"code,omitempty"`
	Message          Message       `json:"message"`
	Data             T             `json:"data,omitempty"`
	Paging           *PageMetadata `json:"paging,omitempty"`
	Errors           string        `json:"errors,omitempty"`
	RequestID        string        `json:"requestId"`
	Timestamp        string        `json:"timestamp"`
	Path             string        `json:"path"`
	DocumentationURL string        `json:"documentationUrl"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging"`
}

type PageMetadata struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItem   int64 `json:"total_item"`
	TotalPage   int64 `json:"total_page"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}
