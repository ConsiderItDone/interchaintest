package penumbra_test

import (
	"context"
	"testing"

	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestAvalancheChainStart(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()
	client, network := interchaintest.DockerSetup(t)

	nv := 5
	nf := 0

	chains, err := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:    "avalanche",
			Version: "v1.9.16",
			ChainConfig: ibc.ChainConfig{
				ChainID: "neto-123123",
				AvalancheSubnets: []ibc.AvalancheSubnetConfig{
					{
						Name:    "timestampvm",
						VMFile:  "",
						Genesis: []byte("{}"),
					},
				},
			},
			NumFullNodes:  &nf,
			NumValidators: &nv,
		},
	},
	).Chains(t.Name())

	require.NoError(t, err, "failed to get avalanche chain")
	require.Len(t, chains, 1)

	chain := chains[0]

	ctx := context.Background()

	err = chain.Initialize(ctx, t.Name(), client, network)
	require.NoError(t, err, "failed to initialize avalanche chain")

	err = chain.Start(t.Name(), ctx)
	require.NoError(t, err, "failed to start avalanche chain")

	err = testutil.WaitForBlocks(ctx, 3, chain)

	require.NoError(t, err, "avalanche chain failed to make blocks")
}