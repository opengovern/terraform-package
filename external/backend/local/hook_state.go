// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package local

import (
	"log"
	"sync"
	"time"

	"github.com/opengovern/terraform-package/external/states"
	"github.com/opengovern/terraform-package/external/states/statemgr"
	"github.com/opengovern/terraform-package/external/terraform"
)

// StateHook is a hook that continuously updates the state by calling
// WriteState on a statemgr.Full.
type StateHook struct {
	terraform.NilHook
	sync.Mutex

	StateMgr statemgr.Writer

	// If PersistInterval is nonzero then for any new state update after
	// the duration has elapsed we'll try to persist a state snapshot
	// to the persistent backend too.
	// That's only possible if field Schemas is valid, because the
	// StateMgr.PersistState function for some backends needs schemas.
	PersistInterval time.Duration

	// Schemas are the schemas to use when persisting state due to
	// PersistInterval. This is ignored if PersistInterval is zero,
	// and PersistInterval is ignored if this is nil.
	Schemas *terraform.Schemas

	lastPersist  time.Time
	forcePersist bool
}

var _ terraform.Hook = (*StateHook)(nil)

func (h *StateHook) PostStateUpdate(new *states.State) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if h.lastPersist.IsZero() {
		// The first PostStateUpdate starts the clock for intermediate
		// calls to PersistState.
		h.lastPersist = time.Now()
	}

	if h.StateMgr != nil {
		if err := h.StateMgr.WriteState(new); err != nil {
			return terraform.HookActionHalt, err
		}
		if mgrPersist, ok := h.StateMgr.(statemgr.Persister); ok && h.PersistInterval != 0 && h.Schemas != nil {
			if h.forcePersist || time.Since(h.lastPersist) >= h.PersistInterval {
				err := mgrPersist.PersistState(h.Schemas)
				if err != nil {
					return terraform.HookActionHalt, err
				}
				h.lastPersist = time.Now()
			}
		}
	}

	return terraform.HookActionContinue, nil
}

func (h *StateHook) Stopping() {
	h.Lock()
	defer h.Unlock()

	// If Terraform has been asked to stop then that might mean that a hard
	// kill signal will follow shortly in case Terraform doesn't stop
	// quickly enough, and so we'll try to persist the latest state
	// snapshot in the hope that it'll give the user less recovery work to
	// do if they _do_ subsequently hard-kill Terraform during an apply.

	if mgrPersist, ok := h.StateMgr.(statemgr.Persister); ok && h.Schemas != nil {
		err := mgrPersist.PersistState(h.Schemas)
		if err != nil {
			// This hook can't affect Terraform Core's ongoing behavior,
			// but it's a best effort thing anyway so we'll just emit a
			// log to aid with debugging.
			log.Printf("[ERROR] Failed to persist state after interruption: %s", err)
		}

		// While we're in the stopping phase we'll try to persist every
		// new state update to maximize every opportunity we get to avoid
		// losing track of objects that have been created or updated.
		// Terraform Core won't start any new operations after it's been
		// stopped, so at most we should see one more PostStateUpdate
		// call per already-active request.
		h.forcePersist = true
	}

}
