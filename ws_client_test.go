package gobe_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dwdwow/gobe"
)

func TestWsClientPrice(t *testing.T) {
	clt := gobe.NewWsClient(gobe.CHAIN_SOLANA, os.Getenv("BIRDEYE_API_KEY"), nil)
	err := clt.Start()
	if err != nil {
		t.Fatal(err)
	}
	ch := clt.NewDataChan(gobe.WS_PRICE_DATA)
	tokenData := gobe.WsPriceSubData{
		QueryType: gobe.QUERY_TYPE_SIMPLE,
		ChartType: gobe.CHART_1m,
		Address:   "HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		Currency:  gobe.WS_CURRENCY_USD,
	}
	pairData := gobe.WsPriceSubData{
		QueryType: gobe.QUERY_TYPE_SIMPLE,
		ChartType: gobe.CHART_1m,
		Address:   "8sN9549P3Zn6xpQRqpApN57xzkCh6sJxLwuEjcG2W4Ji",
		Currency:  gobe.WS_CURRENCY_PAIR,
	}
	// query := gobe.JoinQuery(tokenData.Query(), pairData.Query())
	// complexData := gobe.WsComplexSubData{
	// 	QueryType: gobe.QUERY_TYPE_COMPLEX,
	// 	Query:     query,
	// }
	// clt.WsSub(gobe.WsSubData[gobe.WsComplexSubData]{
	// 	Type: gobe.SUBSCRIBE_PRICE,
	// 	Data: complexData,
	// })
	clt.WsSub(gobe.WsSubData[gobe.WsPriceSubData]{
		Type: gobe.SUBSCRIBE_PRICE,
		Data: pairData,
	})
	clt.WsSub(gobe.WsSubData[gobe.WsPriceSubData]{
		Type: gobe.SUBSCRIBE_PRICE,
		Data: tokenData,
	})
	for data := range ch {
		fmt.Printf("data: %+v\n", data)
	}
}

func TestWsClientTxs(t *testing.T) {
	clt := gobe.NewWsClient(gobe.CHAIN_SOLANA, os.Getenv("BIRDEYE_API_KEY"), nil)
	err := clt.Start()
	if err != nil {
		t.Fatal(err)
	}
	tokenData := gobe.WsTxsSubData{
		QueryType: gobe.QUERY_TYPE_SIMPLE,
		Address:   "So11111111111111111111111111111111111111112",
	}
	// pairData := gobe.WsTxsSubData{
	// 	QueryType:   gobe.QUERY_TYPE_SIMPLE,
	// 	PairAddress: "8sN9549P3Zn6xpQRqpApN57xzkCh6sJxLwuEjcG2W4Ji",
	// }
	// query := gobe.JoinQuery(tokenData.Query(), pairData.Query())
	// complexData := gobe.WsComplexSubData{
	// 	QueryType: gobe.QUERY_TYPE_COMPLEX,
	// 	Query:     query,
	// }
	// clt.WsSub(gobe.WsSubData[gobe.WsComplexSubData]{
	// 	Type: gobe.SUBSCRIBE_TXS,
	// 	Data: complexData,
	// })
	clt.WsSub(gobe.WsSubData[gobe.WsTxsSubData]{
		Type: gobe.SUBSCRIBE_TXS,
		Data: tokenData,
	})
	ch := clt.NewDataChan(gobe.WS_TXS_DATA)
	for data := range ch {
		fmt.Printf("data: %+v\n", data)
	}
}

func TestWsClientBaseQuotePrice(t *testing.T) {
	clt := gobe.NewWsClient(gobe.CHAIN_SOLANA, os.Getenv("BIRDEYE_API_KEY"), nil)
	err := clt.Start()
	if err != nil {
		t.Fatal(err)
	}
	tokenData := gobe.WsBaseQuotePriceSubData{
		BaseAddress:  "So11111111111111111111111111111111111111112",
		QuoteAddress: "HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		ChartType:    gobe.CHART_1m,
	}
	clt.WsSub(gobe.WsSubData[gobe.WsBaseQuotePriceSubData]{
		Type: gobe.SUBSCRIBE_BASE_QUOTE_PRICE,
		Data: tokenData,
	})
	ch := clt.NewDataChan(gobe.WS_BASE_QUOTE_PRICE_DATA)
	for data := range ch {
		fmt.Printf("data: %+v\n", data)
	}
}

