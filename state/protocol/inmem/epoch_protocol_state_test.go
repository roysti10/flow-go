package inmem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/model/flow/filter"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/state/protocol/inmem"
	"github.com/onflow/flow-go/state/protocol/mock"
	"github.com/onflow/flow-go/utils/unittest"
)

// TestEpochProtocolStateAdapter tests if the EpochProtocolStateAdapter returns expected values when created
// using constructor passing a RichEpochStateEntry.
func TestEpochProtocolStateAdapter(t *testing.T) {
	// construct a valid protocol state entry that has semantically correct DKGParticipantKeys
	entry := unittest.EpochStateFixture(unittest.WithValidDKG())

	globalParams := mock.NewGlobalParams(t)
	adapter, err := inmem.NewEpochProtocolStateAdapter(
		inmem.UntrustedEpochProtocolStateAdapter{
			RichEpochStateEntry: entry,
			Params:              globalParams,
		},
	)
	require.NoError(t, err)

	t.Run("clustering", func(t *testing.T) {
		clustering, err := inmem.ClusteringFromSetupEvent(entry.CurrentEpochSetup)
		require.NoError(t, err)
		actual, err := adapter.Clustering()
		require.NoError(t, err)
		assert.Equal(t, clustering, actual)
	})
	t.Run("epoch", func(t *testing.T) {
		assert.Equal(t, entry.CurrentEpochSetup.Counter, adapter.Epoch())
	})
	t.Run("setup", func(t *testing.T) {
		assert.Equal(t, entry.CurrentEpochSetup, adapter.EpochSetup())
	})
	t.Run("commit", func(t *testing.T) {
		assert.Equal(t, entry.CurrentEpochCommit, adapter.EpochCommit())
	})
	t.Run("dkg", func(t *testing.T) {
		dkg, err := adapter.DKG()
		require.NoError(t, err)
		assert.Equal(t, entry.CurrentEpochCommit.DKGGroupKey, dkg.GroupKey())
		assert.Equal(t, len(entry.CurrentEpochCommit.DKGParticipantKeys), int(dkg.Size()))
		dkgParticipants := entry.CurrentEpochSetup.Participants.Filter(filter.IsConsensusCommitteeMember)
		for _, identity := range dkgParticipants {
			keyShare, err := dkg.KeyShare(identity.NodeID)
			require.NoError(t, err)
			index, err := dkg.Index(identity.NodeID)
			require.NoError(t, err)
			assert.Equal(t, entry.CurrentEpochCommit.DKGParticipantKeys[index], keyShare)
		}
	})
	t.Run("entry", func(t *testing.T) {
		actualEntry := adapter.Entry()
		assert.Equal(t, entry, actualEntry, "entry should be equal to the one passed to the constructor")
		assert.NotSame(t, entry, actualEntry, "entry should be a copy of the one passed to the constructor")
	})
	t.Run("identities", func(t *testing.T) {
		assert.Equal(t, entry.CurrentEpochIdentityTable, adapter.Identities())
	})
	t.Run("global-params", func(t *testing.T) {
		expectedChainID := flow.Testnet
		globalParams.On("ChainID").Return(expectedChainID, nil).Once()
		actualChainID := adapter.GlobalParams().ChainID()
		assert.Equal(t, expectedChainID, actualChainID)
	})
	t.Run("epoch-phase-staking", func(t *testing.T) {
		entry := unittest.EpochStateFixture()
		adapter, err := inmem.NewEpochProtocolStateAdapter(
			inmem.UntrustedEpochProtocolStateAdapter{
				RichEpochStateEntry: entry,
				Params:              globalParams,
			},
		)
		require.NoError(t, err)
		assert.Equal(t, flow.EpochPhaseStaking, adapter.EpochPhase())
		assert.True(t, adapter.PreviousEpochExists())
		assert.False(t, adapter.EpochFallbackTriggered())
	})
	t.Run("epoch-phase-setup", func(t *testing.T) {
		entry := unittest.EpochStateFixture(unittest.WithNextEpochProtocolState())
		// cleanup the commit event, so we are in setup phase
		entry.NextEpoch.CommitID = flow.ZeroID
		entry.NextEpochCommit = nil

		adapter, err := inmem.NewEpochProtocolStateAdapter(
			inmem.UntrustedEpochProtocolStateAdapter{
				RichEpochStateEntry: entry,
				Params:              globalParams,
			},
		)
		require.NoError(t, err)
		assert.Equal(t, flow.EpochPhaseSetup, adapter.EpochPhase())
		assert.True(t, adapter.PreviousEpochExists())
		assert.False(t, adapter.EpochFallbackTriggered())
	})
	t.Run("epoch-phase-commit", func(t *testing.T) {
		entry := unittest.EpochStateFixture(unittest.WithNextEpochProtocolState())
		adapter, err := inmem.NewEpochProtocolStateAdapter(
			inmem.UntrustedEpochProtocolStateAdapter{
				RichEpochStateEntry: entry,
				Params:              globalParams,
			},
		)
		require.NoError(t, err)
		assert.Equal(t, flow.EpochPhaseCommitted, adapter.EpochPhase())
		assert.True(t, adapter.PreviousEpochExists())
		assert.False(t, adapter.EpochFallbackTriggered())
	})
	t.Run("epoch-fallback-triggered", func(t *testing.T) {
		t.Run("tentatively staking phase", func(t *testing.T) {
			entry := unittest.EpochStateFixture(func(entry *flow.RichEpochStateEntry) {
				entry.EpochFallbackTriggered = true
			})
			adapter, err := inmem.NewEpochProtocolStateAdapter(
				inmem.UntrustedEpochProtocolStateAdapter{
					RichEpochStateEntry: entry,
					Params:              globalParams,
				},
			)
			require.NoError(t, err)
			assert.True(t, adapter.EpochFallbackTriggered())
			assert.Equal(t, flow.EpochPhaseFallback, entry.EpochPhase())
		})
		t.Run("tentatively committed phase", func(t *testing.T) {
			entry := unittest.EpochStateFixture(unittest.WithNextEpochProtocolState(), func(entry *flow.RichEpochStateEntry) {
				entry.EpochFallbackTriggered = true
			})
			adapter, err := inmem.NewEpochProtocolStateAdapter(
				inmem.UntrustedEpochProtocolStateAdapter{
					RichEpochStateEntry: entry,
					Params:              globalParams,
				},
			)
			require.NoError(t, err)
			assert.True(t, adapter.EpochFallbackTriggered())
			assert.Equal(t, flow.EpochPhaseCommitted, entry.EpochPhase())
		})
	})
	t.Run("no-previous-epoch", func(t *testing.T) {
		entry := unittest.EpochStateFixture(func(entry *flow.RichEpochStateEntry) {
			entry.PreviousEpoch = nil
			entry.PreviousEpochSetup = nil
			entry.PreviousEpochCommit = nil
		})
		adapter, err := inmem.NewEpochProtocolStateAdapter(
			inmem.UntrustedEpochProtocolStateAdapter{
				RichEpochStateEntry: entry,
				Params:              globalParams,
			},
		)
		require.NoError(t, err)
		assert.False(t, adapter.PreviousEpochExists())
	})

	// Invalid input with nil Params
	t.Run("invalid - nil Params", func(t *testing.T) {
		_, err := inmem.NewEpochProtocolStateAdapter(
			inmem.UntrustedEpochProtocolStateAdapter{
				RichEpochStateEntry: unittest.EpochStateFixture(),
				Params:              nil,
			},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "params must not be nil")
	})

	// Invalid input with nil RichEpochStateEntry
	t.Run("invalid - nil RichEpochStateEntry", func(t *testing.T) {
		_, err := inmem.NewEpochProtocolStateAdapter(
			inmem.UntrustedEpochProtocolStateAdapter{
				RichEpochStateEntry: nil,
				Params:              globalParams,
			},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "rich epoch state must not be nil")
	})
}
