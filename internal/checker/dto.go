package checker

import "net/http"

// DTO для UseCase
type CheckLinksInput struct {
	URLs []string
}

type CheckLinksOutput struct {
	Links    map[string]string
	LinksNum int
}

type LinkCheckerUseCase struct {
	httpClient *http.Client
	workerCount int
}