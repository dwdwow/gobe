package gobe

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	PRICE_DATA             WsDataType = "PRICE_DATA"
	TXS_DATA               WsDataType = "TXS_DATA"
	BASE_QUOTE_PRICE_DATA  WsDataType = "BASE_QUOTE_PRICE_DATA"
	TOKEN_NEW_LISTING_DATA WsDataType = "TOKEN_NEW_LISTING_DATA"
	NEW_PAIR_DATA          WsDataType = "NEW_PAIR_DATA"
	LARGE_TRADE_TXS_DATA   WsDataType = "LARGE_TRADE_TXS_DATA"
	WALLET_TXS_DATA        WsDataType = "WALLET_TXS_DATA"
)

type WsPriceSubData struct {
	// QueryType: "simple"
	QueryType string `json:"queryType"`
	// ChartType: "1m", "3m"...
	ChartType string `json:"chartType"`
	Address   string `json:"address"`
	// Currency: "usd"...
	Currency string `json:"currency"`
}

type WsPriceData struct {
	// Open price
	O float64 `json:"o"`
	// High price
	H float64 `json:"h"`
	// Low price
	L float64 `json:"l"`
	// Close price
	C float64 `json:"c"`
	// Event type (e.g. "ohlcv")
	EventType string `json:"eventType"`
	// Chart type (e.g. "1m")
	Type string `json:"type"`
	// UnixTime seconds
	UnixTime int64 `json:"unixTime"`
	// Volume
	V float64 `json:"v"`
	// Token symbol
	Symbol string `json:"symbol"`
	// Token address
	Address string `json:"address"`
}

type WsTxsSubData struct {
	// QueryType: "simple"
	QueryType string `json:"queryType"`
	// Token address
	Address string `json:"address"`
}

type WsTxTokenInfo struct {
	// Token symbol
	Symbol string `json:"symbol"`
	// Token decimals
	Decimals int `json:"decimals"`
	// Token address
	Address string `json:"address"`
	// Raw amount
	Amount int64 `json:"amount"`
	// Transaction type
	Type string `json:"type"`
	// Swap type (from/to)
	TypeSwap string `json:"typeSwap"`
	// UI formatted amount
	UiAmount float64 `json:"uiAmount"`
	// Token price
	Price float64 `json:"price,omitempty"`
	// Nearest price if price is not available
	NearestPrice float64 `json:"nearestPrice"`
	// Raw change amount
	ChangeAmount int64 `json:"changeAmount"`
	// UI formatted change amount
	UiChangeAmount float64 `json:"uiChangeAmount"`
	// Token icon URL
	Icon string `json:"icon,omitempty"`
}

type WsTxsData struct {
	// Block unix time
	BlockUnixTime int64 `json:"blockUnixTime"`
	// Owner address
	Owner string `json:"owner"`
	// Source platform
	Source string `json:"source"`
	// Transaction hash
	TxHash string `json:"txHash"`
	// Alias name
	Alias *string `json:"alias"`
	// Whether trade happened on Birdeye
	IsTradeOnBe bool `json:"isTradeOnBe"`
	// Platform address
	Platform string `json:"platform"`
	// Volume in USD
	VolumeUSD float64 `json:"volumeUSD"`
	// From token info
	From WsTxTokenInfo `json:"from"`
	// To token info
	To WsTxTokenInfo `json:"to"`
}

type WsBaseQuotePriceSubData struct {
	// Base token address
	BaseAddress string `json:"baseAddress"`
	// Quote token address
	QuoteAddress string `json:"quoteAddress"`
	// Chart type/interval
	ChartType string `json:"chartType"`
}

type WsBaseQuotePriceData struct {
	// Open price
	O float64 `json:"o"`
	// High price
	H float64 `json:"h"`
	// Low price
	L float64 `json:"l"`
	// Close price
	C float64 `json:"c"`
	// Event type
	EventType string `json:"eventType"`
	// Chart type/interval
	Type string `json:"type"`
	// Unix timestamp
	UnixTime int64 `json:"unixTime"`
	// Volume
	V float64 `json:"v"`
	// Base token address
	BaseAddress string `json:"baseAddress"`
	// Quote token address
	QuoteAddress string `json:"quoteAddress"`
}

type WsTokenNewListingSubData struct {
	// Whether meme platform is enabled
	MemePlatformEnabled bool `json:"meme_platform_enabled,omitempty"`
	// Minimum liquidity
	MinLiquidity float64 `json:"min_liquidity,omitempty"`
	// Maximum liquidity
	MaxLiquidity float64 `json:"max_liquidity,omitempty"`
}

