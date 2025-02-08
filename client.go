package gobe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dwdwow/golimiter"
)

const (
	BASE_URL = "https://public-api.birdeye.so"
)

var (
	StandardLimiter = golimiter.NewReqLimiter(time.Second, 1)
	StarterLimiter  = golimiter.NewReqLimiter(time.Second, 15)
	PremiumLimiter  = golimiter.NewReqLimiter(time.Minute, 1000)
	BusinessLimiter = golimiter.NewReqLimiter(time.Minute, 1500)
)

var (
	ErrBadRequest          = fmt.Errorf("invalid request parameters or payload")
	ErrUnauthorized        = fmt.Errorf("authentication required")
	ErrForbidden           = fmt.Errorf("you don't have permission to access this resource")
	ErrUnprocessableEntity = fmt.Errorf("please check the provided data and try again")
	ErrTooManyRequests     = fmt.Errorf("too many requests")
	ErrInternalServer      = fmt.Errorf("something went wrong on the server")
)

var statusCodeToError = map[int]error{
	http.StatusBadRequest:          ErrBadRequest,
	http.StatusUnauthorized:        ErrUnauthorized,
	http.StatusForbidden:           ErrForbidden,
	http.StatusUnprocessableEntity: ErrUnprocessableEntity,
	http.StatusTooManyRequests:     ErrTooManyRequests,
	http.StatusInternalServerError: ErrInternalServer,
}

type ChainType string

const (
	CHAIN_SOLANA    = "solana"
	CHAIN_ETHEREUM  = "ethereum"
	CHAIN_ARBITRUM  = "arbitrum"
	CHAIN_AVALANCHE = "avalanche"
	CHAIN_BSC       = "bsc"
	CHAIN_OPTIMISM  = "optimism"
	CHAIN_POLYGON   = "polygon"
	CHAIN_BASE      = "base"
	CHAIN_ZKSYNC    = "zksync"
	CHAIN_SUI       = "sui"
)

type AddressType string

const (
	ADDRESS_TYPE_TOKEN AddressType = "token"
)

type SortType string

const (
	SORT_TYPE_DESC SortType = "desc"
	SORT_TYPE_ASC  SortType = "asc"
)

type TxType string

const (
	TX_TYPE_SWAP   TxType = "swap"
	TX_TYPE_ADD    TxType = "add"
	TX_TYPE_REMOVE TxType = "remove"
	TX_TYPE_ALL    TxType = "all"
)

type ChartType string

const (
	CHART_1m  ChartType = "1m"
	CHART_3m  ChartType = "3m"
	CHART_5m  ChartType = "5m"
	CHART_15m ChartType = "15m"
	CHART_30m ChartType = "30m"
	CHART_1H  ChartType = "1H"
	CHART_2H  ChartType = "2H"
	CHART_4H  ChartType = "4H"
	CHART_6H  ChartType = "6H"
	CHART_8H  ChartType = "8H"
	CHART_12H ChartType = "12H"
	CHART_1D  ChartType = "1D"
	CHART_3D  ChartType = "3D"
	CHART_1W  ChartType = "1W"
	CHART_1M  ChartType = "1M"
)

type TimeType string

const (
	TIME_1h  TimeType = "1h"
	TIME_2h  TimeType = "2h"
	TIME_4h  TimeType = "4h"
	TIME_8h  TimeType = "8h"
	TIME_24h TimeType = "24h"
)

type RankType string

const (
	RANK_RANK         RankType = "rank"
	RANK_LIQUIDITY    RankType = "liquidity"
	RANK_VOLUME24HUSD RankType = "volume24hUSD"
)

type TokenListSortType string

const (
	SORT_V24HUSD             TokenListSortType = "v24hUSD"
	SORT_MARKET_CAP          TokenListSortType = "mc"
	SORT_V24H_CHANGE_PERCENT TokenListSortType = "v24hChangePercent"
)

type MarketListSortType string

const (
	SORT_LIQUIDITY MarketListSortType = "liquidity"
	SORT_VOLUME24H MarketListSortType = "volume24h"
)

type TokenTopTradersSortType string

const (
	SORT_VOLUME TokenTopTradersSortType = "volume"
	SORT_TRADE  TokenTopTradersSortType = "trade"
)

type TopTradersTimeFrame string

const (
	TOP_TRADERS_TIME_30M TopTradersTimeFrame = "30m"
	TOP_TRADERS_TIME_1H  TopTradersTimeFrame = "1h"
	TOP_TRADERS_TIME_2H  TopTradersTimeFrame = "2h"
	TOP_TRADERS_TIME_4H  TopTradersTimeFrame = "4h"
	TOP_TRADERS_TIME_6H  TopTradersTimeFrame = "6h"
	TOP_TRADERS_TIME_8H  TopTradersTimeFrame = "8h"
	TOP_TRADERS_TIME_12H TopTradersTimeFrame = "12h"
	TOP_TRADERS_TIME_24H TopTradersTimeFrame = "24h"
)

type RespData[D any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    D      `json:"data"`
}

type RespItems[Item any] struct {
	Items   []Item `json:"items"`
	HasNext bool   `json:"hasNext"`
	Total   int64  `json:"total"`
}

type RespNetworks []string

type RespPrice struct {
	Value float64 `json:"value" bson:"value"`
	// UpdateUnixTime seconds
	UpdateUnixTime  int64   `json:"updateUnixTime" bson:"updateUnixTime"`
	UpdateHumanTime string  `json:"updateHumanTime" bson:"updateHumanTime"`
	Liquidity       float64 `json:"liquidity" bson:"liquidity"`
}

type RespPriceHistoryItem struct {
	Address  string  `json:"address" bson:"address"`
	UnixTime int64   `json:"unixTime" bson:"unixTime"`
	Value    float64 `json:"value" bson:"value"`
}

type RespMultiPriceInfo struct {
	Value           float64 `json:"value" bson:"value"`
	UpdateUnixTime  int64   `json:"updateUnixTime" bson:"updateUnixTime"`
	UpdateHumanTime string  `json:"updateHumanTime" bson:"updateHumanTime"`
	PriceChange24h  float64 `json:"priceChange24h" bson:"priceChange24h"`
}

// RespMultiPrice key: address
type RespMultiPrice map[string]RespMultiPriceInfo

type RespOHLCVItem struct {
	Address  string    `json:"address" bson:"address"`
	C        float64   `json:"c" bson:"c"`
	H        float64   `json:"h" bson:"h"`
	L        float64   `json:"l" bson:"l"`
	O        float64   `json:"o" bson:"o"`
	Type     ChartType `json:"type" bson:"type"`
	UnixTime int64     `json:"unixTime" bson:"unixTime"`
	V        float64   `json:"v" bson:"v"`
}

type RespOHLCVBaseQuoteItem struct {
	O float64 `json:"o" bson:"o"`
	C float64 `json:"c" bson:"c"`
	H float64 `json:"h" bson:"h"`
	L float64 `json:"l" bson:"l"`
	// BaseAddress  string    `json:"baseAddress"`
	// QuoteAddress string    `json:"quoteAddress"`
	VBase  float64 `json:"vBase" bson:"vBase"`
	VQuote float64 `json:"vQuote" bson:"vQuote"`
	// Type         ChartType `json:"type"`
	UnixTime int64 `json:"unixTime" bson:"unixTime"`
}

type RespTradesByTokenTokenInfo struct {
	Symbol         string   `json:"symbol" bson:"symbol"`
	Decimals       int64    `json:"decimals" bson:"decimals"`
	Address        string   `json:"address" bson:"address"`
	Amount         int64    `json:"amount" bson:"amount"`
	UiAmount       float64  `json:"uiAmount" bson:"uiAmount"`
	Price          *float64 `json:"price" bson:"price"`
	NearestPrice   float64  `json:"nearestPrice" bson:"nearestPrice"`
	ChangeAmount   int64    `json:"changeAmount" bson:"changeAmount"`
	UiChangeAmount float64  `json:"uiChangeAmount" bson:"uiChangeAmount"`
}

type RespTradesByTokenItem struct {
	Quote         RespTradesByTokenTokenInfo `json:"quote" bson:"quote"`
	Base          RespTradesByTokenTokenInfo `json:"base" bson:"base"`
	BasePrice     *float64                   `json:"basePrice" bson:"basePrice"`
	QuotePrice    *float64                   `json:"quotePrice" bson:"quotePrice"`
	TxHash        string                     `json:"txHash" bson:"txHash"`
	Source        string                     `json:"source" bson:"source"`
	BlockUnixTime int64                      `json:"blockUnixTime" bson:"blockUnixTime"`
	TxType        TxType                     `json:"txType" bson:"txType"`
	Owner         string                     `json:"owner" bson:"owner"`
	Side          string                     `json:"side" bson:"side"`
	Alias         *string                    `json:"alias" bson:"alias"`
	PricePair     float64                    `json:"pricePair" bson:"pricePair"`
	From          RespTradesByTokenTokenInfo `json:"from" bson:"from"`
	To            RespTradesByTokenTokenInfo `json:"to" bson:"to"`
	TokenPrice    *float64                   `json:"tokenPrice" bson:"tokenPrice"`
	PoolId        string                     `json:"poolId" bson:"poolId"`
}

