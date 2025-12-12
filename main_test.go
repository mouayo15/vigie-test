package main

import (
	"encoding/json"
	"testing"
	"time"
)

// --- Test du parsing d'une ligne JSON ---
func TestParseOrder(t *testing.T) {
	line := `{"id":"o1","marketplace":"amazon","country":"FR","amount_cents":1500,"created_at":"2024-11-01T10:00:00Z"}`

	var o Order
	if err := json.Unmarshal([]byte(line), &o); err != nil {
		t.Fatalf("Erreur de parsing JSON: %v", err)
	}

	if o.ID != "o1" {
		t.Errorf("ID incorrect, attendu o1, obtenu %s", o.ID)
	}
	if o.AmountCents != 1500 {
		t.Errorf("Amount incorrect, attendu 1500, obtenu %d", o.AmountCents)
	}
}

// --- Test du filtre date ---
func TestDateFilter(t *testing.T) {
	orderDate := "2024-11-01T10:00:00Z"
	tm, _ := time.Parse(time.RFC3339, orderDate)

	from, _ := time.Parse("2006-01-02", "2024-11-02")

	if tm.Before(from) != true {
		t.Errorf("Le filtre date ne fonctionne pas: la commande devrait être filtrée")
	}
}

// --- Test simple de totalisation ---
func TestRevenueSum(t *testing.T) {
	orders := []Order{
		{ID: "1", Marketplace: "amazon", AmountCents: 1000},
		{ID: "2", Marketplace: "amazon", AmountCents: 1500},
	}

	total := int64(0)
	for _, o := range orders {
		total += o.AmountCents
	}

	if total != 2500 {
		t.Errorf("Total incorrect : attendu 2500, obtenu %d", total)
	}
}
