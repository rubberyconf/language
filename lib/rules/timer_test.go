package rules

import (
	"fmt"
	"testing"
)

func TestTimer(t *testing.T) {
	var tests = []struct {
		sent     FeatureTimer
		expected bool
	}{
		{FeatureTimer{"xxxx Jan 123"}, false},
		{FeatureTimer{"Aug 11, 2021 at 10:00pm (CEST)"}, true},
		{FeatureTimer{"Sep 12, 2030 at 11:00pm (CEST)"}, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("version: %s", tt.sent)
		t.Run(testname, func(t *testing.T) {
			var rt RuleTimer
			ok := rt.evaluate(tt.sent)
			if ok != tt.expected {
				t.Errorf("got %t, want %t", ok, tt.expected)
			}
		})
	}

}