type WsTokenNewListingData struct {
	// Token address
	Address string `json:"address"`
	// Token decimals
	Decimals int64 `json:"decimals"`
	// Token name
	Name string `json:"name"`
	// Token symbol
	Symbol string `json:"symbol"`
	// Token liquidity in USD
	Liquidity string `json:"liquidity"`
	// Unix timestamp when liquidity was added
	LiquidityAddedAt int64 `json:"liquidityAddedAt"`
}

type WsNewPairSubData struct {
	// Minimum liquidity
	MinLiquidity float64 `json:"min_liquidity,omitempty"`
	// Maximum liquidity
	MaxLiquidity float64 `json:"max_liquidity,omitempty"`
}

type WsNewPairTokenInfo struct {
	// Token address
	Address string `json:"address"`
	// Token name
	Name string `json:"name"`
	// Token symbol
	Symbol string `json:"symbol"`
	// Token decimals
	Decimals int64 `json:"decimals"`
}

type WsNewPairData struct {
	// Pair address
	Address string `json:"address"`
	// Pair name
	Name string `json:"name"`
	// Source of the pair
	Source string `json:"source"`
	// Base token info
	Base WsNewPairTokenInfo `json:"base"`
	// Quote token info
	Quote WsNewPairTokenInfo `json:"quote"`
	// Transaction hash
	TxHash string `json:"txHash"`
	// Block time
	BlockTime int64 `json:"blockTime"`
}

type WsLargeTradeTxsSubData struct {
	// Minimum volume in USD
	MinVolume float64 `json:"min_volume,omitempty"`
	// Maximum volume in USD
	MaxVolume float64 `json:"max_volume,omitempty"`
}

type WsLargeTradeTxsTokenInfo struct {
	// Token address
	Address string `json:"address"`
	// Token name
	Name string `json:"name"`
	// Token symbol
	Symbol string `json:"symbol"`
	// Token decimals
	Decimals int64 `json:"decimals"`
	// Token amount in UI format
	UiAmount float64 `json:"uiAmount"`
	// Token price
	Price *float64 `json:"price"`
	// Nearest token price
	NearestPrice float64 `json:"nearestPrice"`
	// Token amount change in UI format
	UiChangeAmount float64 `json:"uiChangeAmount"`
}

type WsLargeTradeTxsData struct {
	// Block unix timestamp seconds
	BlockUnixTime int64 `json:"blockUnixTime"`
	// Block time in human readable format
	BlockHumanTime string `json:"blockHumanTime"`
	// Owner address
	Owner string `json:"owner"`
	// Source of the trade
	Source string `json:"source"`
	// Pool address
	PoolAddress string `json:"poolAddress"`
	// Transaction hash
	TxHash string `json:"txHash"`
	// Volume in USD
	VolumeUSD float64 `json:"volumeUSD"`
	// Network name
	Network string `json:"network"`
	// From token info
	From WsLargeTradeTxsTokenInfo `json:"from"`
	// To token info
	To WsLargeTradeTxsTokenInfo `json:"to"`
}

type WsWalletTxsSubData struct {
	// Wallet address to monitor
	Address string `json:"address"`
}

type WsWalletTxsTokenInfo struct {
	// Token symbol
	Symbol string `json:"symbol"`
	// Token decimals
	Decimals int64 `json:"decimals"`
	// Token address
	Address string `json:"address"`
	// Token amount in UI format
	UiAmount float64 `json:"uiAmount"`
}

type WsWalletTxsData struct {
	// Transaction type
	Type string `json:"type"`
	// Block unix timestamp seconds
	BlockUnixTime int64 `json:"blockUnixTime"`
	// Block time in human readable format
	BlockHumanTime string `json:"blockHumanTime"`
	// Owner address
	Owner string `json:"owner"`
	// Source address
	Source string `json:"source"`
	// Transaction hash
	TxHash string `json:"txHash"`
	// Volume in USD
	VolumeUSD float64 `json:"volumeUSD"`
	// Network name
	Network string `json:"network"`
	// Base token info
	Base WsWalletTxsTokenInfo `json:"base"`
	// Quote token info
	Quote WsWalletTxsTokenInfo `json:"quote"`
}