type RespTradesByPairTokenInfo struct {
	Symbol         string   `json:"symbol" bson:"symbol"`
	Decimals       int64    `json:"decimals" bson:"decimals"`
	Address        string   `json:"address" bson:"address"`
	Amount         int64    `json:"amount" bson:"amount"`
	Type           string   `json:"type" bson:"type"`
	TypeSwap       string   `json:"typeSwap" bson:"typeSwap"`
	UiAmount       float64  `json:"uiAmount" bson:"uiAmount"`
	Price          *float64 `json:"price" bson:"price"`
	NearestPrice   float64  `json:"nearestPrice" bson:"nearestPrice"`
	ChangeAmount   int64    `json:"changeAmount" bson:"changeAmount"`
	UiChangeAmount float64  `json:"uiChangeAmount" bson:"uiChangeAmount"`
}

type RespTradesByPairItem struct {
	TxHash        string                    `json:"txHash" bson:"txHash"`
	TxType        TxType                    `json:"txType" bson:"txType"`
	Source        string                    `json:"source" bson:"source"`
	BlockUnixTime int64                     `json:"blockUnixTime" bson:"blockUnixTime"`
	Address       string                    `json:"address" bson:"address"`
	Owner         string                    `json:"owner" bson:"owner"`
	From          RespTradesByPairTokenInfo `json:"from" bson:"from"`
	To            RespTradesByPairTokenInfo `json:"to" bson:"to"`
}

type RespPriceHistoryByTime struct {
	Value          float64 `json:"value" bson:"value"`
	UpdateUnixTime int64   `json:"updateUnixTime" bson:"updateUnixTime"`
	PriceChange24h float64 `json:"priceChange24h" bson:"priceChange24h"`
}

type RespSinglePriceVolume struct {
	// Address just use for multi-token
	Address             string  `json:"address" bson:"address"`
	Price               float64 `json:"price" bson:"price"`
	UpdateUnixTime      int64   `json:"updateUnixTime" bson:"updateUnixTime"`
	UpdateHumanTime     string  `json:"updateHumanTime" bson:"updateHumanTime"`
	VolumeUSD           float64 `json:"volumeUSD" bson:"volumeUSD"`
	VolumeChangePercent float64 `json:"volumeChangePercent" bson:"volumeChangePercent"`
	PriceChangePercent  float64 `json:"priceChangePercent" bson:"priceChangePercent"`
}

type RespTrendingTokensTokenInfo struct {
	Address      string  `json:"address" bson:"address"`
	Decimals     int64   `json:"decimals" bson:"decimals"`
	Liquidity    float64 `json:"liquidity" bson:"liquidity"`
	LogoURI      string  `json:"logoURI" bson:"logoURI"`
	Name         string  `json:"name" bson:"name"`
	Symbol       string  `json:"symbol" bson:"symbol"`
	Volume24hUSD float64 `json:"volume24hUSD" bson:"volume24hUSD"`
	Rank         int64   `json:"rank" bson:"rank"`
}

type RespTrendingTokens struct {
	UpdateUnixTime int64                         `json:"updateUnixTime" bson:"updateUnixTime"`
	UpdateTime     string                        `json:"updateTime" bson:"updateTime"`
	Tokens         []RespTrendingTokensTokenInfo `json:"tokens" bson:"tokens"`
	Total          int64                         `json:"total" bson:"total"`
}

type RespTokenExtension struct {
	CoingeckoId string `json:"coingeckoId" bson:"coingeckoId"`
	SerumV3Usdc string `json:"serumV3Usdc" bson:"serumV3Usdc"`
	SerumV3Usdt string `json:"serumV3Usdt" bson:"serumV3Usdt"`
	Website     string `json:"website" bson:"website"`
	Telegram    string `json:"telegram" bson:"telegram"`
	Twitter     string `json:"twitter" bson:"twitter"`
	Description string `json:"description" bson:"description"`
	Discord     string `json:"discord" bson:"discord"`
	Medium      string `json:"medium" bson:"medium"`
}