func TestWsClientTokenNewListing(t *testing.T) {
	clt := gobe.NewWsClient(gobe.CHAIN_SOLANA, os.Getenv("BIRDEYE_API_KEY"), nil)
	err := clt.Start()
	if err != nil {
		t.Fatal(err)
	}
	subData := gobe.WsTokenNewListingSubData{
		MemePlatformEnabled: true,
		// MinLiquidity:        100,
		// MaxLiquidity:        1000,
	}
	clt.WsSub(gobe.WsSubData[gobe.WsTokenNewListingSubData]{
		Type: gobe.SUBSCRIBE_TOKEN_NEW_LISTING,
		Data: subData,
	})
	ch := clt.NewDataChan(gobe.WS_TOKEN_NEW_LISTING_DATA)
	for data := range ch {
		fmt.Printf("data: %+v\n", data)
	}
}

func TestWsClientNewPair(t *testing.T) {
	clt := gobe.NewWsClient(gobe.CHAIN_SOLANA, os.Getenv("BIRDEYE_API_KEY"), nil)
	err := clt.Start()
	if err != nil {
		t.Fatal(err)
	}
	subData := gobe.WsNewPairSubData{
		MinLiquidity: 100,
		MaxLiquidity: 1000,
	}
	clt.WsSub(gobe.WsSubData[gobe.WsNewPairSubData]{
		Type: gobe.SUBSCRIBE_NEW_PAIR,
		Data: subData,
	})
	ch := clt.NewDataChan(gobe.WS_NEW_PAIR_DATA)
	for data := range ch {
		fmt.Printf("%s data: %+v\n", time.Now().Format("2006-01-02 15:04:05"), data)
	}
}

func TestWsClientLargeTradeTxs(t *testing.T) {
	clt := gobe.NewWsClient(gobe.CHAIN_SOLANA, os.Getenv("BIRDEYE_API_KEY"), nil)
	err := clt.Start()
	if err != nil {
		t.Fatal(err)
	}
	// subData := gobe.WsLargeTradeTxsSubData{
	// 	MinVolume: 2000,
	// 	MaxVolume: 100000,
	// }
	// clt.WsSub(gobe.WsSubData[gobe.WsLargeTradeTxsSubData]{
	// 	Type: gobe.SUBSCRIBE_LARGE_TRADE_TXS,
	// 	Data: subData,
	// })
	clt.WsSub(gobe.WsLargeTradeTxsSubData{
		Type:      gobe.SUBSCRIBE_LARGE_TRADE_TXS,
		MinVolume: 2000,
		MaxVolume: 100000,
	})
	ch := clt.NewDataChan(gobe.WS_TXS_LARGE_TRADE_DATA)
	for data := range ch {
		fmt.Printf("%s data: %+v\n", time.Now().Format("2006-01-02 15:04:05"), data)
	}
}

func TestWsClientWalletTxs(t *testing.T) {
	clt := gobe.NewWsClient(gobe.CHAIN_SOLANA, os.Getenv("BIRDEYE_API_KEY"), nil)
	err := clt.Start()
	if err != nil {
		t.Fatal(err)
	}
	subData := gobe.WsWalletTxsSubData{
		Address: "7rhxnLV8C77o6d8oz26AgK8x8m5ePsdeRawjqvojbjnQ",
	}
	clt.WsSub(gobe.WsSubData[gobe.WsWalletTxsSubData]{
		Type: gobe.SUBSCRIBE_WALLET_TXS,
		Data: subData,
	})
	ch := clt.NewDataChan(gobe.WS_WALLET_TXS_DATA)
	for data := range ch {
		fmt.Printf("%s data: %+v\n", time.Now().Format("2006-01-02 15:04:05"), data)
	}
}