type WsClient struct {
	url string
	ws  *websocket.Conn

	muReConn sync.RWMutex

	muSubers sync.RWMutex
	subers   map[WsSubType][]chan any

	logger *slog.Logger
}

func NewWsClient(chain, apiKey string, logger *slog.Logger) *WsClient {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	url := fmt.Sprintf("wss://public-api.birdeye.so/socket/%s?x-api-key=%s", chain, apiKey)
	return &WsClient{url: url, logger: logger}
}

func (c *WsClient) Connect() error {
	headers := http.Header{}
	headers.Add("Origin", "ws://public-api.birdeye.so")
	headers.Add("Sec-WebSocket-Origin", "ws://public-api.birdeye.so")
	headers.Add("Sec-WebSocket-Protocol", "echo-protocol")
	conn, reps, err := websocket.DefaultDialer.Dial(c.url, headers)
	if err != nil {
		return fmt.Errorf("eyebird: failed to connect to websocket: %w, http status code: %d", err, reps.StatusCode)
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
		c.logger.Info("eyebird: retrying to connect to websocket...")
		err := c.Connect()
		if err != nil {
			c.logger.Error("eyebird: failed to connect to websocket, retrying...", "error", err)
			time.Sleep(5 * time.Second)
		} else {
			c.logger.Info("eyebird: reconnected to websocket")
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
			c.logger.Error("eyebird: websocket read error", "error", err)
			c.reConn()
			continue
		}
		switch t {
		case websocket.BinaryMessage:
			c.logger.Info("eyebird: websocket binary message", "data", string(b))
		case websocket.PingMessage:
			c.logger.Info("eyebird: websocket ping message", "data", string(b))
		case websocket.PongMessage:
			c.logger.Info("eyebird: websocket pong message", "data", string(b))
		case websocket.TextMessage:
			go c.msgHandler(b)
		case websocket.CloseMessage:
			c.logger.Info("eyebird: websocket close message", "data", string(b))
		}
	}
}

func (c *WsClient) msgHandler(b []byte) {
	d := map[string]any{}
	err := json.Unmarshal(b, &d)
	if err != nil {
		c.logger.Error("eyebird: failed to unmarshal message", "error", err)
		return
	}
	t, ok := d["type"].(WsSubType)
	if !ok {
		c.logger.Error("eyebird: message type is not string", "data", string(b))
		return
	}
	b, err = json.Marshal(d["data"])
	if err != nil {
		c.logger.Error("eyebird: failed to marshal data", "error", err)
		return
	}
	var dd any
	switch t {
	case SUBSCRIBE_PRICE:
		dd = &WsPriceData{}
	case SUBSCRIBE_TXS:
		dd = &WsTxsData{}
	case SUBSCRIBE_BASE_QUOTE_PRICE:
		dd = &WsBaseQuotePriceData{}
	case SUBSCRIBE_TOKEN_NEW_LISTING:
		dd = &WsTokenNewListingData{}
	case SUBSCRIBE_NEW_PAIR:
		dd = &WsNewPairData{}
	case SUBSCRIBE_LARGE_TRADE_TXS:
		dd = &WsLargeTradeTxsData{}
	case SUBSCRIBE_WALLET_TXS:
		dd = &WsWalletTxsData{}
	}
	err = json.Unmarshal(b, dd)
	if err != nil {
		c.logger.Error("eyebird: failed to unmarshal data", "error", err)
		return
	}
	c.muSubers.RLock()
	defer c.muSubers.RUnlock()
	subers := c.subers[t]
	for _, suber := range subers {
		suber := suber
		go func() {
			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()
			select {
			case suber <- dd:
			case <-timer.C:
				c.logger.Error("eyebird: failed to send data to suber", "type", t)
			}
		}()
	}
}

func (c *WsClient) Subscribe(t WsSubType) <-chan any {
	c.muSubers.Lock()
	defer c.muSubers.Unlock()
	ch := make(chan any)
	c.subers[t] = append(c.subers[t], ch)
	return ch
}

func (c *WsClient) Unsubscribe(t WsSubType, ch <-chan any) {
	c.muSubers.Lock()
	defer c.muSubers.Unlock()
	for i, suber := range c.subers[t] {
		if suber != ch {
			continue
		}
		close(suber)
		if i == len(c.subers[t])-1 {
			c.subers[t] = c.subers[t][:i]
		} else {
			c.subers[t] = append(c.subers[t][:i], c.subers[t][i+1:]...)
		}
		break
	}
}