type RespTokenOverview struct {
	Address                      string             `json:"address" bson:"address"`
	Decimals                     int64              `json:"decimals" bson:"decimals"`
	Symbol                       string             `json:"symbol" bson:"symbol"`
	Name                         string             `json:"name" bson:"name"`
	Extensions                   RespTokenExtension `json:"extensions" bson:"extensions"`
	LogoURI                      string             `json:"logoURI" bson:"logoURI"`
	Liquidity                    float64            `json:"liquidity" bson:"liquidity"`
	Price                        float64            `json:"price" bson:"price"`
	History30mPrice              float64            `json:"history30mPrice" bson:"history30mPrice"`
	PriceChange30mPercent        float64            `json:"priceChange30mPercent" bson:"priceChange30mPercent"`
	History1hPrice               float64            `json:"history1hPrice" bson:"history1hPrice"`
	PriceChange1hPercent         float64            `json:"priceChange1hPercent" bson:"priceChange1hPercent"`
	History2hPrice               float64            `json:"history2hPrice" bson:"history2hPrice"`
	PriceChange2hPercent         float64            `json:"priceChange2hPercent" bson:"priceChange2hPercent"`
	History4hPrice               float64            `json:"history4hPrice" bson:"history4hPrice"`
	PriceChange4hPercent         float64            `json:"priceChange4hPercent" bson:"priceChange4hPercent"`
	History6hPrice               float64            `json:"history6hPrice" bson:"history6hPrice"`
	PriceChange6hPercent         float64            `json:"priceChange6hPercent" bson:"priceChange6hPercent"`
	History8hPrice               float64            `json:"history8hPrice" bson:"history8hPrice"`
	PriceChange8hPercent         float64            `json:"priceChange8hPercent" bson:"priceChange8hPercent"`
	History12hPrice              float64            `json:"history12hPrice" bson:"history12hPrice"`
	PriceChange12hPercent        float64            `json:"priceChange12hPercent" bson:"priceChange12hPercent"`
	History24hPrice              float64            `json:"history24hPrice" bson:"history24hPrice"`
	PriceChange24hPercent        float64            `json:"priceChange24hPercent" bson:"priceChange24hPercent"`
	UniqueWallet30m              int64              `json:"uniqueWallet30m" bson:"uniqueWallet30m"`
	UniqueWalletHistory30m       int64              `json:"uniqueWalletHistory30m" bson:"uniqueWalletHistory30m"`
	UniqueWallet30mChangePercent float64            `json:"uniqueWallet30mChangePercent" bson:"uniqueWallet30mChangePercent"`
	UniqueWallet1h               int64              `json:"uniqueWallet1h" bson:"uniqueWallet1h"`
	UniqueWalletHistory1h        int64              `json:"uniqueWalletHistory1h" bson:"uniqueWalletHistory1h"`
	UniqueWallet1hChangePercent  float64            `json:"uniqueWallet1hChangePercent" bson:"uniqueWallet1hChangePercent"`
	UniqueWallet2h               int64              `json:"uniqueWallet2h" bson:"uniqueWallet2h"`
	UniqueWalletHistory2h        int64              `json:"uniqueWalletHistory2h" bson:"uniqueWalletHistory2h"`
	UniqueWallet2hChangePercent  float64            `json:"uniqueWallet2hChangePercent" bson:"uniqueWallet2hChangePercent"`
	UniqueWallet4h               int64              `json:"uniqueWallet4h" bson:"uniqueWallet4h"`
	UniqueWalletHistory4h        int64              `json:"uniqueWalletHistory4h" bson:"uniqueWalletHistory4h"`
	UniqueWallet4hChangePercent  float64            `json:"uniqueWallet4hChangePercent" bson:"uniqueWallet4hChangePercent"`
	UniqueWallet6h               int64              `json:"uniqueWallet6h" bson:"uniqueWallet6h"`
	UniqueWalletHistory6h        int64              `json:"uniqueWalletHistory6h" bson:"uniqueWalletHistory6h"`
	UniqueWallet6hChangePercent  float64            `json:"uniqueWallet6hChangePercent" bson:"uniqueWallet6hChangePercent"`
	UniqueWallet8h               int64              `json:"uniqueWallet8h" bson:"uniqueWallet8h"`
	UniqueWalletHistory8h        int64              `json:"uniqueWalletHistory8h" bson:"uniqueWalletHistory8h"`
	UniqueWallet8hChangePercent  float64            `json:"uniqueWallet8hChangePercent" bson:"uniqueWallet8hChangePercent"`
	UniqueWallet12h              int64              `json:"uniqueWallet12h" bson:"uniqueWallet12h"`
	UniqueWalletHistory12h       int64              `json:"uniqueWalletHistory12h" bson:"uniqueWalletHistory12h"`
	UniqueWallet12hChangePercent float64            `json:"uniqueWallet12hChangePercent" bson:"uniqueWallet12hChangePercent"`
	UniqueWallet24h              int64              `json:"uniqueWallet24h" bson:"uniqueWallet24h"`
	UniqueWalletHistory24h       int64              `json:"uniqueWalletHistory24h" bson:"uniqueWalletHistory24h"`
	UniqueWallet24hChangePercent float64            `json:"uniqueWallet24hChangePercent" bson:"uniqueWallet24hChangePercent"`
	LastTradeUnixTime            int64              `json:"lastTradeUnixTime" bson:"lastTradeUnixTime"`
	LastTradeHumanTime           string             `json:"lastTradeHumanTime" bson:"lastTradeHumanTime"`
	Supply                       float64            `json:"supply" bson:"supply"`
	Mc                           float64            `json:"mc" bson:"mc"`
	Trade30m                     int64              `json:"trade30m" bson:"trade30m"`
	TradeHistory30m              int64              `json:"tradeHistory30m" bson:"tradeHistory30m"`
	Trade30mChangePercent        float64            `json:"trade30mChangePercent" bson:"trade30mChangePercent"`
	Sell30m                      int64              `json:"sell30m" bson:"sell30m"`
	SellHistory30m               int64              `json:"sellHistory30m" bson:"sellHistory30m"`
	Sell30mChangePercent         float64            `json:"sell30mChangePercent" bson:"sell30mChangePercent"`
	Buy30m                       int64              `json:"buy30m" bson:"buy30m"`
	BuyHistory30m                int64              `json:"buyHistory30m" bson:"buyHistory30m"`
	Buy30mChangePercent          float64            `json:"buy30mChangePercent" bson:"buy30mChangePercent"`
	V30m                         float64            `json:"v30m" bson:"v30m"`
	V30mUSD                      float64            `json:"v30mUSD" bson:"v30mUSD"`
	VHistory30m                  float64            `json:"vHistory30m" bson:"vHistory30m"`
	VHistory30mUSD               float64            `json:"vHistory30mUSD" bson:"vHistory30mUSD"`
	V30mChangePercent            float64            `json:"v30mChangePercent" bson:"v30mChangePercent"`
	VBuy30m                      float64            `json:"vBuy30m" bson:"vBuy30m"`
	VBuy30mUSD                   float64            `json:"vBuy30mUSD" bson:"vBuy30mUSD"`
	VBuyHistory30m               float64            `json:"vBuyHistory30m" bson:"vBuyHistory30m"`
	VBuyHistory30mUSD            float64            `json:"vBuyHistory30mUSD" bson:"vBuyHistory30mUSD"`
	VBuy30mChangePercent         float64            `json:"vBuy30mChangePercent" bson:"vBuy30mChangePercent"`
	VSell30m                     float64            `json:"vSell30m" bson:"vSell30m"`
	VSell30mUSD                  float64            `json:"vSell30mUSD" bson:"vSell30mUSD"`
	VSellHistory30m              float64            `json:"vSellHistory30m" bson:"vSellHistory30m"`
	VSellHistory30mUSD           float64            `json:"vSellHistory30mUSD" bson:"vSellHistory30mUSD"`
	VSell30mChangePercent        float64            `json:"vSell30mChangePercent" bson:"vSell30mChangePercent"`
	Trade1h                      int64              `json:"trade1h" bson:"trade1h"`
	TradeHistory1h               int64              `json:"tradeHistory1h" bson:"tradeHistory1h"`
	Trade1hChangePercent         float64            `json:"trade1hChangePercent" bson:"trade1hChangePercent"`
	Sell1h                       int64              `json:"sell1h" bson:"sell1h"`
	SellHistory1h                int64              `json:"sellHistory1h" bson:"sellHistory1h"`
	Sell1hChangePercent          float64            `json:"sell1hChangePercent" bson:"sell1hChangePercent"`
	Buy1h                        int64              `json:"buy1h" bson:"buy1h"`
	BuyHistory1h                 int64              `json:"buyHistory1h" bson:"buyHistory1h"`
	Buy1hChangePercent           float64            `json:"buy1hChangePercent" bson:"buy1hChangePercent"`
	V1h                          float64            `json:"v1h" bson:"v1h"`
	V1hUSD                       float64            `json:"v1hUSD" bson:"v1hUSD"`
	VHistory1h                   float64            `json:"vHistory1h" bson:"vHistory1h"`
	VHistory1hUSD                float64            `json:"vHistory1hUSD" bson:"vHistory1hUSD"`
	V1hChangePercent             float64            `json:"v1hChangePercent" bson:"v1hChangePercent"`
	VBuy1h                       float64            `json:"vBuy1h" bson:"vBuy1h"`
	VBuy1hUSD                    float64            `json:"vBuy1hUSD" bson:"vBuy1hUSD"`
	VBuyHistory1h                float64            `json:"vBuyHistory1h" bson:"vBuyHistory1h"`
	VBuyHistory1hUSD             float64            `json:"vBuyHistory1hUSD" bson:"vBuyHistory1hUSD"`
	VBuy1hChangePercent          float64            `json:"vBuy1hChangePercent" bson:"vBuy1hChangePercent"`
	VSell1h                      float64            `json:"vSell1h" bson:"vSell1h"`
	VSell1hUSD                   float64            `json:"vSell1hUSD" bson:"vSell1hUSD"`
	VSellHistory1h               float64            `json:"vSellHistory1h" bson:"vSellHistory1h"`
	VSellHistory1hUSD            float64            `json:"vSellHistory1hUSD" bson:"vSellHistory1hUSD"`
	VSell1hChangePercent         float64            `json:"vSell1hChangePercent" bson:"vSell1hChangePercent"`
	Trade2h                      int64              `json:"trade2h" bson:"trade2h"`
	TradeHistory2h               int64              `json:"tradeHistory2h" bson:"tradeHistory2h"`
	Trade2hChangePercent         float64            `json:"trade2hChangePercent" bson:"trade2hChangePercent"`
	Sell2h                       int64              `json:"sell2h" bson:"sell2h"`
	SellHistory2h                int64              `json:"sellHistory2h" bson:"sellHistory2h"`
	Sell2hChangePercent          float64            `json:"sell2hChangePercent" bson:"sell2hChangePercent"`
	Buy2h                        int64              `json:"buy2h" bson:"buy2h"`
	BuyHistory2h                 int64              `json:"buyHistory2h" bson:"buyHistory2h"`
	Buy2hChangePercent           float64            `json:"buy2hChangePercent" bson:"buy2hChangePercent"`
	V2h                          float64            `json:"v2h" bson:"v2h"`
	V2hUSD                       float64            `json:"v2hUSD" bson:"v2hUSD"`
	VHistory2h                   float64            `json:"vHistory2h" bson:"vHistory2h"`
	VHistory2hUSD                float64            `json:"vHistory2hUSD" bson:"vHistory2hUSD"`
	V2hChangePercent             float64            `json:"v2hChangePercent" bson:"v2hChangePercent"`
	VBuy2h                       float64            `json:"vBuy2h" bson:"vBuy2h"`
	VBuy2hUSD                    float64            `json:"vBuy2hUSD" bson:"vBuy2hUSD"`
	VBuyHistory2h                float64            `json:"vBuyHistory2h" bson:"vBuyHistory2h"`
	VBuyHistory2hUSD             float64            `json:"vBuyHistory2hUSD" bson:"vBuyHistory2hUSD"`
	VBuy2hChangePercent          float64            `json:"vBuy2hChangePercent" bson:"vBuy2hChangePercent"`
	VSell2h                      float64            `json:"vSell2h" bson:"vSell2h"`
	VSell2hUSD                   float64            `json:"vSell2hUSD" bson:"vSell2hUSD"`
	VSellHistory2h               float64            `json:"vSellHistory2h" bson:"vSellHistory2h"`
	VSellHistory2hUSD            float64            `json:"vSellHistory2hUSD" bson:"vSellHistory2hUSD"`
	VSell2hChangePercent         float64            `json:"vSell2hChangePercent" bson:"vSell2hChangePercent"`
	Trade4h                      int64              `json:"trade4h" bson:"trade4h"`
	TradeHistory4h               int64              `json:"tradeHistory4h" bson:"tradeHistory4h"`
	Trade4hChangePercent         float64            `json:"trade4hChangePercent" bson:"trade4hChangePercent"`
	Sell4h                       int64              `json:"sell4h" bson:"sell4h"`
	SellHistory4h                int64              `json:"sellHistory4h" bson:"sellHistory4h"`
	Sell4hChangePercent          float64            `json:"sell4hChangePercent" bson:"sell4hChangePercent"`
	Buy4h                        int64              `json:"buy4h" bson:"buy4h"`
	BuyHistory4h                 int64              `json:"buyHistory4h" bson:"buyHistory4h"`
	Buy4hChangePercent           float64            `json:"buy4hChangePercent" bson:"buy4hChangePercent"`
	V4h                          float64            `json:"v4h" bson:"v4h"`
	V4hUSD                       float64            `json:"v4hUSD" bson:"v4hUSD"`
	VHistory4h                   float64            `json:"vHistory4h" bson:"vHistory4h"`
	VHistory4hUSD                float64            `json:"vHistory4hUSD" bson:"vHistory4hUSD"`
	V4hChangePercent             float64            `json:"v4hChangePercent" bson:"v4hChangePercent"`
	VBuy4h                       float64            `json:"vBuy4h" bson:"vBuy4h"`
	VBuy4hUSD                    float64            `json:"vBuy4hUSD" bson:"vBuy4hUSD"`
	VBuyHistory4h                float64            `json:"vBuyHistory4h" bson:"vBuyHistory4h"`
	VBuyHistory4hUSD             float64            `json:"vBuyHistory4hUSD" bson:"vBuyHistory4hUSD"`
	VBuy4hChangePercent          float64            `json:"vBuy4hChangePercent" bson:"vBuy4hChangePercent"`
	VSell4h                      float64            `json:"vSell4h" bson:"vSell4h"`
	VSell4hUSD                   float64            `json:"vSell4hUSD" bson:"vSell4hUSD"`
	VSellHistory4h               float64            `json:"vSellHistory4h" bson:"vSellHistory4h"`
	VSellHistory4hUSD            float64            `json:"vSellHistory4hUSD" bson:"vSellHistory4hUSD"`
	VSell4hChangePercent         float64            `json:"vSell4hChangePercent" bson:"vSell4hChangePercent"`
	Trade6h                      int64              `json:"trade6h" bson:"trade6h"`
	TradeHistory6h               int64              `json:"tradeHistory6h" bson:"tradeHistory6h"`
	Trade6hChangePercent         float64            `json:"trade6hChangePercent" bson:"trade6hChangePercent"`
	Sell6h                       int64              `json:"sell6h" bson:"sell6h"`
	SellHistory6h                int64              `json:"sellHistory6h" bson:"sellHistory6h"`
	Sell6hChangePercent          float64            `json:"sell6hChangePercent" bson:"sell6hChangePercent"`
	Buy6h                        int64              `json:"buy6h" bson:"buy6h"`
	BuyHistory6h                 int64              `json:"buyHistory6h" bson:"buyHistory6h"`
	Buy6hChangePercent           float64            `json:"buy6hChangePercent" bson:"buy6hChangePercent"`
	V6h                          float64            `json:"v6h" bson:"v6h"`
	V6hUSD                       float64            `json:"v6hUSD" bson:"v6hUSD"`
	VHistory6h                   float64            `json:"vHistory6h" bson:"vHistory6h"`
	VHistory6hUSD                float64            `json:"vHistory6hUSD" bson:"vHistory6hUSD"`
	V6hChangePercent             float64            `json:"v6hChangePercent" bson:"v6hChangePercent"`
	VBuy6h                       float64            `json:"vBuy6h" bson:"vBuy6h"`
	VBuy6hUSD                    float64            `json:"vBuy6hUSD" bson:"vBuy6hUSD"`
	VBuyHistory6h                float64            `json:"vBuyHistory6h" bson:"vBuyHistory6h"`
	VBuyHistory6hUSD             float64            `json:"vBuyHistory6hUSD" bson:"vBuyHistory6hUSD"`
	VBuy6hChangePercent          float64            `json:"vBuy6hChangePercent" bson:"vBuy6hChangePercent"`
	VSell6h                      float64            `json:"vSell6h" bson:"vSell6h"`
	VSell6hUSD                   float64            `json:"vSell6hUSD" bson:"vSell6hUSD"`
	VSellHistory6h               float64            `json:"vSellHistory6h" bson:"vSellHistory6h"`
	VSellHistory6hUSD            float64            `json:"vSellHistory6hUSD" bson:"vSellHistory6hUSD"`
	VSell6hChangePercent         float64            `json:"vSell6hChangePercent" bson:"vSell6hChangePercent"`
	Trade8h                      int64              `json:"trade8h" bson:"trade8h"`
	TradeHistory8h               int64              `json:"tradeHistory8h" bson:"tradeHistory8h"`
	Trade8hChangePercent         float64            `json:"trade8hChangePercent" bson:"trade8hChangePercent"`
	Sell8h                       int64              `json:"sell8h" bson:"sell8h"`
	SellHistory8h                int64              `json:"sellHistory8h" bson:"sellHistory8h"`
	Sell8hChangePercent          float64            `json:"sell8hChangePercent" bson:"sell8hChangePercent"`
	Buy8h                        int64              `json:"buy8h" bson:"buy8h"`
	BuyHistory8h                 int64              `json:"buyHistory8h" bson:"buyHistory8h"`
	Buy8hChangePercent           float64            `json:"buy8hChangePercent" bson:"buy8hChangePercent"`
	V8h                          float64            `json:"v8h" bson:"v8h"`
	V8hUSD                       float64            `json:"v8hUSD" bson:"v8hUSD"`
	VHistory8h                   float64            `json:"vHistory8h" bson:"vHistory8h"`
	VHistory8hUSD                float64            `json:"vHistory8hUSD" bson:"vHistory8hUSD"`
	V8hChangePercent             float64            `json:"v8hChangePercent" bson:"v8hChangePercent"`
	VBuy8h                       float64            `json:"vBuy8h" bson:"vBuy8h"`
	VBuy8hUSD                    float64            `json:"vBuy8hUSD" bson:"vBuy8hUSD"`
	VBuyHistory8h                float64            `json:"vBuyHistory8h" bson:"vBuyHistory8h"`
	VBuyHistory8hUSD             float64            `json:"vBuyHistory8hUSD" bson:"vBuyHistory8hUSD"`
	VBuy8hChangePercent          float64            `json:"vBuy8hChangePercent" bson:"vBuy8hChangePercent"`
	VSell8h                      float64            `json:"vSell8h" bson:"vSell8h"`
	VSell8hUSD                   float64            `json:"vSell8hUSD" bson:"vSell8hUSD"`
	VSellHistory8h               float64            `json:"vSellHistory8h" bson:"vSellHistory8h"`
	VSellHistory8hUSD            float64            `json:"vSellHistory8hUSD" bson:"vSellHistory8hUSD"`
	VSell8hChangePercent         float64            `json:"vSell8hChangePercent" bson:"vSell8hChangePercent"`
	Trade12h                     int64              `json:"trade12h" bson:"trade12h"`
	TradeHistory12h              int64              `json:"tradeHistory12h" bson:"tradeHistory12h"`
	Trade12hChangePercent        float64            `json:"trade12hChangePercent" bson:"trade12hChangePercent"`
	Sell12h                      int64              `json:"sell12h" bson:"sell12h"`
	SellHistory12h               int64              `json:"sellHistory12h" bson:"sellHistory12h"`
	Sell12hChangePercent         float64            `json:"sell12hChangePercent" bson:"sell12hChangePercent"`
	Buy12h                       int64              `json:"buy12h" bson:"buy12h"`
	BuyHistory12h                int64              `json:"buyHistory12h" bson:"buyHistory12h"`
	Buy12hChangePercent          float64            `json:"buy12hChangePercent" bson:"buy12hChangePercent"`
	V12h                         float64            `json:"v12h" bson:"v12h"`
	V12hUSD                      float64            `json:"v12hUSD" bson:"v12hUSD"`
	VHistory12h                  float64            `json:"vHistory12h" bson:"vHistory12h"`
	VHistory12hUSD               float64            `json:"vHistory12hUSD" bson:"vHistory12hUSD"`
	V12hChangePercent            float64            `json:"v12hChangePercent" bson:"v12hChangePercent"`
	VBuy12h                      float64            `json:"vBuy12h" bson:"vBuy12h"`
	VBuy12hUSD                   float64            `json:"vBuy12hUSD" bson:"vBuy12hUSD"`
	VBuyHistory12h               float64            `json:"vBuyHistory12h" bson:"vBuyHistory12h"`
	VBuyHistory12hUSD            float64            `json:"vBuyHistory12hUSD" bson:"vBuyHistory12hUSD"`
	VBuy12hChangePercent         float64            `json:"vBuy12hChangePercent" bson:"vBuy12hChangePercent"`
	VSell12h                     float64            `json:"vSell12h" bson:"vSell12h"`
	VSell12hUSD                  float64            `json:"vSell12hUSD" bson:"vSell12hUSD"`
	VSellHistory12h              float64            `json:"vSellHistory12h" bson:"vSellHistory12h"`
	VSellHistory12hUSD           float64            `json:"vSellHistory12hUSD" bson:"vSellHistory12hUSD"`
	VSell12hChangePercent        float64            `json:"vSell12hChangePercent" bson:"vSell12hChangePercent"`
	Trade24h                     int64              `json:"trade24h" bson:"trade24h"`
	TradeHistory24h              int64              `json:"tradeHistory24h" bson:"tradeHistory24h"`
	Trade24hChangePercent        float64            `json:"trade24hChangePercent" bson:"trade24hChangePercent"`
	Sell24h                      int64              `json:"sell24h" bson:"sell24h"`
	SellHistory24h               int64              `json:"sellHistory24h" bson:"sellHistory24h"`
	Sell24hChangePercent         float64            `json:"sell24hChangePercent" bson:"sell24hChangePercent"`
	Buy24h                       int64              `json:"buy24h" bson:"buy24h"`
	BuyHistory24h                int64              `json:"buyHistory24h" bson:"buyHistory24h"`
	Buy24hChangePercent          float64            `json:"buy24hChangePercent" bson:"buy24hChangePercent"`
	V24h                         float64            `json:"v24h" bson:"v24h"`
	V24hUSD                      float64            `json:"v24hUSD" bson:"v24hUSD"`
	VHistory24h                  float64            `json:"vHistory24h" bson:"vHistory24h"`
	VHistory24hUSD               float64            `json:"vHistory24hUSD" bson:"vHistory24hUSD"`
	V24hChangePercent            float64            `json:"v24hChangePercent" bson:"v24hChangePercent"`
	VBuy24h                      float64            `json:"vBuy24h" bson:"vBuy24h"`
	VBuy24hUSD                   float64            `json:"vBuy24hUSD" bson:"vBuy24hUSD"`
	VBuyHistory24h               float64            `json:"vBuyHistory24h" bson:"vBuyHistory24h"`
	VBuyHistory24hUSD            float64            `json:"vBuyHistory24hUSD" bson:"vBuyHistory24hUSD"`
	VBuy24hChangePercent         float64            `json:"vBuy24hChangePercent" bson:"vBuy24hChangePercent"`
	VSell24h                     float64            `json:"vSell24h" bson:"vSell24h"`
	VSell24hUSD                  float64            `json:"vSell24hUSD" bson:"vSell24hUSD"`
	VSellHistory24h              float64            `json:"vSellHistory24h" bson:"vSellHistory24h"`
	VSellHistory24hUSD           float64            `json:"vSellHistory24hUSD" bson:"vSellHistory24hUSD"`
	VSell24hChangePercent        float64            `json:"vSell24hChangePercent" bson:"vSell24hChangePercent"`
	Watch                        int64              `json:"watch" bson:"watch"`
	View30m                      int64              `json:"view30m" bson:"view30m"`
	ViewHistory30m               int64              `json:"viewHistory30m" bson:"viewHistory30m"`
	View30mChangePercent         float64            `json:"view30mChangePercent" bson:"view30mChangePercent"`
	View1h                       int64              `json:"view1h" bson:"view1h"`
	ViewHistory1h                int64              `json:"viewHistory1h" bson:"viewHistory1h"`
	View1hChangePercent          float64            `json:"view1hChangePercent" bson:"view1hChangePercent"`
	View2h                       int64              `json:"view2h" bson:"view2h"`
	ViewHistory2h                int64              `json:"viewHistory2h" bson:"viewHistory2h"`
	View2hChangePercent          float64            `json:"view2hChangePercent" bson:"view2hChangePercent"`
	View4h                       int64              `json:"view4h" bson:"view4h"`
	ViewHistory4h                int64              `json:"viewHistory4h" bson:"viewHistory4h"`
	View4hChangePercent          float64            `json:"view4hChangePercent" bson:"view4hChangePercent"`
	View6h                       int64              `json:"view6h" bson:"view6h"`
	ViewHistory6h                int64              `json:"viewHistory6h" bson:"viewHistory6h"`
	View6hChangePercent          float64            `json:"view6hChangePercent" bson:"view6hChangePercent"`
	View8h                       int64              `json:"view8h" bson:"view8h"`
	ViewHistory8h                int64              `json:"viewHistory8h" bson:"viewHistory8h"`
	View8hChangePercent          float64            `json:"view8hChangePercent" bson:"view8hChangePercent"`
	View12h                      int64              `json:"view12h" bson:"view12h"`
	ViewHistory12h               int64              `json:"viewHistory12h" bson:"viewHistory12h"`
	View12hChangePercent         float64            `json:"view12hChangePercent" bson:"view12hChangePercent"`
	View24h                      int64              `json:"view24h" bson:"view24h"`
	ViewHistory24h               int64              `json:"viewHistory24h" bson:"viewHistory24h"`
	View24hChangePercent         float64            `json:"view24hChangePercent" bson:"view24hChangePercent"`
	UniqueView30m                int64              `json:"uniqueView30m" bson:"uniqueView30m"`
	UniqueViewHistory30m         int64              `json:"uniqueViewHistory30m" bson:"uniqueViewHistory30m"`
	UniqueView30mChangePercent   float64            `json:"uniqueView30mChangePercent" bson:"uniqueView30mChangePercent"`
	UniqueView1h                 int64              `json:"uniqueView1h" bson:"uniqueView1h"`
	UniqueViewHistory1h          int64              `json:"uniqueViewHistory1h" bson:"uniqueViewHistory1h"`
	UniqueView1hChangePercent    float64            `json:"uniqueView1hChangePercent" bson:"uniqueView1hChangePercent"`
	UniqueView2h                 int64              `json:"uniqueView2h" bson:"uniqueView2h"`
	UniqueViewHistory2h          int64              `json:"uniqueViewHistory2h" bson:"uniqueViewHistory2h"`
	UniqueView2hChangePercent    float64            `json:"uniqueView2hChangePercent" bson:"uniqueView2hChangePercent"`
	UniqueView4h                 int64              `json:"uniqueView4h" bson:"uniqueView4h"`
	UniqueViewHistory4h          int64              `json:"uniqueViewHistory4h" bson:"uniqueViewHistory4h"`
	UniqueView4hChangePercent    float64            `json:"uniqueView4hChangePercent" bson:"uniqueView4hChangePercent"`
	UniqueView6h                 int64              `json:"uniqueView6h" bson:"uniqueView6h"`
	UniqueViewHistory6h          int64              `json:"uniqueViewHistory6h" bson:"uniqueViewHistory6h"`
	UniqueView6hChangePercent    float64            `json:"uniqueView6hChangePercent" bson:"uniqueView6hChangePercent"`
	UniqueView8h                 int64              `json:"uniqueView8h" bson:"uniqueView8h"`
	UniqueViewHistory8h          int64              `json:"uniqueViewHistory8h" bson:"uniqueViewHistory8h"`
	UniqueView8hChangePercent    float64            `json:"uniqueView8hChangePercent" bson:"uniqueView8hChangePercent"`
	UniqueView12h                int64              `json:"uniqueView12h" bson:"uniqueView12h"`
	UniqueViewHistory12h         int64              `json:"uniqueViewHistory12h" bson:"uniqueViewHistory12h"`
	UniqueView12hChangePercent   float64            `json:"uniqueView12hChangePercent" bson:"uniqueView12hChangePercent"`
	UniqueView24h                int64              `json:"uniqueView24h" bson:"uniqueView24h"`
	UniqueViewHistory24h         int64              `json:"uniqueViewHistory24h" bson:"uniqueViewHistory24h"`
	UniqueView24hChangePercent   float64            `json:"uniqueView24hChangePercent" bson:"uniqueView24hChangePercent"`
	NumberMarkets                int64              `json:"numberMarkets" bson:"numberMarkets"`
}

