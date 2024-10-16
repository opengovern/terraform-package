// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloud

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/opengovern/terraform-package/external/states/statefile"
	"github.com/opengovern/terraform-package/external/states/statemgr"
)

func TestState_impl(t *testing.T) {
	var _ statemgr.Reader = new(State)
	var _ statemgr.Writer = new(State)
	var _ statemgr.Persister = new(State)
	var _ statemgr.Refresher = new(State)
	var _ statemgr.OutputReader = new(State)
	var _ statemgr.Locker = new(State)
}

type ExpectedOutput struct {
	Name      string
	Sensitive bool
	IsNull    bool
}

func TestState_GetRootOutputValues(t *testing.T) {
	b, bCleanup := testBackendWithOutputs(t)
	defer bCleanup()

	state := &State{tfeClient: b.client, organization: b.organization, workspace: &tfe.Workspace{
		ID: "ws-abcd",
	}}
	outputs, err := state.GetRootOutputValues()

	if err != nil {
		t.Fatalf("error returned from GetRootOutputValues: %s", err)
	}

	cases := []ExpectedOutput{
		{
			Name:      "sensitive_output",
			Sensitive: true,
			IsNull:    false,
		},
		{
			Name:      "nonsensitive_output",
			Sensitive: false,
			IsNull:    false,
		},
		{
			Name:      "object_output",
			Sensitive: false,
			IsNull:    false,
		},
		{
			Name:      "list_output",
			Sensitive: false,
			IsNull:    false,
		},
	}

	if len(outputs) != len(cases) {
		t.Errorf("Expected %d item but %d were returned", len(cases), len(outputs))
	}

	for _, testCase := range cases {
		so, ok := outputs[testCase.Name]
		if !ok {
			t.Fatalf("Expected key %s but it was not found", testCase.Name)
		}
		if so.Value.IsNull() != testCase.IsNull {
			t.Errorf("Key %s does not match null expectation %v", testCase.Name, testCase.IsNull)
		}
		if so.Sensitive != testCase.Sensitive {
			t.Errorf("Key %s does not match sensitive expectation %v", testCase.Name, testCase.Sensitive)
		}
	}
}

