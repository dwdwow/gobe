package gobe

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsQueryType string

const (
	QUERY_TYPE_SIMPLE  WsQueryType = "simple"
	QUERY_TYPE_COMPLEX WsQueryType = "complex"
)

type WsSubType string

const (
	SUBSCRIBE_PRICE             WsSubType = "SUBSCRIBE_PRICE"
	SUBSCRIBE_TXS               WsSubType = "SUBSCRIBE_TXS"
	SUBSCRIBE_BASE_QUOTE_PRICE  WsSubType = "SUBSCRIBE_BASE_QUOTE_PRICE"
	SUBSCRIBE_TOKEN_NEW_LISTING WsSubType = "SUBSCRIBE_TOKEN_NEW_LISTING"
	SUBSCRIBE_NEW_PAIR          WsSubType = "SUBSCRIBE_NEW_PAIR"
	SUBSCRIBE_LARGE_TRADE_TXS   WsSubType = "SUBSCRIBE_LARGE_TRADE_TXS"
	SUBSCRIBE_WALLET_TXS        WsSubType = "SUBSCRIBE_WALLET_TXS"

	UNSUBSCRIBE_PRICE             WsSubType = "UNSUBSCRIBE_PRICE"
	UNSUBSCRIBE_TXS               WsSubType = "UNSUBSCRIBE_TXS"
	UNSUBSCRIBE_BASE_QUOTE_PRICE  WsSubType = "UNSUBSCRIBE_BASE_QUOTE_PRICE"
	UNSUBSCRIBE_TOKEN_NEW_LISTING WsSubType = "UNSUBSCRIBE_TOKEN_NEW_LISTING"
	UNSUBSCRIBE_NEW_PAIR          WsSubType = "UNSUBSCRIBE_NEW_PAIR"
	UNSUBSCRIBE_LARGE_TRADE_TXS   WsSubType = "UNSUBSCRIBE_LARGE_TRADE_TXS"
	UNSUBSCRIBE_WALLET_TXS        WsSubType = "UNSUBSCRIBE_WALLET_TXS"
)

type WsDataType string

const (
	WS_WELCOME_DATA           WsDataType = "WELCOME"
	WS_ERROR_DATA             WsDataType = "ERROR"
	WS_PRICE_DATA             WsDataType = "PRICE_DATA"
	WS_TXS_DATA               WsDataType = "TXS_DATA"
	WS_BASE_QUOTE_PRICE_DATA  WsDataType = "BASE_QUOTE_PRICE_DATA"
	WS_TOKEN_NEW_LISTING_DATA WsDataType = "TOKEN_NEW_LISTING_DATA"
	WS_NEW_PAIR_DATA          WsDataType = "NEW_PAIR_DATA"
	WS_TXS_LARGE_TRADE_DATA   WsDataType = "TXS_LARGE_TRADE_DATA"
	WS_WALLET_TXS_DATA        WsDataType = "WALLET_TXS_DATA"
)

type WsCurrency string

const (
	WS_CURRENCY_USD  WsCurrency = "usd"
	WS_CURRENCY_PAIR WsCurrency = "pair"
)

type WsSubData[D any] struct {
	Type WsSubType `json:"type"`
	Data D         `json:"data"`
}

type WsComplexSubData struct {
	QueryType WsQueryType `json:"queryType"`
	Query     string      `json:"query"`
}

func JoinQuery(queries ...string) string {
	return strings.Join(queries, " OR ")
}

type WsPriceSubData struct {
	// QueryType: "simple"
	QueryType WsQueryType `json:"queryType"`
	// ChartType: "1m", "3m"...
	ChartType ChartType `json:"chartType"`
	Address   string    `json:"address"`
	// Currency: "usd"...
	Currency WsCurrency `json:"currency"`
}

func (d WsPriceSubData) Query() string {
	return fmt.Sprintf("(address = %s AND chartType = %s AND currency = %s)", d.Address, d.ChartType, d.Currency)
}