type RespToken struct {
	Address           string  `json:"address" bson:"address"`
	Decimals          int64   `json:"decimals" bson:"decimals"`
	Liquidity         float64 `json:"liquidity" bson:"liquidity"`
	Mc                float64 `json:"mc" bson:"mc"`
	Symbol            string  `json:"symbol" bson:"symbol"`
	V24hChangePercent float64 `json:"v24hChangePercent" bson:"v24hChangePercent"`
	V24hUSD           float64 `json:"v24hUSD" bson:"v24hUSD"`
	Name              string  `json:"name" bson:"name"`
	LastTradeUnixTime int64   `json:"lastTradeUnixTime" bson:"lastTradeUnixTime"`
}

type RespTokenListV2Url struct {
	Url string `json:"url" bson:"url"`
}

type RespTokenSecurity struct {
	CreatorAddress                 *string  `json:"creatorAddress" bson:"creatorAddress"`
	OwnerAddress                   *string  `json:"ownerAddress" bson:"ownerAddress"`
	CreationTx                     *string  `json:"creationTx" bson:"creationTx"`
	CreationTime                   *int64   `json:"creationTime" bson:"creationTime"`
	CreationSlot                   *int64   `json:"creationSlot" bson:"creationSlot"`
	MintTx                         *string  `json:"mintTx" bson:"mintTx"`
	MintTime                       *int64   `json:"mintTime" bson:"mintTime"`
	MintSlot                       *int64   `json:"mintSlot" bson:"mintSlot"`
	CreatorBalance                 *float64 `json:"creatorBalance" bson:"creatorBalance"`
	OwnerBalance                   *float64 `json:"ownerBalance" bson:"ownerBalance"`
	OwnerPercentage                *float64 `json:"ownerPercentage" bson:"ownerPercentage"`
	CreatorPercentage              *float64 `json:"creatorPercentage" bson:"creatorPercentage"`
	MetaplexUpdateAuthority        string   `json:"metaplexUpdateAuthority" bson:"metaplexUpdateAuthority"`
	MetaplexUpdateAuthorityBalance float64  `json:"metaplexUpdateAuthorityBalance" bson:"metaplexUpdateAuthorityBalance"`
	MetaplexUpdateAuthorityPercent float64  `json:"metaplexUpdateAuthorityPercent" bson:"metaplexUpdateAuthorityPercent"`
	MutableMetadata                bool     `json:"mutableMetadata" bson:"mutableMetadata"`
	Top10HolderBalance             float64  `json:"top10HolderBalance" bson:"top10HolderBalance"`
	Top10HolderPercent             float64  `json:"top10HolderPercent" bson:"top10HolderPercent"`
	Top10UserBalance               float64  `json:"top10UserBalance" bson:"top10UserBalance"`
	Top10UserPercent               float64  `json:"top10UserPercent" bson:"top10UserPercent"`
	IsTrueToken                    *bool    `json:"isTrueToken" bson:"isTrueToken"`
	TotalSupply                    float64  `json:"totalSupply" bson:"totalSupply"`
	PreMarketHolder                []string `json:"preMarketHolder" bson:"preMarketHolder"`
	LockInfo                       *string  `json:"lockInfo" bson:"lockInfo"`
	Freezeable                     *bool    `json:"freezeable" bson:"freezeable"`
	FreezeAuthority                *string  `json:"freezeAuthority" bson:"freezeAuthority"`
	TransferFeeEnable              *bool    `json:"transferFeeEnable" bson:"transferFeeEnable"`
	TransferFeeData                *string  `json:"transferFeeData" bson:"transferFeeData"`
	IsToken2022                    bool     `json:"isToken2022" bson:"isToken2022"`
	NonTransferable                *bool    `json:"nonTransferable" bson:"nonTransferable"`
}

