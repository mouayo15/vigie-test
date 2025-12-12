package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

type Order struct {
	ID          string `json:"id"`
	Marketplace string `json:"marketplace"`
	Country     string `json:"country"`
	AmountCents int64  `json:"amount_cents"`
	CreatedAt   string `json:"created_at"`
}

func main() {

	from := flag.String("from", "", "filter from date (YYYY-MM-DD)")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: go run main.go [-from=YYYY-MM-DD] orders.json")
		return
	}
	file := flag.Args()[0]

	var fromDate *time.Time
	if *from != "" {
		t, err := time.Parse("2006-01-02", *from)
		if err != nil {
			fmt.Println("Invalid date:", err)
			return
		}
		fromDate = &t
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	total := int64(0)
	revenueByMP := map[string]int64{}
	suspicious := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		var o Order
		if err := json.Unmarshal([]byte(line), &o); err != nil {
			continue
		}

		// Parse date
		t, err := time.Parse(time.RFC3339, o.CreatedAt)
		if err != nil {
			continue
		}

		// Filter by date
		if fromDate != nil && t.Before(*fromDate) {
			continue
		}

		// Suspicious?
		if o.AmountCents < 0 {
			suspicious = append(suspicious, fmt.Sprintf("- %s: negative amount (%d)", o.ID, o.AmountCents))
			continue
		}
		if o.Marketplace == "" {
			suspicious = append(suspicious, fmt.Sprintf("- %s: empty marketplace", o.ID))
			continue
		}

		total += o.AmountCents
		revenueByMP[o.Marketplace] += o.AmountCents
	}

	// Sort marketplaces by revenue
	type kv struct {
		k string
		v int64
	}
	var arr []kv
	for k, v := range revenueByMP {
		arr = append(arr, kv{k, v})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].v > arr[j].v })

	// Output
	fmt.Printf("Total revenue: %.2f EUR\n\n", float64(total)/100)

	fmt.Println("Revenue by marketplace:")
	for _, e := range arr {
		fmt.Printf("- %s: %.2f EUR\n", e.k, float64(e.v)/100)
	}

	fmt.Println("\nSuspicious orders:")
	for _, s := range suspicious {
		fmt.Println(s)
	}
}
