package policy

import (
	"context"
	"testing"
)

func TestEvaluateSLA(t *testing.T) {
	engine := NewOPAEngine()

	// A simple Rego policy for testing
	policyRego := `
package sla

default breach = false

breach {
	input.error_rate > 0.05
}
breach {
	input.p99_latency > 200
}
`

	tests := []struct {
		name       string
		errorRate  float64
		p99Latency float64
		wantBreach bool
	}{
		{
			name:       "Normal Traffic",
			errorRate:  0.01,
			p99Latency: 100.0,
			wantBreach: false,
		},
		{
			name:       "High Error Rate",
			errorRate:  0.06,
			p99Latency: 100.0,
			wantBreach: true,
		},
		{
			name:       "High Latency",
			errorRate:  0.01,
			p99Latency: 250.0,
			wantBreach: true,
		},
		{
			name:       "Both High",
			errorRate:  0.10,
			p99Latency: 300.0,
			wantBreach: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			breach, err := engine.EvaluateSLA(context.Background(), policyRego, tt.errorRate, tt.p99Latency)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if breach != tt.wantBreach {
				t.Errorf("expected breach %v, got %v", tt.wantBreach, breach)
			}
		})
	}
}