type RespTokenCreationInfo struct {
	TxHash         string `json:"txHash" bson:"txHash"`
	Slot           int64  `json:"slot" bson:"slot"`
	TokenAddress   string `json:"tokenAddress" bson:"tokenAddress"`
	Decimals       int64  `json:"decimals" bson:"decimals"`
	Owner          string `json:"owner" bson:"owner"`
	BlockUnixTime  int64  `json:"blockUnixTime" bson:"blockUnixTime"`
	BlockHumanTime string `json:"blockHumanTime" bson:"blockHumanTime"`
}

type RespMarketTokenInfo struct {
	Address  string `json:"address" bson:"address"`
	Decimals int64  `json:"decimals" bson:"decimals"`
	Icon     string `json:"icon" bson:"icon"`
	Symbol   string `json:"symbol" bson:"symbol"`
}

type RespMarketItem struct {
	Address                      string              `json:"address" bson:"address"`
	Base                         RespMarketTokenInfo `json:"base" bson:"base"`
	CreatedAt                    string              `json:"createdAt" bson:"createdAt"`
	Name                         string              `json:"name" bson:"name"`
	Quote                        RespMarketTokenInfo `json:"quote" bson:"quote"`
	Source                       string              `json:"source" bson:"source"`
	Liquidity                    float64             `json:"liquidity" bson:"liquidity"`
	LiquidityChangePercentage24h *float64            `json:"liquidityChangePercentage24h" bson:"liquidityChangePercentage24h"`
	Price                        float64             `json:"price" bson:"price"`
	Trade24h                     int64               `json:"trade24h" bson:"trade24h"`
	Trade24hChangePercent        float64             `json:"trade24hChangePercent" bson:"trade24hChangePercent"`
	UniqueWallet24h              int64               `json:"uniqueWallet24h" bson:"uniqueWallet24h"`
	UniqueWallet24hChangePercent float64             `json:"uniqueWallet24hChangePercent" bson:"uniqueWallet24hChangePercent"`
	Volume24h                    float64             `json:"volume24h" bson:"volume24h"`
	Volume24hChangePercentage24h *float64            `json:"volume24hChangePercentage24h" bson:"volume24hChangePercentage24h"`
}