type WsPriceData struct {
	// Open price
	O float64 `json:"o" bson:"o"`
	// High price
	H float64 `json:"h" bson:"h"`
	// Low price
	L float64 `json:"l" bson:"l"`
	// Close price
	C float64 `json:"c" bson:"c"`
	// Event type (e.g. "ohlcv")
	EventType string `json:"eventType" bson:"eventType"`
	// Chart type (e.g. "1m")
	Type ChartType `json:"type" bson:"type"`
	// UnixTime seconds
	UnixTime int64 `json:"unixTime" bson:"unixTime"`
	// Volume
	V float64 `json:"v" bson:"v"`
	// Token/Pair symbol, e.g. "SOL"/"SOL-USDC"
	Symbol string `json:"symbol" bson:"symbol"`
	// Token/Pair address
	Address string `json:"address" bson:"address"`
}

type WsTxsSubData struct {
	// QueryType: "simple"
	QueryType WsQueryType `json:"queryType" bson:"queryType"`
	// set token address or pair address, not both
	// Token address
	Address string `json:"address,omitempty" bson:"address,omitempty"`
	// Pair address
	PairAddress string `json:"pairAddress,omitempty" bson:"pairAddress,omitempty"`
}

func (d WsTxsSubData) Query() string {
	if d.Address != "" {
		return fmt.Sprintf("address = %s", d.Address)
	}
	return fmt.Sprintf("pairAddress = %s", d.PairAddress)
}

type WsTxTokenInfo struct {
	// Token symbol
	Symbol string `json:"symbol" bson:"symbol"`
	// Token decimals
	Decimals int `json:"decimals" bson:"decimals"`
	// Token address
	Address string `json:"address" bson:"address"`
	// Raw amount
	Amount any `json:"amount" bson:"amount"`
	// Transaction type
	Type string `json:"type" bson:"type"`
	// Swap type (from/to)
	TypeSwap string `json:"typeSwap" bson:"typeSwap"`
	// UI formatted amount
	UiAmount float64 `json:"uiAmount" bson:"uiAmount"`
	// Token price
	Price float64 `json:"price,omitempty" bson:"price,omitempty"`
	// Nearest price if price is not available
	NearestPrice float64 `json:"nearestPrice" bson:"nearestPrice"`
	// Raw change amount
	ChangeAmount any `json:"changeAmount" bson:"changeAmount"`
	// UI formatted change amount
	UiChangeAmount float64 `json:"uiChangeAmount" bson:"uiChangeAmount"`
	// Token icon URL
	Icon string `json:"icon,omitempty" bson:"icon,omitempty"`
}

type WsTxsData struct {
	// Block unix time
	BlockUnixTime int64 `json:"blockUnixTime" bson:"blockUnixTime"`
	// Owner address
	Owner string `json:"owner" bson:"owner"`
	// Source platform
	Source string `json:"source" bson:"source"`
	// Transaction hash
	TxHash string `json:"txHash" bson:"txHash"`
	// Alias name
	Alias *string `json:"alias" bson:"alias"`
	// Whether trade happened on Birdeye
	IsTradeOnBe bool `json:"isTradeOnBe" bson:"isTradeOnBe"`
	// Platform address
	Platform string `json:"platform" bson:"platform"`
	// Volume in USD
	VolumeUSD float64 `json:"volumeUSD" bson:"volumeUSD"`
	// From token info
	From WsTxTokenInfo `json:"from" bson:"from"`
	// To token info
	To WsTxTokenInfo `json:"to" bson:"to"`
}

type WsBaseQuotePriceSubData struct {
	// Base token address
	BaseAddress string `json:"baseAddress" bson:"baseAddress"`
	// Quote token address
	QuoteAddress string `json:"quoteAddress" bson:"quoteAddress"`
	// Chart type/interval
	ChartType ChartType `json:"chartType" bson:"chartType"`
}

type WsBaseQuotePriceData struct {
	// Open price
	O float64 `json:"o" bson:"o"`
	// High price
	H float64 `json:"h" bson:"h"`
	// Low price
	L float64 `json:"l" bson:"l"`
	// Close price
	C float64 `json:"c" bson:"c"`
	// Event type
	EventType string `json:"eventType" bson:"eventType"`
	// Chart type/interval
	Type string `json:"type" bson:"type"`
	// Unix timestamp
	UnixTime int64 `json:"unixTime" bson:"unixTime"`
	// Volume is 0
	V float64 `json:"v" bson:"v"`
	// Base token address
	BaseAddress string `json:"baseAddress" bson:"baseAddress"`
	// Quote token address
	QuoteAddress string `json:"quoteAddress" bson:"quoteAddress"`
}

