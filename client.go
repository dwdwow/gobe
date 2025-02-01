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
	Value float64 `json:"value"`
	// UpdateUnixTime seconds
	UpdateUnixTime  int64   `json:"updateUnixTime"`
	UpdateHumanTime string  `json:"updateHumanTime"`
	Liquidity       float64 `json:"liquidity"`
}

type RespPriceHistoryItem struct {
	Address  string  `json:"address"`
	UnixTime int64   `json:"unixTime"`
	Value    float64 `json:"value"`
}

type RespMultiPriceInfo struct {
	Value           float64 `json:"value"`
	UpdateUnixTime  int64   `json:"updateUnixTime"`
	UpdateHumanTime string  `json:"updateHumanTime"`
	PriceChange24h  float64 `json:"priceChange24h"`
}

// RespMultiPrice key: address
type RespMultiPrice map[string]RespMultiPriceInfo

type RespOHLCVItem struct {
	Address  string    `json:"address"`
	C        float64   `json:"c"`
	H        float64   `json:"h"`
	L        float64   `json:"l"`
	O        float64   `json:"o"`
	Type     ChartType `json:"type"`
	UnixTime int64     `json:"unixTime"`
	V        float64   `json:"v"`
}

type RespOHLCVBaseQuoteItem struct {
	O            float64   `json:"o"`
	C            float64   `json:"c"`
	H            float64   `json:"h"`
	L            float64   `json:"l"`
	BaseAddress  string    `json:"baseAddress"`
	QuoteAddress string    `json:"quoteAddress"`
	VBase        float64   `json:"vBase"`
	VQuote       float64   `json:"vQuote"`
	Type         ChartType `json:"type"`
	UnixTime     int64     `json:"unixTime"`
}

type RespTradesByTokenTokenInfo struct {
	Symbol         string   `json:"symbol"`
	Decimals       int64    `json:"decimals"`
	Address        string   `json:"address"`
	Amount         int64    `json:"amount"`
	UiAmount       float64  `json:"uiAmount"`
	Price          *float64 `json:"price"`
	NearestPrice   float64  `json:"nearestPrice"`
	ChangeAmount   int64    `json:"changeAmount"`
	UiChangeAmount float64  `json:"uiChangeAmount"`
}

type RespTradesByTokenItem struct {
	Quote         RespTradesByTokenTokenInfo `json:"quote"`
	Base          RespTradesByTokenTokenInfo `json:"base"`
	BasePrice     *float64                   `json:"basePrice"`
	QuotePrice    *float64                   `json:"quotePrice"`
	TxHash        string                     `json:"txHash"`
	Source        string                     `json:"source"`
	BlockUnixTime int64                      `json:"blockUnixTime"`
	TxType        string                     `json:"txType"`
	Owner         string                     `json:"owner"`
	Side          string                     `json:"side"`
	Alias         *string                    `json:"alias"`
	PricePair     float64                    `json:"pricePair"`
	From          RespTradesByTokenTokenInfo `json:"from"`
	To            RespTradesByTokenTokenInfo `json:"to"`
	TokenPrice    *float64                   `json:"tokenPrice"`
	PoolId        string                     `json:"poolId"`
}

type RespTradesByPairTokenInfo struct {
	Symbol         string   `json:"symbol"`
	Decimals       int64    `json:"decimals"`
	Address        string   `json:"address"`
	Amount         int64    `json:"amount"`
	Type           string   `json:"type"`
	TypeSwap       string   `json:"typeSwap"`
	UiAmount       float64  `json:"uiAmount"`
	Price          *float64 `json:"price"`
	NearestPrice   float64  `json:"nearestPrice"`
	ChangeAmount   int64    `json:"changeAmount"`
	UiChangeAmount float64  `json:"uiChangeAmount"`
}

type RespTradesByPairItem struct {
	TxHash        string                    `json:"txHash"`
	Source        string                    `json:"source"`
	BlockUnixTime int64                     `json:"blockUnixTime"`
	Address       string                    `json:"address"`
	Owner         string                    `json:"owner"`
	From          RespTradesByPairTokenInfo `json:"from"`
	To            RespTradesByPairTokenInfo `json:"to"`
}

type RespPriceHistoryByTime struct {
	Value          float64 `json:"value"`
	UpdateUnixTime int64   `json:"updateUnixTime"`
	PriceChange24h float64 `json:"priceChange24h"`
}

type RespSinglePriceVolume struct {
	// Address just use for multi-token
	Address             string  `json:"address"`
	Price               float64 `json:"price"`
	UpdateUnixTime      int64   `json:"updateUnixTime"`
	UpdateHumanTime     string  `json:"updateHumanTime"`
	VolumeUSD           float64 `json:"volumeUSD"`
	VolumeChangePercent float64 `json:"volumeChangePercent"`
	PriceChangePercent  float64 `json:"priceChangePercent"`
}

