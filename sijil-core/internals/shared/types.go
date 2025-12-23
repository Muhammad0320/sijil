package shared

import (
	"sijil-core/internals/identity"
	"sijil-core/internals/projects"
)

type Handlers struct {
	Identity *identity.Handler
	Projects *projects.Handler
}
