package deployment

import (
	"context"
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/datastore"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCCIPViewFiredrill(t *testing.T) {
	lggr := logger.TestLogger(t)
	chains, _ := memory.NewMemoryChains(t, 3, 1)
	env := *deployment.NewEnvironment(
		memory.Memory,
		lggr,
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore[
			datastore.DefaultMetadata,
			datastore.DefaultMetadata,
		]().Seal(),
		chains,
		map[uint64]deployment.SolChain{},
		map[uint64]deployment.AptosChain{},
		[]string{},
		nil,
		func() context.Context { return tests.Context(t) },
		deployment.XXXGenerateTestOCRSecrets(),
	)
	chainSels := env.AllChainSelectors()
	require.Len(t, chainSels, 3)
	chainSel := chainSels[0]
	sourceChainSel := chainSels[1]
	changeset, err := DeployFiredrillContracts(env, FiredrillConfig{
		Version:             deployment.Version1_5_0,
		ChainSelector:       chainSel,
		SourceChainSelector: sourceChainSel,
	})
	require.NoError(t, err)
	err = env.ExistingAddresses.Merge(changeset.AddressBook)
	require.NoError(t, err)
	view, err := CCIPViewFiredrill(env)
	require.NoError(t, err)
	_, err = view.MarshalJSON()
	require.NoError(t, err)
}
