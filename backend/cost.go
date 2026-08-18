package backend

import "fmt"

// Input is an old-gold exchange at the counter: weight and purity of the piece
// brought in, today's per-gram rate for pure gold, the shop's wastage/making
// deduction, and the price of the new item being bought.
type Input struct {
	OldWeightG       float64 `json:"oldWeightG"`
	PurityPct        float64 `json:"purityPct"`        // e.g. 91.6 for 22k
	RatePerGram      float64 `json:"ratePerGram"`      // ₹ per gram of pure gold
	WastageDeductPct float64 `json:"wastageDeductPct"` // making/wastage deducted on exchange
	NewItemPrice     float64 `json:"newItemPrice"`
}

// Result is the exchange derivation and the balance due.
type Result struct {
	PureContentG  float64 `json:"pureContentG"`
	GrossValue    float64 `json:"grossValue"`
	ExchangeValue float64 `json:"exchangeValue"`
	BalanceDue    float64 `json:"balanceDue"`
}

// Headline is the balance due on the new purchase.
func (r Result) Headline() float64 { return r.BalanceDue }

// Label flags whether the customer still owes money or is due a refund.
func (r Result) Label() string {
	if r.BalanceDue > 0 {
		return "balance-due"
	}
	return "refund-or-even"
}

// Validate reports whether the Input is well formed.
func (in Input) Validate() error {
	if in.OldWeightG < 0 || in.RatePerGram < 0 || in.NewItemPrice < 0 {
		return fmt.Errorf("weight, rate and price cannot be negative")
	}
	if in.PurityPct <= 0 || in.PurityPct > 100 {
		return fmt.Errorf("purity %% must be between 0 and 100")
	}
	if in.WastageDeductPct < 0 || in.WastageDeductPct >= 100 {
		return fmt.Errorf("wastage deduction %% must be between 0 and 100")
	}
	return nil
}

// Evaluate resolves the old-gold value and the balance due on the new item.
func Evaluate(in Input) Result {
	content := in.OldWeightG * in.PurityPct / 100
	gross := content * in.RatePerGram
	exchange := gross * (1 - in.WastageDeductPct/100)
	return Result{
		PureContentG:  content,
		GrossValue:    gross,
		ExchangeValue: exchange,
		BalanceDue:    in.NewItemPrice - exchange,
	}
}
