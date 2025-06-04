package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	url        = "https://idlc.com/aml/nav.php"
	outputFile = "updated-nav.txt"
)

func main() {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	table := doc.Find("table.table").First()
	if table.Length() == 0 {
		log.Fatal("Could not find NAV table")
	}

	var sb strings.Builder

	// Extract headers
	table.Find("thead tr th").Each(func(i int, s *goquery.Selection) {
		sb.WriteString(strings.TrimSpace(s.Text()))
		sb.WriteString("\t")
	})
	sb.WriteString("\n")

	// Extract rows
	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			sb.WriteString(strings.TrimSpace(cell.Text()))
			sb.WriteString("\t")
		})
		sb.WriteString("\n")
	})

	if err := os.WriteFile(outputFile, []byte(sb.String()), 0644); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("✅ NAV data updated in", outputFile)
}