type RespNewTokenListingItem struct {
	Address          string  `json:"address" bson:"address"`
	Symbol           string  `json:"symbol" bson:"symbol"`
	Name             string  `json:"name" bson:"name"`
	Decimals         int64   `json:"decimals" bson:"decimals"`
	LiquidityAddedAt string  `json:"liquidityAddedAt" bson:"liquidityAddedAt"`
	Liquidity        float64 `json:"liquidity" bson:"liquidity"`
}

type RespTopTraderItem struct {
	Owner        string   `json:"owner" bson:"owner"`
	TokenAddress string   `json:"tokenAddress" bson:"tokenAddress"`
	Trade        int64    `json:"trade" bson:"trade"`
	TradeBuy     int64    `json:"tradeBuy" bson:"tradeBuy"`
	TradeSell    int64    `json:"tradeSell" bson:"tradeSell"`
	Type         string   `json:"type" bson:"type"`
	Volume       float64  `json:"volume" bson:"volume"`
	VolumeBuy    float64  `json:"volumeBuy" bson:"volumeBuy"`
	VolumeSell   float64  `json:"volumeSell" bson:"volumeSell"`
	Tags         []string `json:"tags" bson:"tags"`
}

type RespWalletBalanceChange struct {
	Amount   int64  `json:"amount" bson:"amount"`
	Symbol   string `json:"symbol" bson:"symbol"`
	Name     string `json:"name" bson:"name"`
	Decimals int64  `json:"decimals" bson:"decimals"`
	Address  string `json:"address" bson:"address"`
	LogoURI  string `json:"logoURI" bson:"logoURI"`
}

type RespContractLabel struct {
	Address  string `json:"address" bson:"address"`
	Name     string `json:"name" bson:"name"`
	Metadata struct {
		Icon string `json:"icon" bson:"icon"`
	} `json:"metadata" bson:"metadata"`
}

type RespWalletHistory struct {
	TxHash        string                    `json:"txHash" bson:"txHash"`
	BlockNumber   int64                     `json:"blockNumber" bson:"blockNumber"`
	BlockTime     string                    `json:"blockTime" bson:"blockTime"`
	Status        bool                      `json:"status" bson:"status"`
	From          string                    `json:"from" bson:"from"`
	To            string                    `json:"to" bson:"to"`
	Fee           int64                     `json:"fee" bson:"fee"`
	MainAction    string                    `json:"mainAction" bson:"mainAction"`
	BalanceChange []RespWalletBalanceChange `json:"balanceChange" bson:"balanceChange"`
	ContractLabel RespContractLabel         `json:"contractLabel" bson:"contractLabel"`
}

type RespWalletPortfolioItem struct {
	Address  string  `json:"address" bson:"address"`
	Decimals int64   `json:"decimals" bson:"decimals"`
	Balance  int64   `json:"balance" bson:"balance"`
	UiAmount float64 `json:"uiAmount" bson:"uiAmount"`
	ChainId  string  `json:"chainId" bson:"chainId"`
	Name     string  `json:"name" bson:"name"`
	Symbol   string  `json:"symbol" bson:"symbol"`
	LogoURI  string  `json:"logoURI" bson:"logoURI"`
	PriceUsd float64 `json:"priceUsd" bson:"priceUsd"`
	ValueUsd float64 `json:"valueUsd" bson:"valueUsd"`
}

type RespWalletPortfolio struct {
	Wallet   string                    `json:"wallet" bson:"wallet"`
	TotalUsd float64                   `json:"totalUsd" bson:"totalUsd"`
	Items    []RespWalletPortfolioItem `json:"items" bson:"items"`
}

type Client struct {
	apiKey  string
	limiter *golimiter.ReqLimiter
}

func NewClient(apiKey string, limiter *golimiter.ReqLimiter) *Client {
	return &Client{
		apiKey:  apiKey,
		limiter: limiter,
	}
}

func (c *Client) newHeader(chains ...string) http.Header {
	header := http.Header{}
	header.Set("content-type", "application/json")
	header.Set("x-api-key", c.apiKey)
	if len(chains) > 0 {
		header.Set("x-chain", strings.Join(chains, ","))
	}
	return header
}

