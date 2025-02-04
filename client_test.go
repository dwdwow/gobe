package gobe_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dwdwow/gobe"
)

var (
	to   = time.Now()
	from = to.Add(-time.Hour)
)

func TestClientSupportedNetworks(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	networks, err := clt.SupportedNetworks()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("networks: %v\n", networks)
}

func TestClientPrice(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	price, err := clt.Price(gobe.CHAIN_SOLANA, "HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC", true, 1000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("price: %+v\n", price)
}

func TestClientPriceHistory(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	to := time.Now()
	from := to.Add(-time.Hour)
	fmt.Println(from.Unix(), to.Unix())
	history, err := clt.PriceHistory(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		gobe.ADDRESS_TYPE_TOKEN,
		gobe.CHART_1m,
		from.Unix(),
		to.Unix(),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("history: %+v\n", history)
}

func TestClientMultiPrice(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	price, err := clt.MultiPrice(
		gobe.CHAIN_SOLANA,
		[]string{"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC", "So11111111111111111111111111111111111111112"},
		true,
		1000,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("price: %+v\n", price)
}

func TestClientOHLCVByToken(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	fmt.Println(from.Unix(), to.Unix())
	ohlcv, err := clt.OHLCVByToken(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		gobe.CHART_1m,
		from.Unix(),
		to.Unix(),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, ohlcv: %+v\n", len(ohlcv.Items), ohlcv.Items)
}

func TestClientOHLCVByPair(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	ohlcv, err := clt.OHLCVByPair(
		gobe.CHAIN_SOLANA,
		"8sN9549P3Zn6xpQRqpApN57xzkCh6sJxLwuEjcG2W4Ji",
		gobe.CHART_1m,
		from.Unix(),
		to.Unix(),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, ohlcv: %+v\n", len(ohlcv.Items), ohlcv.Items)
}

func TestClientOHLCVByBaseQuote(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	ohlcv, err := clt.OHLCVByBaseQuote(
		gobe.CHAIN_SOLANA,
		"So11111111111111111111111111111111111111112",
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		gobe.CHART_1m,
		from.Unix(),
		to.Unix(),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, ohlcv: %+v\n", len(ohlcv.Items), ohlcv.Items)
}

func TestClientTradesByToken(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	trades, err := clt.TradesByToken(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		gobe.SORT_TYPE_DESC,
		0,
		10,
		gobe.TX_TYPE_ALL,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, trades: %+v\n", len(trades.Items), trades.Items)
}

func TestClientTradesByPair(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	trades, err := clt.TradesByPair(
		gobe.CHAIN_SOLANA,
		"8sN9549P3Zn6xpQRqpApN57xzkCh6sJxLwuEjcG2W4Ji",
		gobe.SORT_TYPE_DESC,
		0,
		10,
		gobe.TX_TYPE_ALL,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, trades: %+v\n", len(trades.Items), trades.Items)
}

func TestClientHistoricalPriceByUnix(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	price, err := clt.HistoricalPriceByUnix(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		time.Now().Unix(),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("price: %+v\n", price)
}

func TestClientPriceVolumeByToken(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	price, err := clt.PriceVolumeByToken(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		gobe.TIME_1h,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("price: %+v\n", price)
}

func TestClientPriceVolumeByTokens(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	price, err := clt.PriceVolumeByTokens(
		gobe.CHAIN_SOLANA,
		[]string{"8sN9549P3Zn6xpQRqpApN57xzkCh6sJxLwuEjcG2W4Ji", "HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC"},
		gobe.TIME_1h,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("price: %+v\n", price)
}

func TestClientTrendingTokens(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	tokens, err := clt.TrendingTokens(
		gobe.CHAIN_SOLANA,
		gobe.RANK_LIQUIDITY,
		gobe.SORT_TYPE_DESC,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("tokens: %+v\n", tokens)
}

func TestTradeByTokenAndTime(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	trades, err := clt.TradeByTokenAndTime(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		0,
		from.Unix(),
		gobe.TX_TYPE_ALL,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, trades: %+v\n", len(trades.Items), trades.Items)
}

func TestTradesSeekByTime(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	trades, err := clt.TradesByPairAndTime(
		gobe.CHAIN_SOLANA,
		"8sN9549P3Zn6xpQRqpApN57xzkCh6sJxLwuEjcG2W4Ji",
		0,
		from.Unix(),
		gobe.TX_TYPE_ALL,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, trades: %+v\n", len(trades.Items), trades.Items)
}

func TestTokenOverview(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	overview, err := clt.TokenOverview(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("overview: %+v\n", overview)
}

func TestTokenList(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	tokens, err := clt.TokenList(
		gobe.CHAIN_SOLANA,
		gobe.SORT_V24HUSD,
		gobe.SORT_TYPE_DESC,
		0,
		10,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, tokens: %+v\n", len(tokens.Items), tokens.Items)
}

func TestTokenListV2(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	url, err := clt.TokenListV2(gobe.CHAIN_SOLANA)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("url: %+v\n", url)
}

func TestTokenSecurity(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	security, err := clt.TokenSecurity(gobe.CHAIN_SOLANA, "HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("security: %+v\n", security)
}

func TestTokenCreationInfo(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	creationInfo, err := clt.TokenCreationInfo(gobe.CHAIN_SOLANA, "HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("creationInfo: %+v\n", creationInfo)
}

func TestMarketList(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	markets, err := clt.MarketList(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		gobe.SORT_LIQUIDITY,
		gobe.SORT_TYPE_DESC,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, markets: %+v\n", len(markets.Items), markets.Items)
}

func TestNewTokenListing(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	listing, err := clt.NewTokenListing(gobe.CHAIN_SOLANA, time.Now().Unix(), 10, false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, listing: %+v\n", len(listing.Items), listing.Items)
}

func TestTokenTopTraders(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	topTraders, err := clt.TokenTopTraders(
		gobe.CHAIN_SOLANA,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		gobe.SORT_VOLUME,
		gobe.SORT_TYPE_DESC,
		gobe.TOP_TRADERS_TIME_24H,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, topTraders: %+v\n", len(topTraders.Items), topTraders.Items)
}

func TestWalletTxHistories(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	txHistories, err := clt.WalletTxHistories(gobe.CHAIN_SOLANA, "C6uq3kFSMDwudLR2ecaWduDYviWEWiJiag7F6A93FyDL", 0, "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("length: %d, txHistories: %+v\n", len(txHistories[gobe.CHAIN_SOLANA]), txHistories)
}

func TestWalletPortfolio(t *testing.T) {
	clt := gobe.NewClient(os.Getenv("BIRDEYE_API_KEY"), gobe.StarterLimiter)
	portfolio, err := clt.WalletPortfolio(gobe.CHAIN_SOLANA, "C6uq3kFSMDwudLR2ecaWduDYviWEWiJiag7F6A93FyDL")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total", portfolio.TotalUsd)
	for _, item := range portfolio.Items {
		fmt.Println(item.Symbol, item.ValueUsd)
	}
}
