package backend

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvaluate_ExchangeAndBalance(t *testing.T) {
	// 10g @ 91.6% = 9.16g pure ; rate ₹6000/g => gross 54960 ; 8% wastage => 50563.2
	// new item ₹80000 => balance 29436.8
	r := Evaluate(Input{OldWeightG: 10, PurityPct: 91.6, RatePerGram: 6000, WastageDeductPct: 8, NewItemPrice: 80000})
	if math.Abs(r.PureContentG-9.16) > 1e-9 {
		t.Fatalf("content=%v want 9.16", r.PureContentG)
	}
	if math.Abs(r.GrossValue-54960) > 1e-6 {
		t.Fatalf("gross=%v want 54960", r.GrossValue)
	}
	if math.Abs(r.ExchangeValue-50563.2) > 1e-4 {
		t.Fatalf("exchange=%v want 50563.2", r.ExchangeValue)
	}
	if math.Abs(r.BalanceDue-29436.8) > 1e-4 {
		t.Fatalf("balance=%v want 29436.8", r.BalanceDue)
	}
}

func TestValidate(t *testing.T) {
	if err := (Input{OldWeightG: 10, PurityPct: 91.6, RatePerGram: 6000}).Validate(); err != nil {
		t.Fatalf("valid rejected: %v", err)
	}
	for i, bad := range []Input{{PurityPct: 0}, {PurityPct: 91.6, WastageDeductPct: 100}, {PurityPct: 91.6, OldWeightG: -1}} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("bad %d accepted", i)
		}
	}
}

func TestEvaluateEndpoint(t *testing.T) {
	srv := NewServer(nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/evaluate",
		strings.NewReader(`{"oldWeightG":10,"purityPct":91.6,"ratePerGram":6000,"wastageDeductPct":8,"newItemPrice":80000}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var r Result
	json.Unmarshal(rec.Body.Bytes(), &r)
	if math.Abs(r.BalanceDue-29436.8) > 1e-3 {
		t.Fatalf("balance=%v", r.BalanceDue)
	}
}