func get[D any](ctx context.Context, clt *Client, path string, chains []string, params ...any) (D, error) {
	clt.limiter.Wait(ctx)
	ul := fmt.Sprintf("%s%s", BASE_URL, path)
	ps := url.Values{}
	for i := 0; i < len(params)-1; i += 2 {
		key := fmt.Sprintf("%v", params[i])
		var value string
		switch v := params[i+1].(type) {
		case string:
			value = v
		case int64:
			value = strconv.FormatInt(v, 10)
		case float64:
			value = strconv.FormatFloat(v, 'f', -1, 64)
		case []string:
			value = strings.Join(v, ",")
		default:
			value = fmt.Sprintf("%v", v)
		}
		ps.Add(key, value)
	}
	ul += "?" + ps.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", ul, nil)
	if err != nil {
		return *new(D), fmt.Errorf("birdeye: new request: %w", err)
	}
	req.Header = clt.newHeader(chains...)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return *new(D), fmt.Errorf("birdeye: do request: %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	err = statusCodeToError[statusCode]
	if err != nil {
		return *new(D), fmt.Errorf("birdeye: status code: %d, message: %s", statusCode, err)
	}

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return *new(D), fmt.Errorf("birdeye: read response: %w", err)
	// }
	// fmt.Println(string(body))

	var rd RespData[D]
	if err := json.NewDecoder(resp.Body).Decode(&rd); err != nil {
		return *new(D), fmt.Errorf("birdeye: decode response: %w", err)
	}

	if statusCode == http.StatusOK {
		return rd.Data, nil
	}

	return *new(D), fmt.Errorf("birdeye: status code: %d, message: %s", statusCode, rd.Message)
}

func (c *Client) SupportedNetworks() ([]string, error) {
	return get[[]string](context.Background(), c, "/defi/networks", nil)
}

func (c *Client) Price(chain string, token string, includeLiquidity bool, checkLiquidity float64) (RespPrice, error) {
	params := []any{"address", token}
	if includeLiquidity {
		params = append(params, "include_liquidity", includeLiquidity)
	}
	if checkLiquidity > 0 {
		params = append(params, "check_liquidity", checkLiquidity)
	}
	return get[RespPrice](context.Background(), c, "/defi/price", []string{chain}, params...)
}

func (c *Client) PriceHistory(chain string, address string, addressType AddressType, chartType ChartType, timeFrom, timeTo int64) (RespItems[RespPriceHistoryItem], error) {
	return get[RespItems[RespPriceHistoryItem]](context.Background(), c, "/defi/history_price", []string{chain},
		"address", address, "address_type", addressType, "type", chartType, "time_from", timeFrom, "time_to", timeTo)
}

func (c *Client) MultiPrice(chain string, listAddress []string, includeLiquidity bool, checkLiquidity float64) (RespMultiPrice, error) {
	params := []any{"list_address", listAddress}
	if includeLiquidity {
		params = append(params, "include_liquidity", includeLiquidity)
	}
	if checkLiquidity > 0 {
		params = append(params, "check_liquidity", checkLiquidity)
	}
	return get[RespMultiPrice](context.Background(), c, "/defi/multi_price", []string{chain}, params...)
}

// OHLCVByToken retrieves OHLCV (Open, High, Low, Close, Volume) data for a specific token
//
// Parameters:
//   - chain: The blockchain network
//   - address: The token address to get OHLCV data for
//   - chartType: The interval type for OHLCV data (1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H, 6H, 8H, 12H, 1D, 3D, 1W, 1M)
//   - timeFrom: Unix timestamp for start of time range
//   - timeTo: Unix timestamp for end of time range
//
// Returns:
//   - []RespOHLCVItem: Array of OHLCV data points
//   - error: Any error that occurred during the request
func (c *Client) OHLCVByToken(chain string, address string, chartType ChartType, timeFrom, timeTo int64) (RespItems[RespOHLCVItem], error) {
	return get[RespItems[RespOHLCVItem]](context.Background(), c, "/defi/ohlcv", []string{chain},
		"address", address, "type", chartType, "time_from", timeFrom, "time_to", timeTo)
}

// OHLCVByPair retrieves OHLCV (Open, High, Low, Close, Volume) data for a specific trading pair
//
// Parameters:
//   - chain: The blockchain network
//   - address: The trading pair's token address
//   - chartType: The interval type for OHLCV data (1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H, 6H, 8H, 12H, 1D, 3D, 1W, 1M)
//   - timeFrom: Unix timestamp for start of time range
//   - timeTo: Unix timestamp for end of time range
//
// Returns:
//   - []RespOHLCVBaseQuoteItem: Array of OHLCV data points for the trading pair
//   - error: Any error that occurred during the request
func (c *Client) OHLCVByPair(chain string, address string, chartType ChartType, timeFrom, timeTo int64) (RespItems[RespOHLCVItem], error) {
	return get[RespItems[RespOHLCVItem]](context.Background(), c, "/defi/ohlcv/pair", []string{chain},
		"address", address, "type", chartType, "time_from", timeFrom, "time_to", timeTo)
}

// OHLCVByBaseQuote retrieves OHLCV (Open, High, Low, Close, Volume) data for a trading pair specified by base and quote token addresses
//
// Parameters:
//   - chain: The blockchain network
//   - baseAddress: The base token's address
//   - quoteAddress: The quote token's address
//   - chartType: The interval type for OHLCV data (1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H, 6H, 8H, 12H, 1D, 3D, 1W, 1M)
//   - timeFrom: Unix timestamp for start of time range
//   - timeTo: Unix timestamp for end of time range
//
// Returns:
//   - []RespOHLCVBaseQuoteItem: Array of OHLCV data points for the trading pair
//   - error: Any error that occurred during the request
func (c *Client) OHLCVByBaseQuote(chain string, baseAddress string, quoteAddress string, chartType ChartType, timeFrom, timeTo int64) (RespItems[RespOHLCVBaseQuoteItem], error) {
	return get[RespItems[RespOHLCVBaseQuoteItem]](context.Background(), c, "/defi/ohlcv/base_quote", []string{chain},
		"base_address", baseAddress, "quote_address", quoteAddress, "type", chartType, "time_from", timeFrom, "time_to", timeTo)
}

// TradesByToken retrieves transaction records for a specific token
//
// Parameters:
//   - chain: The blockchain network
//   - address: The token address to retrieve transaction records for (default: So11111111111111111111111111111111111111112)
//   - sortType: Sort order for trades by UNIX time ("desc" or "asc", default: "desc")
//   - offset: Starting index for the list (0-1000, default: 0)
//   - limit: Maximum number of records to retrieve (1-50, default: 50)
//   - txType: Type of transactions to filter by ("swap", "add", "remove", "all", default: "swap")
//
// Returns:
//   - RespItems[RespTradesByTokenItem]: Paginated list of trade records
//   - error: Any error that occurred during the request
func (c *Client) TradesByToken(chain string, address string, sortType SortType, offset int, limit int, txType TxType) (RespItems[RespTradesByTokenItem], error) {
	return get[RespItems[RespTradesByTokenItem]](context.Background(), c, "/defi/txs/token", []string{chain},
		"address", address,
		"sort_type", sortType,
		"offset", offset,
		"limit", limit,
		"tx_type", txType)
}

// TradesByPair retrieves transaction records for a specific trading pair
//
// Parameters:
//   - chain: The blockchain network
//   - address: The pair address to retrieve transaction records for
//   - sortType: Sort order for trades by UNIX time ("desc" or "asc", default: "desc")
//   - offset: Starting index for the list (0-1000, default: 0)
//   - limit: Maximum number of records to retrieve (1-50, default: 50)
//   - txType: Type of transactions to filter by ("swap", "add", "remove", "all", default: "swap")
//
// Returns:
//   - RespItems[RespTradesByPairItem]: Paginated list of trade records
//   - error: Any error that occurred during the request
func (c *Client) TradesByPair(chain string, address string, sortType SortType, offset int, limit int, txType TxType) (RespItems[RespTradesByPairItem], error) {
	return get[RespItems[RespTradesByPairItem]](context.Background(), c, "/defi/txs/pair", []string{chain},
		"address", address,
		"sort_type", sortType,
		"offset", offset,
		"limit", limit,
		"tx_type", txType)
}

// HistoricalPriceByUnix retrieves the historical price of a token at a specific Unix timestamp
//
// Parameters:
//   - chain: The blockchain network
//   - address: The Solana token address to retrieve historical prices for
//   - unixTime: The Unix timestamp representing the specific point in time
//
// Returns:
//   - RespPriceHistoryByTime: Historical price data at the specified timestamp
//   - error: Any error that occurred during the request
func (c *Client) HistoricalPriceByUnix(chain string, address string, unixTime int64) (RespPriceHistoryByTime, error) {
	return get[RespPriceHistoryByTime](context.Background(), c, "/defi/historical_price_unix", []string{chain},
		"address", address,
		"unixtime", unixTime)
}

// PriceVolumeByToken retrieves price and volume data for a specific token over a time period
//
// Parameters:
//   - chain: The blockchain network
//   - address: The Solana token address to retrieve price and volume data for
//   - timeType: The time period for the data ("1h", "2h", "4h", "8h", "24h", default: "24h")
//
// Returns:
//   - RespSinglePriceVolume: Price and volume data for the specified token and time period
//   - error: Any error that occurred during the request
func (c *Client) PriceVolumeByToken(chain string, address string, timeType TimeType) (RespSinglePriceVolume, error) {
	return get[RespSinglePriceVolume](context.Background(), c, "/defi/price_volume/single", []string{chain},
		"address", address,
		"type", timeType)
}

// PriceVolumeByTokens retrieves price and volume data for multiple tokens over a time period
//
// Parameters:
//   - chain: The blockchain network
//   - listAddress: Comma-separated list of token addresses to retrieve price and volume data for
//   - timeType: The time period for the data ("1h", "2h", "4h", "8h", "24h", default: "24h")
//
// Returns:
//   - []RespSinglePriceVolume: Price and volume data for the specified tokens and time period
//   - error: Any error that occurred during the request
func (c *Client) PriceVolumeByTokens(chain string, listAddress []string, timeType TimeType) ([]RespSinglePriceVolume, error) {
	return get[[]RespSinglePriceVolume](context.Background(), c, "/defi/price_volume/multi", []string{chain},
		"list_address", listAddress,
		"type", timeType)
}

// TrendingTokens retrieves a list of trending tokens with sorting and pagination options
//
// Parameters:
//   - chain: The blockchain network
//   - sortBy: The attribute to sort tokens by ("rank", "liquidity", "volume24hUSD", default: "rank")
//   - sortType: The sort order ("asc", "desc", default: "asc")
//   - offset: Number of records to skip (default: 0)
//   - limit: Maximum number of records to return (max: 20, default: 20)
//
// Returns:
//   - RespTrendingTokens: List of trending tokens with metadata
//   - error: Any error that occurred during the request
func (c *Client) TrendingTokens(chain string, sortBy RankType, sortType SortType, offset int, limit int) (RespTrendingTokens, error) {
	if limit > 20 {
		limit = 20
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return get[RespTrendingTokens](context.Background(), c, "/defi/token_trending", []string{chain},
		"sort_by", sortBy,
		"sort_type", sortType,
		"offset", offset,
		"limit", limit)
}

// TradeByTokenAndTime retrieves token transaction data based on Unix time
//
// Parameters:
//   - chain: The blockchain network
//   - address: The token address to retrieve transaction records for
//   - beforeTime: Filters transactions that occurred before this UNIX timestamp (in seconds)
//   - afterTime: Filters transactions that occurred after this UNIX timestamp (in seconds)
//   - txType: The type of transactions to filter by ("swap", "add", "remove", "all", default: "swap")
//   - offset: Number of records to skip (default: 0)
//   - limit: Maximum number of records to return (max: 50, default: 50)
//
// Returns:
//   - RespItems[RespTradesByTokenItem]: List of token transactions with metadata
//   - error: Any error that occurred during the request
//
// Note: beforeTime and afterTime cannot be used simultaneously
func (c *Client) TradeByTokenAndTime(chain string, address string, beforeTime, afterTime int64, txType TxType, offset int, limit int) (RespItems[RespTradesByTokenItem], error) {
	if beforeTime > 0 && afterTime > 0 {
		return RespItems[RespTradesByTokenItem]{}, fmt.Errorf("beforeTime and afterTime cannot be used simultaneously")
	}

	if limit > 50 {
		limit = 50
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	params := []any{
		"address", address,
		"offset", offset,
		"limit", limit,
	}

	if txType != "" {
		params = append(params, "tx_type", txType)
	}

	if beforeTime > 0 {
		params = append(params, "before_time", beforeTime)
	}
	if afterTime > 0 {
		params = append(params, "after_time", afterTime)
	}

	d, err := get[RespItems[RespTradesByTokenItem]](context.Background(), c, "/defi/txs/token/seek_by_time", []string{chain}, params...)
	if err != nil {
		return RespItems[RespTradesByTokenItem]{}, err
	}
	return d, nil
}

// TradesByPairAndTime retrieves transaction records for a specific pair address based on Unix time.
// It allows querying up to 10,000 records from a specified point in time.
//
// Parameters:
//   - chain: The blockchain network to query
//   - address: The pair address to retrieve transaction records for
//   - beforeTime: Filters transactions that occurred before this UNIX timestamp (in seconds)
//   - afterTime: Filters transactions that occurred after this UNIX timestamp (in seconds)
//   - txType: The type of transactions to filter by ("swap", "add", "remove", "all", default: "swap")
//   - offset: Starting index for pagination (default: 0, max: 1000)
//   - limit: Maximum number of records per request (default: 50, max: 50)
//
// Returns:
//   - RespItems[RespTradesByPairItem]: List of pair transactions with metadata
//   - error: Any error that occurred during the request
//
// Usage Notes:
//   - Use beforeTime to retrieve historical records before a specified Unix time
//   - Use afterTime to retrieve recent records after a specified Unix time
//   - beforeTime and afterTime cannot be used simultaneously (will return error 422)
//   - For sequential queries, adjust beforeTime/afterTime based on the last record's timestamp
//   - Maximum of 10,000 records can be retrieved from a specified time point
func (c *Client) TradesByPairAndTime(chain string, address string, beforeTime, afterTime int64, txType TxType, offset int, limit int) (RespItems[RespTradesByPairItem], error) {
	if beforeTime > 0 && afterTime > 0 {
		return RespItems[RespTradesByPairItem]{}, fmt.Errorf("beforeTime and afterTime cannot be used simultaneously (error 422)")
	}

	if beforeTime == 0 && afterTime == 0 {
		return RespItems[RespTradesByPairItem]{}, fmt.Errorf("either beforeTime or afterTime must be specified")
	}

	if limit > 50 {
		limit = 50
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if offset > 1000 {
		offset = 1000
	}

	params := []any{
		"address", address,
		"offset", offset,
		"limit", limit,
	}

	if txType != "" {
		params = append(params, "tx_type", txType)
	}

	if beforeTime > 0 {
		params = append(params, "before_time", beforeTime)
	}
	if afterTime > 0 {
		params = append(params, "after_time", afterTime)
	}

	d, err := get[RespItems[RespTradesByPairItem]](context.Background(), c, "/defi/txs/pair/seek_by_time", []string{chain}, params...)
	if err != nil {
		return RespItems[RespTradesByPairItem]{}, err
	}
	return d, nil
}

// TokenOverview returns detailed information about a token, including price changes, volume, and social metrics
func (c *Client) TokenOverview(chain string, address string) (RespTokenOverview, error) {
	return get[RespTokenOverview](context.Background(), c, "/defi/token_overview", []string{chain}, "address", address)
}

// TokenList retrieves a list of tokens sorted by specified criteria
//
// Parameters:
//   - chain: The blockchain network
//   - sortType: Attribute to sort tokens by ("v24hUSD", "mc", "v24hChangePercent", default: "v24hUSD")
//   - sortOrder: Sort order ("asc" or "desc", default: "desc")
//   - offset: Starting index for the list (0-1000, default: 0)
//   - limit: Maximum number of tokens to retrieve (1-50, default: 50)
//   - minLiquidity: Minimum liquidity filter (default: 100)
//
// Returns:
//   - RespItems[RespToken]: Paginated list of token information
//   - error: Any error that occurred during the request
func (c *Client) TokenList(chain string, sortBy TokenListSortType, sortType SortType, offset int, limit int, minLiquidity float64) (RespItems[RespToken], error) {
	if limit > 50 {
		limit = 50
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if offset > 1000 {
		offset = 1000
	}

	params := []any{
		"sort_by", sortBy,
		"sort_type", sortType,
		"offset", offset,
		"limit", limit,
	}

	if minLiquidity > 0 {
		params = append(params, "min_liquidity", minLiquidity)
	}

	return get[RespItems[RespToken]](context.Background(), c, "/defi/tokenlist", []string{chain}, params...)
}

// TokenListV2 retrieves a URL to download the complete token list
//
// Parameters:
//   - chain: The blockchain network
//
// Returns:
//   - RespTokenListV2Url: URL to download the complete token list
//   - error: Any error that occurred during the request
//
// Note: The returned URL can be used to download a JSON file containing
// the complete list of tokens and their metadata for the specified chain.
func (c *Client) TokenListV2(chain string) (RespTokenListV2Url, error) {
	return get[RespTokenListV2Url](context.Background(), c, "/defi/v2/tokens/all", []string{chain})
}

// TokenSecurity retrieves security information for a specific token
//
// Parameters:
//   - chain: The blockchain network
//   - address: Token address to get security info for
//
// Returns:
//   - RespTokenSecurity: Security information for the token
//   - error: Any error that occurred during the request
func (c *Client) TokenSecurity(chain string, address string) (RespTokenSecurity, error) {
	return get[RespTokenSecurity](context.Background(), c, "/defi/token_security", []string{chain},
		"address", address)
}

// TokenCreationInfo retrieves creation information for a specific token
//
// Parameters:
//   - chain: The blockchain network
//   - address: Token address to get creation info for
//
// Returns:
//   - RespTokenSecurity: Creation information for the token
//   - error: Any error that occurred during the request
func (c *Client) TokenCreationInfo(chain string, address string) (RespTokenCreationInfo, error) {
	return get[RespTokenCreationInfo](context.Background(), c, "/defi/token_creation_info", []string{chain},
		"address", address)
}

// MarketList retrieves a list of markets for a specific token
//
// Parameters:
//   - chain: The blockchain network
//   - address: Token address to get market list for (default: So11111111111111111111111111111111111111112)
//   - sortBy: Sort markets by "liquidity" or "volume24h" (default: "liquidity")
//   - sortType: Sort order "desc" or "asc" (default: "desc")
//   - offset: Starting index for pagination (default: 0)
//   - limit: Maximum number of markets to return (1-10, default: 10)
//
// Returns:
//   - RespItems[RespToken]: Paginated list of market data
//   - error: Any error that occurred during the request
func (c *Client) MarketList(chain string, address string, sortBy MarketListSortType, sortType SortType, offset int, limit int) (RespItems[RespMarketItem], error) {
	if limit > 10 {
		limit = 10
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return get[RespItems[RespMarketItem]](context.Background(), c, "/defi/v2/markets", []string{chain},
		"address", address,
		"sort_by", sortBy,
		"sort_type", sortType,
		"offset", offset,
		"limit", limit)
}

// NewTokenListing retrieves newly listed tokens up to a specified time
//
// Parameters:
//   - chain: The blockchain network
//   - timeTo: Unix timestamp to retrieve tokens listed up to
//   - limit: Maximum number of records to return (1-10, default: 10)
//   - memePlatformEnabled: Include tokens from meme platforms (Solana only, default: false)
//
// Returns:
//   - RespItems[RespToken]: Paginated list of newly listed tokens
//   - error: Any error that occurred during the request
func (c *Client) NewTokenListing(chain string, timeTo int64, limit int, memePlatformEnabled bool) (RespItems[RespNewTokenListingItem], error) {
	if limit > 10 {
		limit = 10
	}
	if limit <= 0 {
		limit = 10
	}

	params := []any{
		"time_to", timeTo,
		"limit", limit,
	}

	if memePlatformEnabled {
		params = append(params, "meme_platform_enabled", true)
	}

	return get[RespItems[RespNewTokenListingItem]](context.Background(), c, "/defi/v2/tokens/new_listing", []string{chain}, params...)
}

// TokenTopTraders retrieves the top traders for a specific token based on volume or trade count
//
// Parameters:
//   - chain: The blockchain network
//   - address: The token address to retrieve top traders for (required)
//   - sortBy: Attribute to sort traders by ("volume" or "trade", default: "volume")
//   - sortType: Sort order ("asc" or "desc", default: "desc")
//   - timeFrame: Time period for the data ("30m","1h","2h","4h","6h","8h","12h","24h", default: "24h")
//   - offset: Number of records to skip (default: 0)
//   - limit: Maximum number of records to return (1-10, default: 10)
//
// Returns:
//   - RespItems[RespToken]: Paginated list of top traders
//   - error: Any error that occurred during the request
func (c *Client) TokenTopTraders(chain string, address string, sortBy TokenTopTradersSortType, sortType SortType, timeFrame TopTradersTimeFrame, offset int64, limit int64) (RespItems[RespTopTraderItem], error) {
	if limit > 10 {
		limit = 10
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return get[RespItems[RespTopTraderItem]](context.Background(), c, "/defi/v2/tokens/top_traders", []string{chain},
		"address", address,
		"sort_by", sortBy,
		"sort_type", sortType,
		"time_frame", timeFrame,
		"offset", offset,
		"limit", limit)
}

// WalletTxHistories retrieves transaction history for a specific wallet address
//
// Parameters:
//   - chain: The blockchain network
//   - wallet: The wallet address to retrieve transaction history for
//   - limit: Maximum number of transactions to return (default: 50)
//   - before: Transaction hash cursor for pagination (transactions before this hash will be returned)
//
// Returns:
//   - RespItems[RespTradesByTokenItem]: List of transactions for the wallet
//   - error: Any error that occurred during the request
func (c *Client) WalletTxHistories(chain string, wallet string, limit int, before string) (map[ChainType][]RespWalletHistory, error) {
	if limit <= 0 {
		limit = 50
	}
	params := []any{
		"wallet", wallet,
		"limit", limit,
	}
	if before != "" {
		params = append(params, "before", before)
	}
	return get[map[ChainType][]RespWalletHistory](context.Background(), c, "/v1/wallet/tx_list", []string{chain}, params...)
}

// WalletPortfolio retrieves the token portfolio for a specific wallet address
//
// Parameters:
//   - chain: The blockchain network
//   - wallet: The wallet address to retrieve portfolio for
//
// Returns:
//   - RespItems[RespToken]: List of tokens held in the wallet
//   - error: Any error that occurred during the request
func (c *Client) WalletPortfolio(chain string, wallet string) (RespWalletPortfolio, error) {
	return get[RespWalletPortfolio](context.Background(), c, "/v1/wallet/token_list", []string{chain}, "wallet", wallet)
}