// WsTokenNewListingSubData represents subscription data for new token listing notifications
//
// Parameters:
//   - MemePlatformEnabled: Optional. Set to true to receive new meme token listings from platforms like pump.fun.
//     If not set, no listings from meme platforms will be received.
//   - MinLiquidity: Optional. Minimum liquidity threshold for notifications. Must be set higher than system minimum of 10.
//   - MaxLiquidity: Optional. Maximum liquidity threshold for notifications. When provided, must be higher than MinLiquidity.
type WsTokenNewListingSubData struct {
	// Whether meme platform is enabled
	MemePlatformEnabled bool `json:"meme_platform_enabled,omitempty" bson:"meme_platform_enabled,omitempty"`
	// Minimum liquidity
	MinLiquidity float64 `json:"min_liquidity,omitempty" bson:"min_liquidity,omitempty"`
	// Maximum liquidity
	MaxLiquidity float64 `json:"max_liquidity,omitempty" bson:"max_liquidity,omitempty"`
}

type WsTokenNewListingData struct {
	// Token address
	Address string `json:"address" bson:"address"`
	// Token decimals
	Decimals int64 `json:"decimals" bson:"decimals"`
	// Token name
	Name string `json:"name" bson:"name"`
	// Token symbol
	Symbol string `json:"symbol" bson:"symbol"`
	// Token liquidity in USD
	Liquidity float64 `json:"liquidity" bson:"liquidity"`
	// Unix timestamp when liquidity was added
	LiquidityAddedAt string `json:"liquidityAddedAt" bson:"liquidityAddedAt"`
}

type WsNewPairSubData struct {
	// Minimum liquidity
	MinLiquidity float64 `json:"min_liquidity,omitempty" bson:"min_liquidity,omitempty"`
	// Maximum liquidity
	MaxLiquidity float64 `json:"max_liquidity,omitempty" bson:"max_liquidity,omitempty"`
}

type WsNewPairTokenInfo struct {
	// Token address
	Address string `json:"address" bson:"address"`
	// Token name
	Name string `json:"name" bson:"name"`
	// Token symbol
	Symbol string `json:"symbol" bson:"symbol"`
	// Token decimals
	Decimals int64 `json:"decimals" bson:"decimals"`
}

type WsNewPairData struct {
	// Pair address
	Address string `json:"address" bson:"address"`
	// Pair name
	Name string `json:"name" bson:"name"`
	// Source of the pair
	Source string `json:"source" bson:"source"`
	// Base token info
	Base WsNewPairTokenInfo `json:"base" bson:"base"`
	// Quote token info
	Quote WsNewPairTokenInfo `json:"quote" bson:"quote"`
	// Transaction hash
	TxHash string `json:"txHash" bson:"txHash"`
	// Block time
	BlockTime string `json:"blockTime" bson:"blockTime"`
}

// WsLargeTradeTxsSubData represents the subscription data for large trade transactions.
//
// The min_volume parameter is mandatory and sets the lower bound for trade volume in USD.
// It must be at least 1000 USD.
//
// The max_volume parameter is optional but when provided must be greater than min_volume.
// Trades with volumes outside the specified range will be filtered out.
//
// The subscription will return trades for all tokens meeting the volume criteria,
// regardless of the specific tokens or trading pairs involved.
type WsLargeTradeTxsSubData struct {
	// SubType: "SUBSCRIBE_LARGE_TRADE_TXS"
	Type WsSubType `json:"type" bson:"type"`
	// Minimum volume in USD
	MinVolume float64 `json:"min_volume" bson:"min_volume"`
	// Maximum volume in USD
	MaxVolume float64 `json:"max_volume,omitempty" bson:"max_volume,omitempty"`
}

