package agent

import "github.com/donCannoli-burns/kingdomsitter/internal/analyze"

// Envelope is provider-neutral. It is meant to be sent to an LLM/agent layer
// without granting that agent direct runtime authority.
type Envelope struct {
	Schema   string         `json:"schema"`
	Task     string         `json:"task"`
	Analysis analyze.Report `json:"analysis"`
	Policy   Policy         `json:"policy"`
}

type Policy struct {
	ExecutionAuthority bool     `json:"execution_authority"`
	AllowedOutputs     []string `json:"allowed_outputs"`
}

func New(task string, report analyze.Report) Envelope {
	return Envelope{
		Schema:   "kingdomsitter.agent-envelope/v0",
		Task:     task,
		Analysis: report,
		Policy:   Policy{ExecutionAuthority: false, AllowedOutputs: []string{"analysis", "patch", "typescript", "explanation"}},
	}
}
