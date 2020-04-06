package scenario

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"testing"
)

func TestTargetModifier(t *testing.T) {
	f := func(tgt vegeta.Target) vegeta.Target {
		tgt.Method = "POST"
		return tgt

	}
	targeter := StaticInterceptedTargeter("", f, vegeta.Target{
		Method: "GET",
		URL:    "http://example.com",
	})

	tgt := &vegeta.Target{}
	err := targeter.Decode(tgt)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	want := &vegeta.Target{
		Method: "POST",
		URL:    "http://example.com",
	}
	if !tgt.Equal(want) {
		t.Errorf("Not equal! got: %v; want: %v", *tgt, *want)
	}
}

func TestTargetModifierWithOverride(t *testing.T) {
	f := func(tgt vegeta.Target) vegeta.Target {
		tgt.Method = "POST"
		return tgt

	}
	targeter := StaticInterceptedTargeter("http://override.com", f, vegeta.Target{
		Method: "GET",
		URL:    "http://example.com",
	})

	tgt := &vegeta.Target{}
	err := targeter.Decode(tgt)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	want := &vegeta.Target{
		Method: "POST",
		URL:    "http://override.com",
	}
	if !tgt.Equal(want) {
		t.Errorf("Not equal! got: %v; want: %v", *tgt, *want)
	}
}
