package shared

import (
	"sijil-core/internals/core/identity"
	"sijil-core/internals/core/observability"
	"sijil-core/internals/core/projects"
)

type Handlers struct {
	IdentityRepo  identity.Repository
	Identity      *identity.Handler
	Projects      *projects.Handler
	Observability *observability.Handler
}
