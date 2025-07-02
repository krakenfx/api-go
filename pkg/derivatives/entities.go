package derivatives

import (
	"encoding/json"
	"time"

	"github.com/krakenfx/api-go/v2/pkg/decimal"
)

type MarginSchedule struct {
	Contracts           int              `json:"contracts,omitempty"`
	NumNonContractUnits *decimal.Decimal `json:"numNonContractUnits,omitempty"`
	InitialMargin       *decimal.Decimal `json:"initialMargin,omitempty"`
	MaintenanceMargin   *decimal.Decimal `json:"maintenanceMargin,omitempty"`
}

type Instrument struct {
	Category                    string                    `json:"category,omitempty"`
	ContractSize                *decimal.Decimal          `json:"contractSize,omitempty"`
	ContractValueTradePrecision *decimal.Decimal          `json:"contractValueTradePrecision,omitempty"`
	FundingRateCoefficient      *decimal.Decimal          `json:"fundingRateCoefficient,omitempty"`
	ImpactMidSize               *decimal.Decimal          `json:"impactMidSize,omitempty"`
	ISIN                        string                    `json:"isin,omitempty"`
	LastTradingTime             time.Time                 `json:"lastTradingTime,omitempty"`
	MarginSchedules             map[string]MarginSchedule `json:"marginSchedules,omitempty"`
	RetailMarginLevels          []MarginSchedule          `json:"retailMarginLevels,omitempty"`
	MarginLevels                []MarginSchedule          `json:"marginLevels,omitempty"`
	MaxPositionSize             *decimal.Decimal          `json:"maxPositionSize,omitempty"`
	MaxRelativeFundingRate      *decimal.Decimal          `json:"maxRelativeFundingRate,omitempty"`
	OpeningDate                 time.Time                 `json:"openingDate,omitempty"`
	PostOnly                    bool                      `json:"postOnly,omitempty"`
	FeeScheduleUid              string                    `json:"feeScheduleUid,omitempty"`
	Symbol                      string                    `json:"symbol,omitempty"`
	Pair                        string                    `json:"pair,omitempty"`
	Base                        string                    `json:"base,omitempty"`
	Quote                       string                    `json:"quote,omitempty"`
	Tags                        []string                  `json:"tags,omitempty"`
	TickSize                    *decimal.Decimal          `json:"tickSize,omitempty"`
	Tradeable                   bool                      `json:"tradeable,omitempty"`
	Type                        string                    `json:"type,omitempty"`
	Underlying                  string                    `json:"underlying,omitempty"`
	UnderlyingFuture            string                    `json:"underlyingFuture,omitempty"`
	TradFi                      bool                      `json:"tradfi,omitempty"`
	Mtf                         bool                      `json:"mtf,omitempty"`
}

type Greeks struct {
	IV *decimal.Decimal `json:"iv,omitempty"`
}

type TickerData struct {
	Symbol                string           `json:"symbol,omitempty"`
	Last                  *decimal.Decimal `json:"last,omitempty"`
	LastTime              time.Time        `json:"lastTime,omitempty"`
	LastSize              *decimal.Decimal `json:"lastSize,omitempty"`
	Tag                   string           `json:"tag,omitempty"`
	Pair                  string           `json:"pair,omitempty"`
	MarkPrice             *decimal.Decimal `json:"markPrice,omitempty"`
	Bid                   *decimal.Decimal `json:"bid,omitempty"`
	BidSize               *decimal.Decimal `json:"bidSize,omitempty"`
	Ask                   *decimal.Decimal `json:"ask,omitempty"`
	AskSize               *decimal.Decimal `json:"askSize,omitempty"`
	Vol24h                *decimal.Decimal `json:"vol24h,omitempty"`
	VolumeQuote           *decimal.Decimal `json:"volumeQuote,omitempty"`
	OpenInterest          *decimal.Decimal `json:"openInterest,omitempty"`
	Open24h               *decimal.Decimal `json:"open24h,omitempty"`
	High24h               *decimal.Decimal `json:"high24h,omitempty"`
	Low24h                *decimal.Decimal `json:"low24h,omitempty"`
	ExtrinsicValue        *decimal.Decimal `json:"extrinsicValue,omitempty"`
	FundingRate           *decimal.Decimal `json:"fundingRate,omitempty"`
	FundingRatePrediction *decimal.Decimal `json:"fundingRatePrediction,omitempty"`
	Suspended             bool             `json:"suspended,omitempty"`
	IndexPrice            *decimal.Decimal `json:"indexPrice,omitempty"`
	PostOnly              bool             `json:"postOnly,omitempty"`
	Change24h             *decimal.Decimal `json:"change24h,omitempty"`
}

