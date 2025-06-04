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
	targetFund = "IDLC Asset Management Shariah Fund"
)

func main() {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("❌ Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("❌ Unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("❌ Failed to parse HTML: %v", err)
	}

	var foundLine string
	doc.Find("table.table tbody tr").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		cells := row.Find("td")
		if cells.Length() == 0 {
			return true // skip empty rows
		}

		fundName := strings.TrimSpace(cells.Eq(0).Text())
		if fundName == targetFund {
			var sb strings.Builder
			cells.Each(func(i int, cell *goquery.Selection) {
				sb.WriteString(strings.TrimSpace(cell.Text()))
				sb.WriteString("\t")
			})
			foundLine = strings.TrimRight(sb.String(), "\t") + "\n"
			return false // stop after finding the match
		}
		return true
	})

	if foundLine == "" {
		log.Fatalf("❌ Target fund not found: %s", targetFund)
	}

	// Write the exact line to updated-nav.txt
	err = os.WriteFile(outputFile, []byte(foundLine), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write to file: %v", err)
	}

	fmt.Println("✅ NAV data written to", outputFile)
}
