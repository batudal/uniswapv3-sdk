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
	println("poolAddress: ", poolAddress.String())
	contractPool, err := contract.NewUniswapv3Pool(poolAddress, client)
	if err != nil {
		return nil, err
	}
	liquidity, err := contractPool.Liquidity(nil)
	if err != nil {
		return nil, err
	}
	println("liquidity: ", liquidity.String())

	slot0, err := contractPool.Slot0(nil)
	if err != nil {
		return nil, err
	}
	println("slot check")

	pooltick, err := contractPool.Ticks(nil, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	println("pooltick check")
	println("pooltick check")
	println("pooltick check")
	println("pooltick check")
	println("pooltick check")

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
	println("ticks check")

	// create tick data provider
	p, err := entities.NewTickListDataProvider(ticks, constants.TickSpacings[feeAmount])
	if err != nil {
		return nil, err
	}
	println("slot0.SqrtPriceX96: ", slot0.SqrtPriceX96.String())
	return entities.NewPool(token0, token1, constants.FeeAmount(poolFee),
		slot0.SqrtPriceX96, liquidity, int(slot0.Tick.Int64()), p)
}
