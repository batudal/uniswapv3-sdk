package helper

import (
	"errors"
	"math/big"

	"github.com/batudal/uniswapv3-sdk/examples/contract"

	"github.com/batudal/uniswapv3-sdk/constants"
	"github.com/batudal/uniswapv3-sdk/entities"
	sdkutils "github.com/batudal/uniswapv3-sdk/utils"
	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetPoolAddress(client *ethclient.Client, factory common.Address, token0, token1 common.Address, fee *big.Int) (common.Address, error) {
	f, err := contract.NewUniswapv3Factory(factory, client)
	if err != nil {
		return common.Address{}, err
	}
	poolAddr, err := f.GetPool(nil, token0, token1, fee)
	if err != nil {
		return common.Address{}, err
	}
	if poolAddr == (common.Address{}) {
		return common.Address{}, errors.New("pool is not exist")
	}

	return poolAddr, nil
}

func ConstructV3Pool(client *ethclient.Client, factory common.Address, token0, token1 *coreEntities.Token, poolFee uint64) (*entities.Pool, error) {
	poolAddress, err := GetPoolAddress(client, factory, token0.Address, token1.Address, new(big.Int).SetUint64(poolFee))
	if err != nil {
		return nil, err
	}
	contractPool, err := contract.NewPancakev3(poolAddress, client)
	if err != nil {
		return nil, err
	}
	liquidity, err := contractPool.Liquidity(nil)
	if err != nil {
		return nil, err
	}
	slot0, err := contractPool.Slot0(nil)
	if err != nil {
		return nil, err
	}
	pooltick, err := contractPool.Ticks(nil, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	feeAmount := constants.FeeAmount(poolFee)
	ticks := []entities.Tick{
		{
			Index: entities.NearestUsableTick(sdkutils.MinTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
		{
			Index: entities.NearestUsableTick(sdkutils.MaxTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
	}
	// create tick data provider
	ticks[1].LiquidityNet = big.NewInt(0).Sub(big.NewInt(0), pooltick.LiquidityNet)
	p, err := entities.NewTickListDataProvider(ticks, constants.TickSpacings[feeAmount])
	if err != nil {
		return nil, err
	}
	return entities.NewPool(token0, token1, constants.FeeAmount(poolFee),
		slot0.SqrtPriceX96, liquidity, int(slot0.Tick.Int64()), p)
}
