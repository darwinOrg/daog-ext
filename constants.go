package daogext

import "github.com/rolandhe/daog"

var (
	IdOrderAsc  = daog.NewOrder("id")
	IdOrderDesc = daog.NewDescOrder("id")

	CreatedAtOrderAsc  = daog.NewOrder("created_at")
	CreatedAtOrderDesc = daog.NewDescOrder("created_at")

	ModifiedAtOrderAsc  = daog.NewOrder("modified_at")
	ModifiedAtOrderDesc = daog.NewDescOrder("modified_at")
)
