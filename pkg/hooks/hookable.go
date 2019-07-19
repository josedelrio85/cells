package leads

import (
	model "github.com/bysidecar/leads/pkg/model"
)

// HookResponse defines the information available for a Hook to return data to
// the main workflow, its main responsibility is to alter the default workflow
// for the leads (returning error?) when they don't meet the Hook criteria.
type HookResponse struct {
	Err        error
	StatusCode int
	Result     bool
}

// Hookable defines the interface to perform Hook actions on leads. A Hook is a
// custom code that will be executed on a precise and defined moment on Lead's
// lifecycle that will alter the default workflow for that Lead.
//
// A Hookable also has an activation function that will inform if the Hookable
// should be triggered for the given Lead.
type Hookable interface {
	Active(model.Lead) bool
	Perform(*model.Lead) HookResponse
	Test(interface{})
}
