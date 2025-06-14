package derivatives

import (
	"time"

	"github.com/krakenfx/api-go/pkg/kraken"
)

type MarginSchedule struct {
	Contracts           int           `json:"contracts,omitempty"`
	NumNonContractUnits *kraken.Money `json:"numNonContractUnits,omitempty"`
	InitialMargin       *kraken.Money `json:"initialMargin,omitempty"`
	MaintenanceMargin   *kraken.Money `json:"maintenanceMargin,omitempty"`
}

type Instrument struct {
	Category                    string                    `json:"category,omitempty"`
	ContractSize                *kraken.Money             `json:"contractSize,omitempty"`
	ContractValueTradePrecision *kraken.Money             `json:"contractValueTradePrecision,omitempty"`
	FundingRateCoefficient      *kraken.Money             `json:"fundingRateCoefficient,omitempty"`
	ImpactMidSize               *kraken.Money             `json:"impactMidSize,omitempty"`
	ISIN                        string                    `json:"isin,omitempty"`
	LastTradingTime             time.Time                 `json:"lastTradingTime,omitempty"`
	MarginSchedules             map[string]MarginSchedule `json:"marginSchedules,omitempty"`
	RetailMarginLevels          []MarginSchedule          `json:"retailMarginLevels,omitempty"`
	MarginLevels                []MarginSchedule          `json:"marginLevels,omitempty"`
	MaxPositionSize             *kraken.Money             `json:"maxPositionSize,omitempty"`
	MaxRelativeFundingRate      *kraken.Money             `json:"maxRelativeFundingRate,omitempty"`
	OpeningDate                 time.Time                 `json:"openingDate,omitempty"`
	PostOnly                    bool                      `json:"postOnly,omitempty"`
	FeeScheduleUid              string                    `json:"feeScheduleUid,omitempty"`
	Symbol                      string                    `json:"symbol,omitempty"`
	Pair                        string                    `json:"pair,omitempty"`
	Base                        string                    `json:"base,omitempty"`
	Quote                       string                    `json:"quote,omitempty"`
	Tags                        []string                  `json:"tags,omitempty"`
	TickSize                    *kraken.Money             `json:"tickSize,omitempty"`
	Tradeable                   bool                      `json:"tradeable,omitempty"`
	Type                        string                    `json:"type,omitempty"`
	Underlying                  string                    `json:"underlying,omitempty"`
	UnderlyingFuture            string                    `json:"underlyingFuture,omitempty"`
	TradFi                      bool                      `json:"tradfi,omitempty"`
	Mtf                         bool                      `json:"mtf,omitempty"`
}

type Greeks struct {
	IV *kraken.Money `json:"iv,omitempty"`
}

type TickerData struct {
	Symbol                string        `json:"symbol,omitempty"`
	Last                  *kraken.Money `json:"last,omitempty"`
	LastTime              time.Time     `json:"lastTime,omitempty"`
	LastSize              *kraken.Money `json:"lastSize,omitempty"`
	Tag                   string        `json:"tag,omitempty"`
	Pair                  string        `json:"pair,omitempty"`
	MarkPrice             *kraken.Money `json:"markPrice,omitempty"`
	Bid                   *kraken.Money `json:"bid,omitempty"`
	BidSize               *kraken.Money `json:"bidSize,omitempty"`
	Ask                   *kraken.Money `json:"ask,omitempty"`
	AskSize               *kraken.Money `json:"askSize,omitempty"`
	Vol24h                *kraken.Money `json:"vol24h,omitempty"`
	VolumeQuote           *kraken.Money `json:"volumeQuote,omitempty"`
	OpenInterest          *kraken.Money `json:"openInterest,omitempty"`
	Open24h               *kraken.Money `json:"open24h,omitempty"`
	High24h               *kraken.Money `json:"high24h,omitempty"`
	Low24h                *kraken.Money `json:"low24h,omitempty"`
	ExtrinsicValue        *kraken.Money `json:"extrinsicValue,omitempty"`
	FundingRate           *kraken.Money `json:"fundingRate,omitempty"`
	FundingRatePrediction *kraken.Money `json:"fundingRatePrediction,omitempty"`
	Suspended             bool          `json:"suspended,omitempty"`
	IndexPrice            *kraken.Money `json:"indexPrice,omitempty"`
	PostOnly              bool          `json:"postOnly,omitempty"`
	Change24h             *kraken.Money `json:"change24h,omitempty"`
}

type OrderBook struct {
	Asks [][]*kraken.Money `json:"asks,omitempty"`
	Bids [][]*kraken.Money `json:"bids,omitempty"`
}