type RespTrendingTokensTokenInfo struct {
	Address      string  `json:"address"`
	Decimals     int64   `json:"decimals"`
	Liquidity    float64 `json:"liquidity"`
	LogoURI      string  `json:"logoURI"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	Volume24hUSD float64 `json:"volume24hUSD"`
	Rank         int64   `json:"rank"`
}

type RespTrendingTokens struct {
	UpdateUnixTime int64                         `json:"updateUnixTime"`
	UpdateTime     string                        `json:"updateTime"`
	Tokens         []RespTrendingTokensTokenInfo `json:"tokens"`
	Total          int64                         `json:"total"`
}

type RespTokenExtension struct {
	CoingeckoId string `json:"coingeckoId"`
	SerumV3Usdc string `json:"serumV3Usdc"`
	SerumV3Usdt string `json:"serumV3Usdt"`
	Website     string `json:"website"`
	Telegram    string `json:"telegram"`
	Twitter     string `json:"twitter"`
	Description string `json:"description"`
	Discord     string `json:"discord"`
	Medium      string `json:"medium"`
}

type RespTokenOverview struct {
	Address                      string             `json:"address"`
	Decimals                     int64              `json:"decimals"`
	Symbol                       string             `json:"symbol"`
	Name                         string             `json:"name"`
	Extensions                   RespTokenExtension `json:"extensions"`
	LogoURI                      string             `json:"logoURI"`
	Liquidity                    float64            `json:"liquidity"`
	Price                        float64            `json:"price"`
	History30mPrice              float64            `json:"history30mPrice"`
	PriceChange30mPercent        float64            `json:"priceChange30mPercent"`
	History1hPrice               float64            `json:"history1hPrice"`
	PriceChange1hPercent         float64            `json:"priceChange1hPercent"`
	History2hPrice               float64            `json:"history2hPrice"`
	PriceChange2hPercent         float64            `json:"priceChange2hPercent"`
	History4hPrice               float64            `json:"history4hPrice"`
	PriceChange4hPercent         float64            `json:"priceChange4hPercent"`
	History6hPrice               float64            `json:"history6hPrice"`
	PriceChange6hPercent         float64            `json:"priceChange6hPercent"`
	History8hPrice               float64            `json:"history8hPrice"`
	PriceChange8hPercent         float64            `json:"priceChange8hPercent"`
	History12hPrice              float64            `json:"history12hPrice"`
	PriceChange12hPercent        float64            `json:"priceChange12hPercent"`
	History24hPrice              float64            `json:"history24hPrice"`
	PriceChange24hPercent        float64            `json:"priceChange24hPercent"`
	UniqueWallet30m              int64              `json:"uniqueWallet30m"`
	UniqueWalletHistory30m       int64              `json:"uniqueWalletHistory30m"`
	UniqueWallet30mChangePercent float64            `json:"uniqueWallet30mChangePercent"`
	UniqueWallet1h               int64              `json:"uniqueWallet1h"`
	UniqueWalletHistory1h        int64              `json:"uniqueWalletHistory1h"`
	UniqueWallet1hChangePercent  float64            `json:"uniqueWallet1hChangePercent"`
	UniqueWallet2h               int64              `json:"uniqueWallet2h"`
	UniqueWalletHistory2h        int64              `json:"uniqueWalletHistory2h"`
	UniqueWallet2hChangePercent  float64            `json:"uniqueWallet2hChangePercent"`
	UniqueWallet4h               int64              `json:"uniqueWallet4h"`
	UniqueWalletHistory4h        int64              `json:"uniqueWalletHistory4h"`
	UniqueWallet4hChangePercent  float64            `json:"uniqueWallet4hChangePercent"`
	UniqueWallet6h               int64              `json:"uniqueWallet6h"`
	UniqueWalletHistory6h        int64              `json:"uniqueWalletHistory6h"`
	UniqueWallet6hChangePercent  float64            `json:"uniqueWallet6hChangePercent"`
	UniqueWallet8h               int64              `json:"uniqueWallet8h"`
	UniqueWalletHistory8h        int64              `json:"uniqueWalletHistory8h"`
	UniqueWallet8hChangePercent  float64            `json:"uniqueWallet8hChangePercent"`
	UniqueWallet12h              int64              `json:"uniqueWallet12h"`
	UniqueWalletHistory12h       int64              `json:"uniqueWalletHistory12h"`
	UniqueWallet12hChangePercent float64            `json:"uniqueWallet12hChangePercent"`
	UniqueWallet24h              int64              `json:"uniqueWallet24h"`
	UniqueWalletHistory24h       int64              `json:"uniqueWalletHistory24h"`
	UniqueWallet24hChangePercent float64            `json:"uniqueWallet24hChangePercent"`
	LastTradeUnixTime            int64              `json:"lastTradeUnixTime"`
	LastTradeHumanTime           string             `json:"lastTradeHumanTime"`
	Supply                       float64            `json:"supply"`
	Mc                           float64            `json:"mc"`
	Trade30m                     int64              `json:"trade30m"`
	TradeHistory30m              int64              `json:"tradeHistory30m"`
	Trade30mChangePercent        float64            `json:"trade30mChangePercent"`
	Sell30m                      int64              `json:"sell30m"`
	SellHistory30m               int64              `json:"sellHistory30m"`
	Sell30mChangePercent         float64            `json:"sell30mChangePercent"`
	Buy30m                       int64              `json:"buy30m"`
	BuyHistory30m                int64              `json:"buyHistory30m"`
	Buy30mChangePercent          float64            `json:"buy30mChangePercent"`
	V30m                         float64            `json:"v30m"`
	V30mUSD                      float64            `json:"v30mUSD"`
	VHistory30m                  float64            `json:"vHistory30m"`
	VHistory30mUSD               float64            `json:"vHistory30mUSD"`
	V30mChangePercent            float64            `json:"v30mChangePercent"`
	VBuy30m                      float64            `json:"vBuy30m"`
	VBuy30mUSD                   float64            `json:"vBuy30mUSD"`
	VBuyHistory30m               float64            `json:"vBuyHistory30m"`
	VBuyHistory30mUSD            float64            `json:"vBuyHistory30mUSD"`
	VBuy30mChangePercent         float64            `json:"vBuy30mChangePercent"`
	VSell30m                     float64            `json:"vSell30m"`
	VSell30mUSD                  float64            `json:"vSell30mUSD"`
	VSellHistory30m              float64            `json:"vSellHistory30m"`
	VSellHistory30mUSD           float64            `json:"vSellHistory30mUSD"`
	VSell30mChangePercent        float64            `json:"vSell30mChangePercent"`
	Trade1h                      int64              `json:"trade1h"`
	TradeHistory1h               int64              `json:"tradeHistory1h"`
	Trade1hChangePercent         float64            `json:"trade1hChangePercent"`
	Sell1h                       int64              `json:"sell1h"`
	SellHistory1h                int64              `json:"sellHistory1h"`
	Sell1hChangePercent          float64            `json:"sell1hChangePercent"`
	Buy1h                        int64              `json:"buy1h"`
	BuyHistory1h                 int64              `json:"buyHistory1h"`
	Buy1hChangePercent           float64            `json:"buy1hChangePercent"`
	V1h                          float64            `json:"v1h"`
	V1hUSD                       float64            `json:"v1hUSD"`
	VHistory1h                   float64            `json:"vHistory1h"`
	VHistory1hUSD                float64            `json:"vHistory1hUSD"`
	V1hChangePercent             float64            `json:"v1hChangePercent"`
	VBuy1h                       float64            `json:"vBuy1h"`
	VBuy1hUSD                    float64            `json:"vBuy1hUSD"`
	VBuyHistory1h                float64            `json:"vBuyHistory1h"`
	VBuyHistory1hUSD             float64            `json:"vBuyHistory1hUSD"`
	VBuy1hChangePercent          float64            `json:"vBuy1hChangePercent"`
	VSell1h                      float64            `json:"vSell1h"`
	VSell1hUSD                   float64            `json:"vSell1hUSD"`
	VSellHistory1h               float64            `json:"vSellHistory1h"`
	VSellHistory1hUSD            float64            `json:"vSellHistory1hUSD"`
	VSell1hChangePercent         float64            `json:"vSell1hChangePercent"`
	Trade2h                      int64              `json:"trade2h"`
	TradeHistory2h               int64              `json:"tradeHistory2h"`
	Trade2hChangePercent         float64            `json:"trade2hChangePercent"`
	Sell2h                       int64              `json:"sell2h"`
	SellHistory2h                int64              `json:"sellHistory2h"`
	Sell2hChangePercent          float64            `json:"sell2hChangePercent"`
	Buy2h                        int64              `json:"buy2h"`
	BuyHistory2h                 int64              `json:"buyHistory2h"`
	Buy2hChangePercent           float64            `json:"buy2hChangePercent"`
	V2h                          float64            `json:"v2h"`
	V2hUSD                       float64            `json:"v2hUSD"`
	VHistory2h                   float64            `json:"vHistory2h"`
	VHistory2hUSD                float64            `json:"vHistory2hUSD"`
	V2hChangePercent             float64            `json:"v2hChangePercent"`
	VBuy2h                       float64            `json:"vBuy2h"`
	VBuy2hUSD                    float64            `json:"vBuy2hUSD"`
	VBuyHistory2h                float64            `json:"vBuyHistory2h"`
	VBuyHistory2hUSD             float64            `json:"vBuyHistory2hUSD"`
	VBuy2hChangePercent          float64            `json:"vBuy2hChangePercent"`
	VSell2h                      float64            `json:"vSell2h"`
	VSell2hUSD                   float64            `json:"vSell2hUSD"`
	VSellHistory2h               float64            `json:"vSellHistory2h"`
	VSellHistory2hUSD            float64            `json:"vSellHistory2hUSD"`
	VSell2hChangePercent         float64            `json:"vSell2hChangePercent"`
	Trade4h                      int64              `json:"trade4h"`
	TradeHistory4h               int64              `json:"tradeHistory4h"`
	Trade4hChangePercent         float64            `json:"trade4hChangePercent"`
	Sell4h                       int64              `json:"sell4h"`
	SellHistory4h                int64              `json:"sellHistory4h"`
	Sell4hChangePercent          float64            `json:"sell4hChangePercent"`
	Buy4h                        int64              `json:"buy4h"`
	BuyHistory4h                 int64              `json:"buyHistory4h"`
	Buy4hChangePercent           float64            `json:"buy4hChangePercent"`
	V4h                          float64            `json:"v4h"`
	V4hUSD                       float64            `json:"v4hUSD"`
	VHistory4h                   float64            `json:"vHistory4h"`
	VHistory4hUSD                float64            `json:"vHistory4hUSD"`
	V4hChangePercent             float64            `json:"v4hChangePercent"`
	VBuy4h                       float64            `json:"vBuy4h"`
	VBuy4hUSD                    float64            `json:"vBuy4hUSD"`
	VBuyHistory4h                float64            `json:"vBuyHistory4h"`
	VBuyHistory4hUSD             float64            `json:"vBuyHistory4hUSD"`
	VBuy4hChangePercent          float64            `json:"vBuy4hChangePercent"`
	VSell4h                      float64            `json:"vSell4h"`
	VSell4hUSD                   float64            `json:"vSell4hUSD"`
	VSellHistory4h               float64            `json:"vSellHistory4h"`
	VSellHistory4hUSD            float64            `json:"vSellHistory4hUSD"`
	VSell4hChangePercent         float64            `json:"vSell4hChangePercent"`
	Trade6h                      int64              `json:"trade6h"`
	TradeHistory6h               int64              `json:"tradeHistory6h"`
	Trade6hChangePercent         float64            `json:"trade6hChangePercent"`
	Sell6h                       int64              `json:"sell6h"`
	SellHistory6h                int64              `json:"sellHistory6h"`
	Sell6hChangePercent          float64            `json:"sell6hChangePercent"`
	Buy6h                        int64              `json:"buy6h"`
	BuyHistory6h                 int64              `json:"buyHistory6h"`
	Buy6hChangePercent           float64            `json:"buy6hChangePercent"`
	V6h                          float64            `json:"v6h"`
	V6hUSD                       float64            `json:"v6hUSD"`
	VHistory6h                   float64            `json:"vHistory6h"`
	VHistory6hUSD                float64            `json:"vHistory6hUSD"`
	V6hChangePercent             float64            `json:"v6hChangePercent"`
	VBuy6h                       float64            `json:"vBuy6h"`
	VBuy6hUSD                    float64            `json:"vBuy6hUSD"`
	VBuyHistory6h                float64            `json:"vBuyHistory6h"`
	VBuyHistory6hUSD             float64            `json:"vBuyHistory6hUSD"`
	VBuy6hChangePercent          float64            `json:"vBuy6hChangePercent"`
	VSell6h                      float64            `json:"vSell6h"`
	VSell6hUSD                   float64            `json:"vSell6hUSD"`
	VSellHistory6h               float64            `json:"vSellHistory6h"`
	VSellHistory6hUSD            float64            `json:"vSellHistory6hUSD"`
	VSell6hChangePercent         float64            `json:"vSell6hChangePercent"`
	Trade8h                      int64              `json:"trade8h"`
	TradeHistory8h               int64              `json:"tradeHistory8h"`
	Trade8hChangePercent         float64            `json:"trade8hChangePercent"`
	Sell8h                       int64              `json:"sell8h"`
	SellHistory8h                int64              `json:"sellHistory8h"`
	Sell8hChangePercent          float64            `json:"sell8hChangePercent"`
	Buy8h                        int64              `json:"buy8h"`
	BuyHistory8h                 int64              `json:"buyHistory8h"`
	Buy8hChangePercent           float64            `json:"buy8hChangePercent"`
	V8h                          float64            `json:"v8h"`
	V8hUSD                       float64            `json:"v8hUSD"`
	VHistory8h                   float64            `json:"vHistory8h"`
	VHistory8hUSD                float64            `json:"vHistory8hUSD"`
	V8hChangePercent             float64            `json:"v8hChangePercent"`
	VBuy8h                       float64            `json:"vBuy8h"`
	VBuy8hUSD                    float64            `json:"vBuy8hUSD"`
	VBuyHistory8h                float64            `json:"vBuyHistory8h"`
	VBuyHistory8hUSD             float64            `json:"vBuyHistory8hUSD"`
	VBuy8hChangePercent          float64            `json:"vBuy8hChangePercent"`
	VSell8h                      float64            `json:"vSell8h"`
	VSell8hUSD                   float64            `json:"vSell8hUSD"`
	VSellHistory8h               float64            `json:"vSellHistory8h"`
	VSellHistory8hUSD            float64            `json:"vSellHistory8hUSD"`
	VSell8hChangePercent         float64            `json:"vSell8hChangePercent"`
	Trade12h                     int64              `json:"trade12h"`
	TradeHistory12h              int64              `json:"tradeHistory12h"`
	Trade12hChangePercent        float64            `json:"trade12hChangePercent"`
	Sell12h                      int64              `json:"sell12h"`
	SellHistory12h               int64              `json:"sellHistory12h"`
	Sell12hChangePercent         float64            `json:"sell12hChangePercent"`
	Buy12h                       int64              `json:"buy12h"`
	BuyHistory12h                int64              `json:"buyHistory12h"`
	Buy12hChangePercent          float64            `json:"buy12hChangePercent"`
	V12h                         float64            `json:"v12h"`
	V12hUSD                      float64            `json:"v12hUSD"`
	VHistory12h                  float64            `json:"vHistory12h"`
	VHistory12hUSD               float64            `json:"vHistory12hUSD"`
	V12hChangePercent            float64            `json:"v12hChangePercent"`
	VBuy12h                      float64            `json:"vBuy12h"`
	VBuy12hUSD                   float64            `json:"vBuy12hUSD"`
	VBuyHistory12h               float64            `json:"vBuyHistory12h"`
	VBuyHistory12hUSD            float64            `json:"vBuyHistory12hUSD"`
	VBuy12hChangePercent         float64            `json:"vBuy12hChangePercent"`
	VSell12h                     float64            `json:"vSell12h"`
	VSell12hUSD                  float64            `json:"vSell12hUSD"`
	VSellHistory12h              float64            `json:"vSellHistory12h"`
	VSellHistory12hUSD           float64            `json:"vSellHistory12hUSD"`
	VSell12hChangePercent        float64            `json:"vSell12hChangePercent"`
	Trade24h                     int64              `json:"trade24h"`
	TradeHistory24h              int64              `json:"tradeHistory24h"`
	Trade24hChangePercent        float64            `json:"trade24hChangePercent"`
	Sell24h                      int64              `json:"sell24h"`
	SellHistory24h               int64              `json:"sellHistory24h"`
	Sell24hChangePercent         float64            `json:"sell24hChangePercent"`
	Buy24h                       int64              `json:"buy24h"`
	BuyHistory24h                int64              `json:"buyHistory24h"`
	Buy24hChangePercent          float64            `json:"buy24hChangePercent"`
	V24h                         float64            `json:"v24h"`
	V24hUSD                      float64            `json:"v24hUSD"`
	VHistory24h                  float64            `json:"vHistory24h"`
	VHistory24hUSD               float64            `json:"vHistory24hUSD"`
	V24hChangePercent            float64            `json:"v24hChangePercent"`
	VBuy24h                      float64            `json:"vBuy24h"`
	VBuy24hUSD                   float64            `json:"vBuy24hUSD"`
	VBuyHistory24h               float64            `json:"vBuyHistory24h"`
	VBuyHistory24hUSD            float64            `json:"vBuyHistory24hUSD"`
	VBuy24hChangePercent         float64            `json:"vBuy24hChangePercent"`
	VSell24h                     float64            `json:"vSell24h"`
	VSell24hUSD                  float64            `json:"vSell24hUSD"`
	VSellHistory24h              float64            `json:"vSellHistory24h"`
	VSellHistory24hUSD           float64            `json:"vSellHistory24hUSD"`
	VSell24hChangePercent        float64            `json:"vSell24hChangePercent"`
	Watch                        int64              `json:"watch"`
	View30m                      int64              `json:"view30m"`
	ViewHistory30m               int64              `json:"viewHistory30m"`
	View30mChangePercent         float64            `json:"view30mChangePercent"`
	View1h                       int64              `json:"view1h"`
	ViewHistory1h                int64              `json:"viewHistory1h"`
	View1hChangePercent          float64            `json:"view1hChangePercent"`
	View2h                       int64              `json:"view2h"`
	ViewHistory2h                int64              `json:"viewHistory2h"`
	View2hChangePercent          float64            `json:"view2hChangePercent"`
	View4h                       int64              `json:"view4h"`
	ViewHistory4h                int64              `json:"viewHistory4h"`
	View4hChangePercent          float64            `json:"view4hChangePercent"`
	View6h                       int64              `json:"view6h"`
	ViewHistory6h                int64              `json:"viewHistory6h"`
	View6hChangePercent          float64            `json:"view6hChangePercent"`
	View8h                       int64              `json:"view8h"`
	ViewHistory8h                int64              `json:"viewHistory8h"`
	View8hChangePercent          float64            `json:"view8hChangePercent"`
	View12h                      int64              `json:"view12h"`
	ViewHistory12h               int64              `json:"viewHistory12h"`
	View12hChangePercent         float64            `json:"view12hChangePercent"`
	View24h                      int64              `json:"view24h"`
	ViewHistory24h               int64              `json:"viewHistory24h"`
	View24hChangePercent         float64            `json:"view24hChangePercent"`
	UniqueView30m                int64              `json:"uniqueView30m"`
	UniqueViewHistory30m         int64              `json:"uniqueViewHistory30m"`
	UniqueView30mChangePercent   float64            `json:"uniqueView30mChangePercent"`
	UniqueView1h                 int64              `json:"uniqueView1h"`
	UniqueViewHistory1h          int64              `json:"uniqueViewHistory1h"`
	UniqueView1hChangePercent    float64            `json:"uniqueView1hChangePercent"`
	UniqueView2h                 int64              `json:"uniqueView2h"`
	UniqueViewHistory2h          int64              `json:"uniqueViewHistory2h"`
	UniqueView2hChangePercent    float64            `json:"uniqueView2hChangePercent"`
	UniqueView4h                 int64              `json:"uniqueView4h"`
	UniqueViewHistory4h          int64              `json:"uniqueViewHistory4h"`
	UniqueView4hChangePercent    float64            `json:"uniqueView4hChangePercent"`
	UniqueView6h                 int64              `json:"uniqueView6h"`
	UniqueViewHistory6h          int64              `json:"uniqueViewHistory6h"`
	UniqueView6hChangePercent    float64            `json:"uniqueView6hChangePercent"`
	UniqueView8h                 int64              `json:"uniqueView8h"`
	UniqueViewHistory8h          int64              `json:"uniqueViewHistory8h"`
	UniqueView8hChangePercent    float64            `json:"uniqueView8hChangePercent"`
	UniqueView12h                int64              `json:"uniqueView12h"`
	UniqueViewHistory12h         int64              `json:"uniqueViewHistory12h"`
	UniqueView12hChangePercent   float64            `json:"uniqueView12hChangePercent"`
	UniqueView24h                int64              `json:"uniqueView24h"`
	UniqueViewHistory24h         int64              `json:"uniqueViewHistory24h"`
	UniqueView24hChangePercent   float64            `json:"uniqueView24hChangePercent"`
	NumberMarkets                int64              `json:"numberMarkets"`
}

type RespToken struct {
	Address           string  `json:"address"`
	Decimals          int64   `json:"decimals"`
	Liquidity         float64 `json:"liquidity"`
	Mc                float64 `json:"mc"`
	Symbol            string  `json:"symbol"`
	V24hChangePercent float64 `json:"v24hChangePercent"`
	V24hUSD           float64 `json:"v24hUSD"`
	Name              string  `json:"name"`
	LastTradeUnixTime int64   `json:"lastTradeUnixTime"`
}

type RespTokenSecurity struct {
	CreatorAddress                 *string  `json:"creatorAddress"`
	OwnerAddress                   *string  `json:"ownerAddress"`
	CreationTx                     *string  `json:"creationTx"`
	CreationTime                   *int64   `json:"creationTime"`
	CreationSlot                   *int64   `json:"creationSlot"`
	MintTx                         *string  `json:"mintTx"`
	MintTime                       *int64   `json:"mintTime"`
	MintSlot                       *int64   `json:"mintSlot"`
	CreatorBalance                 *float64 `json:"creatorBalance"`
	OwnerBalance                   *float64 `json:"ownerBalance"`
	OwnerPercentage                *float64 `json:"ownerPercentage"`
	CreatorPercentage              *float64 `json:"creatorPercentage"`
	MetaplexUpdateAuthority        string   `json:"metaplexUpdateAuthority"`
	MetaplexUpdateAuthorityBalance float64  `json:"metaplexUpdateAuthorityBalance"`
	MetaplexUpdateAuthorityPercent float64  `json:"metaplexUpdateAuthorityPercent"`
	MutableMetadata                bool     `json:"mutableMetadata"`
	Top10HolderBalance             float64  `json:"top10HolderBalance"`
	Top10HolderPercent             float64  `json:"top10HolderPercent"`
	Top10UserBalance               float64  `json:"top10UserBalance"`
	Top10UserPercent               float64  `json:"top10UserPercent"`
	IsTrueToken                    *bool    `json:"isTrueToken"`
	TotalSupply                    float64  `json:"totalSupply"`
	PreMarketHolder                []string `json:"preMarketHolder"`
	LockInfo                       *string  `json:"lockInfo"`
	Freezeable                     *bool    `json:"freezeable"`
	FreezeAuthority                *string  `json:"freezeAuthority"`
	TransferFeeEnable              *bool    `json:"transferFeeEnable"`
	TransferFeeData                *string  `json:"transferFeeData"`
	IsToken2022                    bool     `json:"isToken2022"`
	NonTransferable                *bool    `json:"nonTransferable"`
}

type RespTokenCreationInfo struct {
	TxHash         string `json:"txHash"`
	Slot           int64  `json:"slot"`
	TokenAddress   string `json:"tokenAddress"`
	Decimals       int64  `json:"decimals"`
	Owner          string `json:"owner"`
	BlockUnixTime  int64  `json:"blockUnixTime"`
	BlockHumanTime string `json:"blockHumanTime"`
}

type RespMarketTokenInfo struct {
	Address  string `json:"address"`
	Decimals int64  `json:"decimals"`
	Icon     string `json:"icon"`
	Symbol   string `json:"symbol"`
}

type RespMarketItem struct {
	Address                      string              `json:"address"`
	Base                         RespMarketTokenInfo `json:"base"`
	CreatedAt                    string              `json:"createdAt"`
	Name                         string              `json:"name"`
	Quote                        RespMarketTokenInfo `json:"quote"`
	Source                       string              `json:"source"`
	Liquidity                    float64             `json:"liquidity"`
	LiquidityChangePercentage24h *float64            `json:"liquidityChangePercentage24h"`
	Price                        float64             `json:"price"`
	Trade24h                     int64               `json:"trade24h"`
	Trade24hChangePercent        float64             `json:"trade24hChangePercent"`
	UniqueWallet24h              int64               `json:"uniqueWallet24h"`
	UniqueWallet24hChangePercent float64             `json:"uniqueWallet24hChangePercent"`
	Volume24h                    float64             `json:"volume24h"`
	Volume24hChangePercentage24h *float64            `json:"volume24hChangePercentage24h"`
}

type RespNewTokenListingItem struct {
	Address          string  `json:"address"`
	Symbol           string  `json:"symbol"`
	Name             string  `json:"name"`
	Decimals         int64   `json:"decimals"`
	LiquidityAddedAt int64   `json:"liquidityAddedAt"`
	Liquidity        float64 `json:"liquidity"`
}

type RespTopTraderItem struct {
	Owner        string   `json:"owner"`
	TokenAddress string   `json:"tokenAddress"`
	Trade        int64    `json:"trade"`
	TradeBuy     int64    `json:"tradeBuy"`
	TradeSell    int64    `json:"tradeSell"`
	Type         string   `json:"type"`
	Volume       float64  `json:"volume"`
	VolumeBuy    float64  `json:"volumeBuy"`
	VolumeSell   float64  `json:"volumeSell"`
	Tags         []string `json:"tags"`
}

type RespWalletBalanceChange struct {
	Amount   int64  `json:"amount"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int64  `json:"decimals"`
	Address  string `json:"address"`
	LogoURI  string `json:"logoURI"`
}