func TestState(t *testing.T) {
	var buf bytes.Buffer
	s := statemgr.TestFullInitialState()
	sf := statefile.New(s, "stub-lineage", 2)
	err := statefile.Write(sf, &buf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	data := buf.Bytes()

	state := testCloudState(t)

	jsonState, err := ioutil.ReadFile("../command/testdata/show-json-state/sensitive-variables/output.json")
	if err != nil {
		t.Fatal(err)
	}

	jsonStateOutputs := []byte(`
{
	"outputs": {
			"foo": {
					"type": "string",
					"value": "bar"
			}
	}
}`)

	if err := state.uploadState(state.lineage, state.serial, state.forcePush, data, jsonState, jsonStateOutputs); err != nil {
		t.Fatalf("put: %s", err)
	}

	payload, err := state.getStatePayload()
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	if !bytes.Equal(payload.Data, data) {
		t.Fatalf("expected full state %q\n\ngot: %q", string(payload.Data), string(data))
	}

	if err := state.Delete(true); err != nil {
		t.Fatalf("delete: %s", err)
	}

	p, err := state.getStatePayload()
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	if p != nil {
		t.Fatalf("expected empty state, got: %q", string(p.Data))
	}
}

func TestCloudLocks(t *testing.T) {
	back, bCleanup := testBackendWithName(t)
	defer bCleanup()

	a, err := back.StateMgr(testBackendSingleWorkspaceName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	b, err := back.StateMgr(testBackendSingleWorkspaceName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	lockerA, ok := a.(statemgr.Locker)
	if !ok {
		t.Fatal("client A not a statemgr.Locker")
	}

	lockerB, ok := b.(statemgr.Locker)
	if !ok {
		t.Fatal("client B not a statemgr.Locker")
	}

	infoA := statemgr.NewLockInfo()
	infoA.Operation = "test"
	infoA.Who = "clientA"

	infoB := statemgr.NewLockInfo()
	infoB.Operation = "test"
	infoB.Who = "clientB"

	lockIDA, err := lockerA.Lock(infoA)
	if err != nil {
		t.Fatal("unable to get initial lock:", err)
	}

	_, err = lockerB.Lock(infoB)
	if err == nil {
		lockerA.Unlock(lockIDA)
		t.Fatal("client B obtained lock while held by client A")
	}
	if _, ok := err.(*statemgr.LockError); !ok {
		t.Errorf("expected a LockError, but was %t: %s", err, err)
	}

	if err := lockerA.Unlock(lockIDA); err != nil {
		t.Fatal("error unlocking client A", err)
	}

	lockIDB, err := lockerB.Lock(infoB)
	if err != nil {
		t.Fatal("unable to obtain lock from client B")
	}

	if lockIDB == lockIDA {
		t.Fatalf("duplicate lock IDs: %q", lockIDB)
	}

	if err = lockerB.Unlock(lockIDB); err != nil {
		t.Fatal("error unlocking client B:", err)
	}
}

func TestDelete_SafeDeleteNotSupported(t *testing.T) {
	state := testCloudState(t)
	workspaceId := state.workspace.ID
	state.workspace.Permissions.CanForceDelete = nil
	state.workspace.ResourceCount = 5

	// Typically delete(false) should safe-delete a cloud workspace, which should fail on this workspace with resources
	// However, since we have set the workspace canForceDelete permission to nil, we should fall back to force delete
	if err := state.Delete(false); err != nil {
		t.Fatalf("delete: %s", err)
	}
	workspace, err := state.tfeClient.Workspaces.ReadByID(context.Background(), workspaceId)
	if workspace != nil || err != tfe.ErrResourceNotFound {
		t.Fatalf("workspace %s not deleted", workspaceId)
	}
}

func TestDelete_ForceDelete(t *testing.T) {
	state := testCloudState(t)
	workspaceId := state.workspace.ID
	state.workspace.Permissions.CanForceDelete = tfe.Bool(true)
	state.workspace.ResourceCount = 5

	if err := state.Delete(true); err != nil {
		t.Fatalf("delete: %s", err)
	}
	workspace, err := state.tfeClient.Workspaces.ReadByID(context.Background(), workspaceId)
	if workspace != nil || err != tfe.ErrResourceNotFound {
		t.Fatalf("workspace %s not deleted", workspaceId)
	}
}

func TestDelete_SafeDelete(t *testing.T) {
	state := testCloudState(t)
	workspaceId := state.workspace.ID
	state.workspace.Permissions.CanForceDelete = tfe.Bool(false)
	state.workspace.ResourceCount = 5

	// safe-deleting a workspace with resources should fail
	err := state.Delete(false)
	if err == nil {
		t.Fatalf("workspace should have failed to safe delete")
	}

	// safe-deleting a workspace with resources should succeed once it has no resources
	state.workspace.ResourceCount = 0
	if err = state.Delete(false); err != nil {
		t.Fatalf("workspace safe-delete err: %s", err)
	}

	workspace, err := state.tfeClient.Workspaces.ReadByID(context.Background(), workspaceId)
	if workspace != nil || err != tfe.ErrResourceNotFound {
		t.Fatalf("workspace %s not deleted", workspaceId)
	}
}

func TestState_PersistState(t *testing.T) {
	cloudState := testCloudState(t)

	t.Run("Initial PersistState", func(t *testing.T) {
		if cloudState.readState != nil {
			t.Fatal("expected nil initial readState")
		}

		err := cloudState.PersistState(nil)
		if err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		var expectedSerial uint64 = 1
		if cloudState.readSerial != expectedSerial {
			t.Fatalf("expected initial state readSerial to be %d, got %d", expectedSerial, cloudState.readSerial)
		}
	})
}
