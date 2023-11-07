package penumbra_test

import (
	"context"
	_ "embed"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/errgroup"

	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	subnetevm "github.com/strangelove-ventures/interchaintest/v8/examples/avalanche/subnet-evm"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
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
			Version: "v1.10.1",
			ChainConfig: ibc.ChainConfig{
				ChainID: "neto-123123",
				AvalancheSubnets: []ibc.AvalancheSubnetConfig{
					{
						Name:                "subnetevm",
						VM:                  subnetevm.VM,
						Genesis:             subnetevm.Genesis,
						SubnetClientFactory: subnetevm.NewSubnetEvmClient,
					},
				},
			},
			NumFullNodes:  &nf,
			NumValidators: &nv,
		},
		{
			ChainConfig: ibc.ChainConfig{
				Type:           "cosmos",
				Name:           "ibc-go-simd",
				ChainID:        "simd",
				Bin:            "simd",
				Bech32Prefix:   "cosmos",
				Denom:          "stake",
				GasPrices:      "0.00stake",
				GasAdjustment:  1.3,
				TrustingPeriod: "504h",
			},
		},
	},
	).Chains(t.Name())

	require.NoError(t, err, "failed to get avalanche chain")
	require.Len(t, chains, 1)

	avalanche, simd := chains[0], chains[1]

	ctx := context.Background()

	err = avalanche.Initialize(ctx, t.Name(), client, network)
	require.NoError(t, err, "failed to initialize avalanche chain")

	err = avalanche.Start(t.Name(), ctx)
	require.NoError(t, err, "failed to start avalanche chain")

	subnetCtx := context.WithValue(ctx, "subnet", "0")

	eg := new(errgroup.Group)
	eg.Go(func() error {
		err := avalanche.SendFunds(subnetCtx, "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027", ibc.WalletAmount{
			Address: "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC",
			Amount:  math.NewInt(1000000),
		})
		if err != nil {
			return err
		}
		return avalanche.SendFunds(subnetCtx, "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027", ibc.WalletAmount{
			Address: "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FD",
			Amount:  math.NewInt(2000000),
		})
	})
	eg.Go(func() error {
		return testutil.WaitForBlocks(subnetCtx, 1, chain)
	})

	require.NoError(t, eg.Wait(), "avalanche chain failed to make blocks")
}