type RespContractLabel struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Metadata struct {
		Icon string `json:"icon"`
	} `json:"metadata"`
}

type RespWalletPortfolioItem struct {
	Address  string  `json:"address"`
	Decimals int64   `json:"decimals"`
	Balance  int64   `json:"balance"`
	UiAmount float64 `json:"uiAmount"`
	ChainId  string  `json:"chainId"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	LogoURI  string  `json:"logoURI"`
	PriceUsd float64 `json:"priceUsd"`
	ValueUsd float64 `json:"valueUsd"`
}

type RespWalletPortfolio struct {
	Wallet   string                    `json:"wallet"`
	TotalUsd float64                   `json:"totalUsd"`
	Items    []RespWalletPortfolioItem `json:"items"`
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
		header.Set("chain", strings.Join(chains, ","))
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
		case bool:
			if v {
				value = "True"
			} else {
				value = "False"
			}
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
		return *new(D), err
	}
	req.Header = clt.newHeader(chains...)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return *new(D), err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	err = statusCodeToError[statusCode]
	if err != nil {
		return *new(D), err
	}

	var rd RespData[D]
	if err := json.NewDecoder(resp.Body).Decode(&rd); err != nil {
		return *new(D), err
	}

	if statusCode == http.StatusOK {
		return rd.Data, nil
	}

	return *new(D), fmt.Errorf("birdeye error: status code: %d, message: %s", statusCode, rd.Message)
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

func (c *Client) PriceHistory(chain string, address string, addressType string, chartType ChartType, timeFrom, timeTo int64) ([]RespOHLCVItem, error) {
	return get[[]RespOHLCVItem](context.Background(), c, "/defi/history_price", []string{chain},
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
func (c *Client) OHLCVByToken(chain string, address string, chartType ChartType, timeFrom, timeTo int64) ([]RespOHLCVItem, error) {
	return get[[]RespOHLCVItem](context.Background(), c, "/defi/ohlcv", []string{chain},
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
func (c *Client) OHLCVByPair(chain string, address string, chartType ChartType, timeFrom, timeTo int64) ([]RespOHLCVBaseQuoteItem, error) {
	return get[[]RespOHLCVBaseQuoteItem](context.Background(), c, "/defi/ohlcv_pair", []string{chain},
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
func (c *Client) OHLCVByBaseQuote(chain string, baseAddress string, quoteAddress string, chartType ChartType, timeFrom, timeTo int64) ([]RespOHLCVBaseQuoteItem, error) {
	return get[[]RespOHLCVBaseQuoteItem](context.Background(), c, "/defi/ohlcv/base_quote", []string{chain},
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
func (c *Client) TradesByToken(chain string, address string, sortType string, offset int, limit int, txType string) (RespItems[RespTradesByTokenItem], error) {
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
func (c *Client) TradesByPair(chain string, address string, sortType string, offset int, limit int, txType string) (RespItems[RespTradesByPairItem], error) {
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
func (c *Client) PriceVolumeByToken(chain string, address string, timeType string) (RespSinglePriceVolume, error) {
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
func (c *Client) PriceVolumeByTokens(chain string, listAddress string, timeType string) ([]RespSinglePriceVolume, error) {
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
func (c *Client) TrendingTokens(chain string, sortBy string, sortType string, offset int, limit int) (RespTrendingTokens, error) {
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
func (c *Client) TradeByTokenAndTime(chain string, address string, beforeTime, afterTime int64, txType string, offset int, limit int) ([]RespTradesByTokenItem, error) {
	if beforeTime > 0 && afterTime > 0 {
		return []RespTradesByTokenItem{}, fmt.Errorf("beforeTime and afterTime cannot be used simultaneously")
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
		return []RespTradesByTokenItem{}, err
	}
	return d.Items, nil
}

// TradesSeekByTime retrieves transaction records for a specific pair address based on Unix time.
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
func (c *Client) TradesSeekByTime(chain string, address string, beforeTime, afterTime int64, txType string, offset int, limit int) ([]RespTradesByPairItem, error) {
	if beforeTime > 0 && afterTime > 0 {
		return []RespTradesByPairItem{}, fmt.Errorf("beforeTime and afterTime cannot be used simultaneously (error 422)")
	}

	if beforeTime == 0 && afterTime == 0 {
		return []RespTradesByPairItem{}, fmt.Errorf("either beforeTime or afterTime must be specified")
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
		return []RespTradesByPairItem{}, err
	}
	return d.Items, nil
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
func (c *Client) TokenList(chain string, sortType string, sortOrder string, offset int, limit int, minLiquidity float64) (RespItems[RespToken], error) {
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

	return get[RespItems[RespToken]](context.Background(), c, "/defi/tokenlist", []string{chain},
		"sort_type", sortType,
		"sort_order", sortOrder,
		"offset", offset,
		"limit", limit,
		"min_liquidity", minLiquidity)
}

// TokenListV2 retrieves a list of tokens using the v2 API endpoint
//
// Parameters:
//   - chain: The blockchain network
//   - sortBy: Attribute to sort tokens by ("v24hUSD", "mc", "v24hChangePercent", default: "v24hUSD")
//   - sortOrder: Sort order ("asc" or "desc", default: "desc")
//   - offset: Starting index for the list (0-1000, default: 0)
//   - limit: Maximum number of tokens to retrieve (1-50, default: 50)
//   - minLiquidity: Minimum liquidity filter (default: 100)
//
// Returns:
//   - RespItems[RespToken]: Paginated list of token information
//   - error: Any error that occurred during the request
func (c *Client) TokenListV2(chain string, sortBy string, sortOrder string, offset int, limit int, minLiquidity float64) (RespItems[RespToken], error) {
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

	return get[RespItems[RespToken]](context.Background(), c, "/defi/v2/tokens/all", []string{chain},
		"sort_by", sortBy,
		"sort_order", sortOrder,
		"offset", offset,
		"limit", limit,
		"min_liquidity", minLiquidity)
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
func (c *Client) TokenCreationInfo(chain string, address string) (RespTokenSecurity, error) {
	return get[RespTokenSecurity](context.Background(), c, "/defi/token_creation_info", []string{chain},
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
func (c *Client) MarketList(chain string, address string, sortBy string, sortType string, offset int, limit int) (RespItems[RespToken], error) {
	if limit > 10 {
		limit = 10
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return get[RespItems[RespToken]](context.Background(), c, "/defi/v2/markets", []string{chain},
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
func (c *Client) NewTokenListing(chain string, timeTo int64, limit int, memePlatformEnabled bool) (RespItems[RespToken], error) {
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

	return get[RespItems[RespToken]](context.Background(), c, "/defi/v2/tokens/new_listing", []string{chain}, params...)
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
func (c *Client) TokenTopTraders(chain string, address string, sortBy string, sortType string, timeFrame string, offset int, limit int) (RespItems[RespToken], error) {
	if limit > 10 {
		limit = 10
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return get[RespItems[RespToken]](context.Background(), c, "/defi/v2/tokens/top_traders", []string{chain},
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
func (c *Client) WalletTxHistories(chain string, wallet string, limit int, before string) ([]RespTradesByTokenItem, error) {
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

	d, err := get[RespItems[RespTradesByTokenItem]](context.Background(), c, "/v1/wallet/tx_list", []string{chain}, params...)
	if err != nil {
		return []RespTradesByTokenItem{}, err
	}
	return d.Items, nil
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
func (c *Client) WalletPortfolio(chain string, wallet string) ([]RespToken, error) {
	d, err := get[RespItems[RespToken]](context.Background(), c, "/v1/wallet/token_list", []string{chain}, "wallet", wallet)
	if err != nil {
		return []RespToken{}, err
	}
	return d.Items, nil
}
