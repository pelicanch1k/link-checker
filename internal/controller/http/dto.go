package http

type CheckLinksRequest struct {
	Links []string `json:"links"`
}

type CheckLinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"` // ID
}

type CheckLinksByIDsRequest struct {
	LinksList []int `json:"links_list"`
}

type CheckLinksByIDsResponse struct {
	Info string `json:"info"`
}