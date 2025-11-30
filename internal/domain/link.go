package domain

type LinkStatus string

const (
    StatusAvailable    LinkStatus = "available"
    StatusNotAvailable LinkStatus = "not available"
)

type Link struct {
    URL    string     `json:"url"`
    Status LinkStatus `json:"status"`
}