type WsLargeTradeTxsTokenInfo struct {
	// Token address
	Address string `json:"address" bson:"address"`
	// Token name
	Name string `json:"name" bson:"name"`
	// Token symbol
	Symbol string `json:"symbol" bson:"symbol"`
	// Token decimals
	Decimals int64 `json:"decimals" bson:"decimals"`
	// Token amount in UI format
	UiAmount float64 `json:"uiAmount" bson:"uiAmount"`
	// Token price
	Price float64 `json:"price" bson:"price"`
	// Nearest token price
	NearestPrice float64 `json:"nearestPrice" bson:"nearestPrice"`
	// Token amount change in UI format
	UiChangeAmount float64 `json:"uiChangeAmount" bson:"uiChangeAmount"`
}

type WsLargeTradeTxsData struct {
	// Block unix timestamp seconds
	BlockUnixTime int64 `json:"blockUnixTime" bson:"blockUnixTime"`
	// Block time in human readable format
	BlockHumanTime string `json:"blockHumanTime" bson:"blockHumanTime"`
	// Owner address
	Owner string `json:"owner" bson:"owner"`
	// Source of the trade
	Source string `json:"source" bson:"source"`
	// Pool address
	PoolAddress string `json:"poolAddress" bson:"poolAddress"`
	// Transaction hash
	TxHash string `json:"txHash" bson:"txHash"`
	// Volume in USD
	VolumeUSD float64 `json:"volumeUSD" bson:"volumeUSD"`
	// Network name
	Network string `json:"network" bson:"network"`
	// From token info
	From WsLargeTradeTxsTokenInfo `json:"from" bson:"from"`
	// To token info
	To WsLargeTradeTxsTokenInfo `json:"to" bson:"to"`
}

type WsWalletTxsSubData struct {
	// Wallet address to monitor
	Address string `json:"address" bson:"address"`
}

type WsWalletTxsTokenInfo struct {
	// Token symbol
	Symbol string `json:"symbol" bson:"symbol"`
	// Token decimals
	Decimals int64 `json:"decimals" bson:"decimals"`
	// Token address
	Address string `json:"address" bson:"address"`
	// Token amount in UI format
	UiAmount float64 `json:"uiAmount" bson:"uiAmount"`
}

type WsWalletTxsData struct {
	// Transaction type
	Type string `json:"type" bson:"type"`
	// Block unix timestamp seconds
	BlockUnixTime int64 `json:"blockUnixTime" bson:"blockUnixTime"`
	// Block time in human readable format
	BlockHumanTime string `json:"blockHumanTime" bson:"blockHumanTime"`
	// Owner address
	Owner string `json:"owner" bson:"owner"`
	// Source address
	Source string `json:"source" bson:"source"`
	// Transaction hash
	TxHash string `json:"txHash" bson:"txHash"`
	// Volume in USD
	VolumeUSD float64 `json:"volumeUSD" bson:"volumeUSD"`
	// Network name
	Network string `json:"network" bson:"network"`
	// Base token info
	Base WsWalletTxsTokenInfo `json:"base" bson:"base"`
	// Quote token info
	Quote WsWalletTxsTokenInfo `json:"quote" bson:"quote"`
}

type WsClient struct {
	url string
	ws  *websocket.Conn

	muReConn sync.RWMutex

	muSubers sync.RWMutex
	subers   map[WsDataType][]chan any

	muRW sync.Mutex

	logger *slog.Logger
}

func NewWsClient(chain, apiKey string, logger *slog.Logger) *WsClient {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	if apiKey == "" {
		panic("birdeye: api key is required")
	}
	url := fmt.Sprintf("wss://public-api.birdeye.so/socket/%s?x-api-key=%s", chain, apiKey)
	return &WsClient{url: url, subers: make(map[WsDataType][]chan any), logger: logger}
}

func (c *WsClient) Start() error {
	headers := http.Header{}
	headers.Add("Origin", "ws://public-api.birdeye.so")
	headers.Add("Sec-WebSocket-Origin", "ws://public-api.birdeye.so")
	headers.Add("Sec-WebSocket-Protocol", "echo-protocol")
	conn, reps, err := websocket.DefaultDialer.Dial(c.url, headers)
	if err != nil {
		return fmt.Errorf("birdeye: failed to connect to websocket: %w, http status code: %d", err, reps.StatusCode)
	}
	c.ws = conn
	go c.waiter()
	return nil
}

