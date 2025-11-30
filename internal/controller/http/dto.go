package http

type CheckLinksRequest struct {
	Links []string `json:"links"`
}

type CheckLinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}