type JSONOrderBook struct {
	Asks [][]*decimal.Decimal `json:"asks,omitempty"`
	Bids [][]*decimal.Decimal `json:"bids,omitempty"`
}

func (job JSONOrderBook) OrderBook() (book OrderBook) {
	book.Asks = make([]PriceLevel, len(job.Asks))
	for i, ask := range job.Asks {
		book.Asks[i] = PriceLevel{
			Price:  ask[0],
			Volume: ask[1],
		}
	}
	book.Bids = make([]PriceLevel, len(job.Bids))
	for i, bid := range job.Bids {
		book.Bids[i] = PriceLevel{
			Price:  bid[0],
			Volume: bid[1],
		}
	}
	return
}

type OrderBook struct {
	Asks []PriceLevel `json:"asks,omitempty"`
	Bids []PriceLevel `json:"bids,omitempty"`
}

func (ob *OrderBook) UnmarshalJSON(data []byte) error {
	var v JSONOrderBook
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*ob = v.OrderBook()
	return nil
}

type PriceLevel struct {
	Price  *decimal.Decimal `json:"price,omitempty"`
	Volume *decimal.Decimal `json:"volume,omitempty"`
}

type Trade struct {
	Price                         *decimal.Decimal `json:"price,omitempty"`
	Side                          string           `json:"side,omitempty"`
	Size                          *decimal.Decimal `json:"size,omitempty"`
	Time                          time.Time        `json:"time,omitempty"`
	TradeID                       int              `json:"trade_id,omitempty"`
	Type                          string           `json:"type,omitempty"`
	UID                           string           `json:"uid,omitempty"`
	InstrumentIdentificationType  string           `json:"instrument_identification_type,omitempty"`
	ISIN                          string           `json:"isin,omitempty"`
	ExecutionVenue                string           `json:"execution_venue,omitempty"`
	PriceNotation                 string           `json:"price_notation,omitempty"`
	PriceCurrency                 string           `json:"price_currency,omitempty"`
	NotionalAmount                *decimal.Decimal `json:"notional_amount,omitempty"`
	NotionalCurrency              string           `json:"notional_currency,omitempty"`
	PublicationTime               string           `json:"publication_time,omitempty"`
	PublicationVenue              string           `json:"publication_venue,omitempty"`
	TransactionIdentificationCode string           `json:"transaction_identification_code,omitempty"`
	ToBeCleared                   bool             `json:"to_be_cleared,omitempty"`
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
	OrderID        string           `json:"order_id,omitempty"`
	ClientOrderID  string           `json:"cliOrdId,omitempty"`
	Status         string           `json:"status,omitempty"`
	Side           string           `json:"side,omitempty"`
	OrderType      string           `json:"orderType,omitempty"`
	Symbol         string           `json:"symbol,omitempty"`
	LimitPrice     *decimal.Decimal `json:"limitPrice,omitempty"`
	StopPrice      *decimal.Decimal `json:"stopPrice,omitempty"`
	FilledSize     *decimal.Decimal `json:"filledSize,omitempty"`
	UnfilledSize   *decimal.Decimal `json:"unfilledSize,omitempty"`
	ReduceOnly     bool             `json:"reduceOnly,omitempty"`
	TriggerSignal  string           `json:"triggerSignal,omitempty"`
	LastUpdateTime time.Time        `json:"lastUpdateTime,omitempty"`
	ReceivedTime   time.Time        `json:"receivedTime,omitempty"`
}

type DerivativesResponse struct {
	Result     string    `json:"result,omitempty"`
	ServerTime time.Time `json:"serverTime,omitempty"`
}