type Trade struct {
	Price                         *kraken.Money `json:"price,omitempty"`
	Side                          string        `json:"side,omitempty"`
	Size                          *kraken.Money `json:"size,omitempty"`
	Time                          time.Time     `json:"time,omitempty"`
	TradeID                       int           `json:"trade_id,omitempty"`
	Type                          string        `json:"type,omitempty"`
	UID                           string        `json:"uid,omitempty"`
	InstrumentIdentificationType  string        `json:"instrument_identification_type,omitempty"`
	ISIN                          string        `json:"isin,omitempty"`
	ExecutionVenue                string        `json:"execution_venue,omitempty"`
	PriceNotation                 string        `json:"price_notation,omitempty"`
	PriceCurrency                 string        `json:"price_currency,omitempty"`
	NotionalAmount                *kraken.Money `json:"notional_amount,omitempty"`
	NotionalCurrency              string        `json:"notional_currency,omitempty"`
	PublicationTime               string        `json:"publication_time,omitempty"`
	PublicationVenue              string        `json:"publication_venue,omitempty"`
	TransactionIdentificationCode string        `json:"transaction_identification_code,omitempty"`
	ToBeCleared                   bool          `json:"to_be_cleared,omitempty"`
}

type OrderStatus struct {
	ClientOrderID string           `json:"cliOrdId,omitempty"`
	OrderEvents   []map[string]any `json:"orderEvents,omitempty"`
	OrderID       string           `json:"order_id,omitempty"`
	ReceivedTime  time.Time        `json:"receivedTime,omitempty"`
	Status        string           `json:"status,omitempty"`
}

type BatchOrderInstruction struct {
	Order                     string `json:"order,omitempty"`
	OrderTag                  string `json:"order_tag,omitempty"`
	OrderID                   string `json:"order_id,omitempty"`
	OrderType                 string `json:"orderType,omitempty"`
	Symbol                    string `json:"symbol,omitempty"`
	Side                      string `json:"side,omitempty"`
	Size                      string `json:"size,omitempty"`
	LimitPrice                string `json:"limitPrice,omitempty"`
	StopPrice                 string `json:"stopPrice,omitempty"`
	ClientOrderID             string `json:"cliOrdId,omitempty"`
	TriggerSignal             string `json:"triggerSignal,omitempty"`
	ReduceOnly                bool   `json:"reduceOnly,omitempty"`
	TrailingStopMaxDeviation  string `json:"trailingStopMaxDeviation,omitempty"`
	TrailingStopDeviationUnit string `json:"trailingStopDeviationUnit,omitempty"`
}

type BatchOrderJson struct {
	BatchOrder []*BatchOrderInstruction `json:"batchOrder,omitempty"`
}

type BatchStatusInfo struct {
	ClientOrderID    string           `json:"cliOrdId,omitempty"`
	DateTimeReceived time.Time        `json:"dateTimeReceived,omitempty"`
	OrderEvents      []map[string]any `json:"orderEvents,omitempty"`
	OrderID          string           `json:"order_id,omitempty"`
	OrderTag         string           `json:"order_tag,omitempty"`
	Status           string           `json:"status,omitempty"`
}

type CancelledOrder struct {
	ClientOrderID string `json:"cliOrdId,omitempty"`
	OrderID       string `json:"order_id,omitempty"`
}

type CancelStatus struct {
	CancelOnly      string           `json:"cancelOnly,omitempty"`
	CancelledOrders []CancelledOrder `json:"cancelledOrders,omitempty"`
	OrderEvents     []map[string]any `json:"orderEvents,omitempty"`
	ReceivedTime    time.Time        `json:"receivedTime,omitempty"`
	Status          string           `json:"status,omitempty"`
}

type OpenOrder struct {
	OrderID        string        `json:"order_id,omitempty"`
	ClientOrderID  string        `json:"cliOrdId,omitempty"`
	Status         string        `json:"status,omitempty"`
	Side           string        `json:"side,omitempty"`
	OrderType      string        `json:"orderType,omitempty"`
	Symbol         string        `json:"symbol,omitempty"`
	LimitPrice     *kraken.Money `json:"limitPrice,omitempty"`
	StopPrice      *kraken.Money `json:"stopPrice,omitempty"`
	FilledSize     *kraken.Money `json:"filledSize,omitempty"`
	UnfilledSize   *kraken.Money `json:"unfilledSize,omitempty"`
	ReduceOnly     bool          `json:"reduceOnly,omitempty"`
	TriggerSignal  string        `json:"triggerSignal,omitempty"`
	LastUpdateTime time.Time     `json:"lastUpdateTime,omitempty"`
	ReceivedTime   time.Time     `json:"receivedTime,omitempty"`
}

type DerivativesResponse struct {
	Result     string    `json:"result,omitempty"`
	ServerTime time.Time `json:"serverTime,omitempty"`
}
