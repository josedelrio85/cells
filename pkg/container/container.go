package leads

import model "github.com/bysidecar/leads/pkg/model"

// Container is an auxiliar struct used to pass params to Perform action
type Container struct {
	Storer model.Storer
	Lead   model.Lead
	Redis  model.Redis
}
