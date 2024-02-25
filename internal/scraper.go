package internal

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/alitto/pond"
	"github.com/schollz/progressbar/v3"
	"html-text-extractor/package/scrape"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func ScrapeDomainsFromFile(fname string, requestTimeout time.Duration, maxConcurrentRequests int, ignoredTags map[string]struct{}, outputFile string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	outFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	logFile, err := os.OpenFile("log.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	log.SetOutput(logFile)

	defer logFile.Close()
	defer file.Close()
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	writerLock := sync.Mutex{}

	scanner := bufio.NewScanner(file)
	pgbar := progressbar.Default(-1, fmt.Sprintf("Crawling (max concurrent requests: %d)", maxConcurrentRequests))
	pool := pond.New(maxConcurrentRequests, 0, pond.MinWorkers(maxConcurrentRequests))
	client := http.Client{Timeout: requestTimeout}

	extractor := scrape.Extractor{IgnoredTags: ignoredTags}

	for scanner.Scan() {
		domain := scanner.Text()
		pool.Submit(func() {
			req, err := http.NewRequest("GET", fmt.Sprintf("http://%s", domain), nil)
			req.Header.Set("Accept", "text/html")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
			if err != nil {
				log.Printf("Error occurred during scraping of %s : %s\n", domain, err.Error())
				return
			}
			res, err := client.Do(req)
			if err != nil {
				log.Printf("Error occurred during scraping of %s : %s\n", domain, err.Error())
				return
			}
			result := extractor.ExtractFromReader(res.Body)
			res.Body.Close()
			writerLock.Lock()
			err = writer.Write([]string{domain, result})
			if err != nil {
				log.Printf("Error occurred during scraping of %s : %s\n", domain, err.Error())
			}
			writerLock.Unlock()
			pgbar.Add(1)
		})
	}
	pool.StopAndWait()
	writer.Flush()
	return nil
}
