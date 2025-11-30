package checker

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pelicanch1k/link-checker/internal/adapter/repository"
	"github.com/pelicanch1k/link-checker/internal/domain"
)

type LinkCheckerUseCase struct {
	httpClient   *http.Client
	workerCount  int
	taskRepo     repository.TaskRepository
	pdfGenerator PDFGenerator
}

type PDFGenerator interface {
	GenerateReport(tasks []*domain.Task) ([]byte, error)
}

func NewLinkCheckerUseCase(timeout time.Duration, workerCount int, repo repository.TaskRepository, pdfGen PDFGenerator) *LinkCheckerUseCase {
	return &LinkCheckerUseCase{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		workerCount:  workerCount,
		taskRepo:     repo,
		pdfGenerator: pdfGen,
	}
}

func (uc *LinkCheckerUseCase) checkURL(url string) domain.LinkStatus {
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

	go func() {
		for i := range links {
			jobs <- i
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		links[result.index].Status = result.status
	}

	// Получаем новый ID и сохраняем задачу
	taskID := uc.taskRepo.GetNextID()
	task := &domain.Task{
		ID:        taskID,
		Links:     links,
		CreatedAt: time.Now(),
	}

	if err := uc.taskRepo.Save(task); err != nil {
		return nil, err
	}

	// Формируем результат
	output := &CheckLinksOutput{
		TaskID: taskID,
		Links:  make(map[string]string),
	}

	for _, link := range links {
		output.Links[link.URL] = string(link.Status)
	}

	return output, nil
}

func (uc *LinkCheckerUseCase) CheckLinksByIDs(input CheckLinksByIDsInput) (*CheckLinksByIDsOutput, error) {
	if len(input.LinksList) == 0 {
		log.Fatal("нет ссылок")
		return nil, domain.ErrEmptyURLs
	}

	tasks, err := uc.taskRepo.FindByIDs(input.LinksList)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if len(tasks) == 0 {
		log.Fatal(err)
		return nil, domain.ErrTaskNotFound
	}

	// Генерируем PDF
	pdfData, err := uc.pdfGenerator.GenerateReport(tasks)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &CheckLinksByIDsOutput{
		PDFData: pdfData,
	}, nil
}
