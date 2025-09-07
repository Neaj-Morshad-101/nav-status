package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	url         = "https://idlc.com/aml/nav.php"
	outputFile  = "updated-nav.txt"
	targetFund  = "IDLC Asset Management Shariah Fund"
	avgBuyPrice = 10.2633
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
			return true // skip if no <td>
		}

		fundName := strings.TrimSpace(cells.Eq(0).Text())
		if fundName == targetFund {
			var sb strings.Builder
			cells.Each(func(i int, cell *goquery.Selection) {
				sb.WriteString(strings.TrimSpace(cell.Text()))
				sb.WriteString("\t")
			})
			foundLine = strings.TrimRight(sb.String(), "\t") + "\n"
			return false // stop iteration once found
		}
		return true
	})

	if foundLine == "" {
		log.Fatalf("❌ Target fund not found: %s", targetFund)
	}

	// Extract the NAV per unit from the found line.
	// The line has columns: FundName, NAVasOn, NAVPerUnit, InvestorBuyPrice, InvestorSalePrice
	parts := strings.Split(foundLine, "\t")
	if len(parts) < 3 {
		log.Fatalf("❌ Unexpected format in foundLine: %s", foundLine)
	}
	navStr := parts[2]
	navValue, err := strconv.ParseFloat(navStr, 64)
	if err != nil {
		log.Fatalf("❌ Unable to parse NAV value '%s': %v", navStr, err)
	}

	// Calculate profit percentage: ((NAV - avgBuyPrice) / avgBuyPrice) * 100
	diffPct := (navValue - avgBuyPrice) / avgBuyPrice * 100

	// Format profit line, rounding to one decimal place
	profitLine := fmt.Sprintf("Profit since buy @%.2f: %.1f%%\n", avgBuyPrice, diffPct)

	// Write both lines to updated-nav.txt
	output := foundLine + profitLine
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		log.Fatalf("❌ Failed to write to file: %v", err)
	}

	fmt.Println("✅ NAV data and profit line written to", outputFile)
}
