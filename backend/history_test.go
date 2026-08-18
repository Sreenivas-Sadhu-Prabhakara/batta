package backend

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

// Integration test for the Postgres store. Runs only when DATABASE_URL is set
// (scripts/test.sh spins up a throwaway Postgres, then tears it down).
func TestPostgresStore_RoundTrip(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping Postgres integration test")
	}
	store, err := NewPostgresStore(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer store.Close()

	in, _ := json.Marshal(Input{OldWeightG: 10, PurityPct: 91.6, RatePerGram: 6000, WastageDeductPct: 8, NewItemPrice: 80000})
	saved, err := store.Save(Record{Input: in, Headline: 29436.8, Label: "balance-due"})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.ID == 0 {
		t.Fatal("expected assigned id")
	}
	items, err := store.List(10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) == 0 || items[0].ID != saved.ID || items[0].Label != "balance-due" {
		t.Fatalf("unexpected list head: %+v", items)
	}
}