func (c *WsClient) reConn() {
	if !c.muReConn.TryLock() {
		return
	}
	defer c.muReConn.Unlock()
	for {
		c.logger.Info("birdeye: retrying to connect to websocket...")
		err := c.Start()
		if err != nil {
			c.logger.Error("birdeye: failed to connect to websocket, retrying...", "error", err)
			time.Sleep(5 * time.Second)
		} else {
			c.logger.Info("birdeye: reconnected to websocket")
			return
		}
	}
}

func (c *WsClient) Close() error {
	return c.ws.Close()
}

func (c *WsClient) waiter() {
	for {
		t, b, err := c.ws.ReadMessage()
		if err != nil {
			c.logger.Error("birdeye: websocket read error", "error", err)
			c.reConn()
			continue
		}
		switch t {
		case websocket.BinaryMessage:
			c.logger.Info("birdeye: websocket binary message", "data", string(b))
		case websocket.PingMessage:
			c.logger.Info("birdeye: websocket ping message", "data", string(b))
		case websocket.PongMessage:
			c.logger.Info("birdeye: websocket pong message", "data", string(b))
		case websocket.TextMessage:
			go c.msgHandler(b)
		case websocket.CloseMessage:
			c.logger.Info("birdeye: websocket close message", "data", string(b))
		}
	}
}

func (c *WsClient) msgHandler(b []byte) {
	d := map[string]any{}
	err := json.Unmarshal(b, &d)
	if err != nil {
		c.logger.Error("birdeye: failed to unmarshal message", "error", err)
		return
	}
	t, ok := d["type"].(string)
	if !ok {
		c.logger.Error("birdeye: message type is not string", "data", string(b))
		return
	}
	b, err = json.Marshal(d["data"])
	if err != nil {
		c.logger.Error("birdeye: failed to marshal data", "error", err)
		return
	}
	var dd any
	switch WsDataType(t) {
	case WS_WELCOME_DATA:
		c.logger.Info("birdeye: welcome message", "data", string(b))
		return
	case WS_ERROR_DATA:
		c.logger.Error("birdeye: error message", "data", string(b))
		return
	case WS_PRICE_DATA:
		dd = &WsPriceData{}
	case WS_TXS_DATA:
		dd = &WsTxsData{}
	case WS_BASE_QUOTE_PRICE_DATA:
		dd = &WsBaseQuotePriceData{}
	case WS_TOKEN_NEW_LISTING_DATA:
		dd = &WsTokenNewListingData{}
	case WS_NEW_PAIR_DATA:
		dd = &WsNewPairData{}
	case WS_TXS_LARGE_TRADE_DATA:
		dd = &WsLargeTradeTxsData{}
	case WS_WALLET_TXS_DATA:
		dd = &WsWalletTxsData{}
	default:
		c.logger.Error("birdeye: unknown message type", "type", t, "data", string(b))
		return
	}
	err = json.Unmarshal(b, dd)
	if err != nil {
		c.logger.Error("birdeye: failed to unmarshal data", "error", err)
		return
	}
	c.muSubers.RLock()
	defer c.muSubers.RUnlock()
	subers := c.subers[WsDataType(t)]
	for _, suber := range subers {
		suber := suber
		go func() {
			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()
			select {
			case suber <- dd:
			case <-timer.C:
				c.logger.Error("birdeye: failed to send data to suber", "type", t)
			}
		}()
	}
}

func (c *WsClient) WsSub(d any) error {
	c.muRW.Lock()
	defer c.muRW.Unlock()
	return c.ws.WriteJSON(d)
}

func (c *WsClient) NewDataChan(t WsDataType) <-chan any {
	c.muSubers.Lock()
	defer c.muSubers.Unlock()
	ch := make(chan any, 100)
	c.subers[t] = append(c.subers[t], ch)
	return ch
}

// func (c *WsClient) Unsubscribe(t WsSubType, ch <-chan any) {
// 	c.muSubers.Lock()
// 	defer c.muSubers.Unlock()
// 	for i, suber := range c.subers[t] {
// 		if suber != ch {
// 			continue
// 		}
// 		close(suber)
// 		if i == len(c.subers[t])-1 {
// 			c.subers[t] = c.subers[t][:i]
// 		} else {
// 			c.subers[t] = append(c.subers[t][:i], c.subers[t][i+1:]...)
// 		}
// 		break
// 	}
// }
