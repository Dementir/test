package logparser

import (
	"bufio"
	"github.com/Dementir/test/internal/store"
	"log"
	"os"
	"regexp"
	"time"
)

func LogParse(path string) ([]store.Statistic, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	reg := regexp.MustCompile(`[0-9]{4}\/[0-9]{2}\/[0-9]{2}\W[0-9]{2}:[0-9]{2}:[0-9]{2}\WcbHandler\WServeHTTP:\Wargs\Wmap\[channel:\[([^"]*)\]\Wuser:\[([^"]*)\]\Wregion:\[([^"]*)\]\Wserviceid:\[([^"]*)\]\Wstatus:\[([^"]*)\]\Wtime:\[([0-9-,:]*)\]`)
	scanner := bufio.NewScanner(f)

	stats := make([]store.Statistic, 0, 1000000)
	for scanner.Scan() {
		logStr := scanner.Text()

		matches := reg.FindStringSubmatch(logStr)
		if len(matches) < 7 {
			log.Printf("parse log error")
		}

		logTime, err := time.Parse("2006-01-02,15:04:05", matches[6])
		if err != nil {
			return nil, err
		}

		stat := store.Statistic{
			Channel:   matches[1],
			User:      matches[2],
			Region:    matches[3],
			ServiceID: matches[4],
			Status:    matches[5],
			Time:      logTime,
		}

		stats = append(stats, stat)
	}

	return stats, nil
}
