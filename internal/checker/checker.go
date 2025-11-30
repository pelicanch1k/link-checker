package checker

import (
	"net/http"
	"sync"
	"time"

	"github.com/pelicanch1k/link-checker/internal/domain"
)

func NewLinkCheckerUseCase(timeout time.Duration, workerCount int) *LinkCheckerUseCase {
	return &LinkCheckerUseCase{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		workerCount: workerCount,
	}
}

func (uc *LinkCheckerUseCase) checkURL(url string) domain.LinkStatus {
	// Добавляем схему, если её нет
	if len(url) > 0 && url[0] != 'h' {
		url = "http://" + url
	}

	resp, err := uc.httpClient.Get(url)
	if err != nil {
		return domain.StatusNotAvailable
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return domain.StatusAvailable
	}
	return domain.StatusNotAvailable
}

func (uc *LinkCheckerUseCase) CheckLinks(input CheckLinksInput) (*CheckLinksOutput, error) {
	if len(input.URLs) == 0 {
		return nil, domain.ErrEmptyURLs
	}

	links := make([]domain.Link, len(input.URLs))
	for i, url := range input.URLs {
		links[i] = domain.Link{
			URL:    url,
			Status: "",
		}
	}

	type response struct {
		index  int
		status domain.LinkStatus
	}

	jobs := make(chan int, len(links))
	results := make(chan response, len(links))

	var wg sync.WaitGroup

	for w := 0; w < uc.workerCount; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for index := range jobs {
				status := uc.checkURL(links[index].URL)
				results <- response{index, status}
			}
		}()
	}

	// Отправляем задачи
	go func() {
		for i := range links {
			jobs <- i
		}
		close(jobs)
	}()

	// Ждем результаты
	go func() {
		wg.Wait()
		close(results)
	}()

	// Обновляем статсы
	for result := range results {
		links[result.index].Status = result.status
	}

	// Формируем результат
	output := &CheckLinksOutput{
		Links:    make(map[string]string),
		LinksNum: len(links),
	}

	for _, link := range links {
		output.Links[link.URL] = string(link.Status)
	}

	return output, nil
}