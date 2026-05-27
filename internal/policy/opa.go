package policy

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

type OPAEngine struct{}

func NewOPAEngine() *OPAEngine {
	return &OPAEngine{}
}

// EvaluateSLA checks if the current metrics violate the SLA policy
func (e *OPAEngine) EvaluateSLA(ctx context.Context, policyRego string, errorRate, p99Latency float64) (bool, error) {
	query, err := rego.New(
		rego.Query("data.sla.breach"),
		rego.Module("policy.rego", policyRego),
	).PrepareForEval(ctx)
	
	if err != nil {
		return false, fmt.Errorf("failed to prepare rego query: %w", err)
	}

	input := map[string]interface{}{
		"error_rate": errorRate,
		"p99_latency": p99Latency,
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	if len(results) == 0 {
		return false, nil
	}

	breach, ok := results[0].Expressions[0].Value.(bool)
	if !ok {
		return false, fmt.Errorf("policy breach result is not a boolean")
	}

	return breach, nil
}
