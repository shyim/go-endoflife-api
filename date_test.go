package endoflife

import (
	"encoding/json"
	"testing"
)

func TestDateRoundTrip(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(`"2022-04-21"`), &d); err != nil {
		t.Fatal(err)
	}
	if d.String() != "2022-04-21" {
		t.Errorf("String() = %q", d.String())
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `"2022-04-21"` {
		t.Errorf("Marshal = %s", out)
	}
}

func TestDateNull(t *testing.T) {
	var p struct {
		D *Date `json:"d"`
	}
	if err := json.Unmarshal([]byte(`{"d": null}`), &p); err != nil {
		t.Fatal(err)
	}
	if p.D != nil {
		t.Errorf("expected nil *Date, got %v", p.D)
	}

	// Zero Date marshals to null.
	out, err := json.Marshal(Date{})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "null" {
		t.Errorf("zero Date Marshal = %s", out)
	}
}

func TestDateInvalid(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(`"not-a-date"`), &d); err == nil {
		t.Error("expected error for invalid date")
	}
}
