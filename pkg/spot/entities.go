package spot

import (
	"github.com/krakenfx/api-go/pkg/kraken"
)

// FullName contains the first, middle, and last name of the client.
type FullName struct {
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
}

// TaxID contains the tax identification string and issuing country.
type TaxID struct {
	ID             string `json:"id,omitempty"`
	IssuingCountry string `json:"issuing_country,omitempty"`
}

// Residence contains the address of the client.
type Residence struct {
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Province   string `json:"province,omitempty"`
	Country    string `json:"country,omitempty"`
}

// UserInfo contains the profile of the client.
type UserInfo struct {
	Email              string     `json:"email,omitempty"`
	ExternalID         string     `json:"external_id,omitempty"`
	TOSVersionAccepted int        `json:"tos_version_accepted,omitempty"`
	FullName           *FullName  `json:"full_name,omitempty"`
	DateOfBirth        string     `json:"date_of_birth,omitempty"`
	Residence          *Residence `json:"residence,omitempty"`
	Phone              string     `json:"phone,omitempty"`
	Nationalities      []string   `json:"nationalities,omitempty"`
	Occupation         string     `json:"occupation,omitempty"`
	CityOfBirth        string     `json:"city_of_birth,omitempty"`
	CountryOfBirth     string     `json:"country_of_birth,omitempty"`
	TaxIDs             []*TaxID   `json:"tax_ids,omitempty"`
	Language           string     `json:"language,omitempty"`
}

// ActionDetailsType contains the type and version of the [UserRequiredAction] object.
type ActionDetailsType struct {
	Type    string `json:"type,omitempty"`
	Version int    `json:"version,omitempty"`
}

// UserRequiredAction contains information regarding a required action.
type UserRequiredAction struct {
	ActionType       string             `json:"action_type,omitempty"`
	VerificationType string             `json:"verification_type,omitempty"`
	WaitReasonCode   string             `json:"wait_reason_code,omitempty"`
	DetailsType      *ActionDetailsType `json:"details_type,omitempty"`
	Reasons          []string           `json:"reasons,omitempty"`
	Deadline         *string            `json:"deadline,omitempty"`
}

// UserStatus contains the current status of the client.
type UserStatus struct {
	State           string                `json:"state,omitempty"`
	Reasons         []string              `json:"reasons,omitempty"`
	RequiredActions []*UserRequiredAction `json:"required_actions,omitempty"`
}

// Identity contains the full name and date of birth of the client.
type Identity struct {
	FullName    *FullName `json:"full_name,omitempty"`
	DateOfBirth string    `json:"date_of_birth,omitempty"`
}

// Watchlist contains the details of a watchlist.
type Watchlist struct {
	Status                 string `json:"status,omitempty"`
	Verifier               string `json:"verifier,omitempty"`
	VerifiedAt             string `json:"verified_at,omitempty"`
	VerifierResponse       any    `json:"verifier_response,omitempty"`
	ExternalVerificationID string `json:"external_verification_id,omitempty"`
	ExpirationDate         string `json:"expiration_date,omitempty"`
}

// VerificationMetadata contains metadata regarding a verification.
type VerificationMetadata struct {
	Identity               *Identity  `json:"identity,omitempty"`
	Address                *Residence `json:"address,omitempty"`
	Sanctions              *Watchlist `json:"sanctions,omitempty"`
	NegativeNews           *Watchlist `json:"negative_news,omitempty"`
	Pep                    *Watchlist `json:"pep,omitempty"`
	SelfieType             string     `json:"selfie_type,omitempty"`
	DocumentType           string     `json:"document_type,omitempty"`
	DocumentNumber         string     `json:"document_number,omitempty"`
	IssuingCountry         string     `json:"issuing_country,omitempty"`
	Nationality            string     `json:"nationality,omitempty"`
	Verifier               string     `json:"verifier,omitempty"`
	VerifiedAt             string     `json:"verified_at,omitempty"`
	VerifierResponse       any        `json:"verifier_response,omitempty"`
	ExternalVerificationID string     `json:"external_verification_id,omitempty"`
	ExpirationDate         string     `json:"expiration_date,omitempty"`
}

type Trade struct {
	OrderID        string        `json:"ordertxid,omitempty"`
	PositionID     string        `json:"postxid,omitempty"`
	Pair           string        `json:"pair,omitempty"`
	Time           *kraken.Money `json:"time,omitempty"`
	Type           string        `json:"type,omitempty"`
	OrderType      string        `json:"ordertype,omitempty"`
	Price          *kraken.Money `json:"price,omitempty"`
	Cost           *kraken.Money `json:"cost,omitempty"`
	Fee            *kraken.Money `json:"fee,omitempty"`
	Volume         *kraken.Money `json:"vol,omitempty"`
	Margin         *kraken.Money `json:"margin,omitempty"`
	Leverage       *kraken.Money `json:"leverage,omitempty"`
	Misc           string        `json:"misc,omitempty"`
	Ledgers        []string      `json:"ledgers,omitempty"`
	TradeID        *kraken.Money `json:"trade_id,omitempty"`
	Maker          bool          `json:"maker,omitempty"`
	PositionStatus string        `json:"posstatus,omitempty"`
	CPrice         *kraken.Money `json:"cprice,omitempty"`
	CCost          *kraken.Money `json:"ccost,omitempty"`
	CFee           *kraken.Money `json:"cfee,omitempty"`
	CVol           *kraken.Money `json:"cvol,omitempty"`
	CMargin        *kraken.Money `json:"cmargin,omitempty"`
	Net            *kraken.Money `json:"net,omitempty"`
	Trades         []string      `json:"trades,omitempty"`
}

type OrderDescriptionInfo struct {
	Order string `json:"order,omitempty"`
}

type OrderDescriptionInfoWithClose struct {
	Close string `json:"close,omitempty"`
	OrderDescriptionInfo
}

type OrderDescription struct {
	Pair           string        `json:"pair,omitempty"`
	Type           string        `json:"type,omitempty"`
	OrderType      string        `json:"ordertype,omitempty"`
	Price          *kraken.Money `json:"price,omitempty"`
	SecondaryPrice *kraken.Money `json:"price2,omitempty"`
	Leverage       string        `json:"leverage,omitempty"`
	OrderDescriptionInfoWithClose
}

type Order struct {
	RefID          string            `json:"refid,omitempty"`
	UserRef        *kraken.Money     `json:"userref,omitempty"`
	ClOrdID        string            `json:"cl_ord_id,omitempty"`
	Status         string            `json:"status,omitempty"`
	OpenTm         *kraken.Money     `json:"opentm,omitempty"`
	StartTm        *kraken.Money     `json:"starttm,omitempty"`
	ExpireTm       *kraken.Money     `json:"expiretm,omitempty"`
	Description    *OrderDescription `json:"descr,omitempty"`
	Volume         *kraken.Money     `json:"vol,omitempty"`
	VolumeExecuted *kraken.Money     `json:"vol_exec,omitempty"`
	Cost           *kraken.Money     `json:"cost,omitempty"`
	Fee            *kraken.Money     `json:"fee,omitempty"`
	Price          *kraken.Money     `json:"price,omitempty"`
	StopPrice      *kraken.Money     `json:"stopprice,omitempty"`
	LimitPrice     *kraken.Money     `json:"limitprice,omitempty"`
	Trigger        string            `json:"trigger,omitempty"`
	Margin         bool              `json:"margin,omitempty"`
	Misc           string            `json:"misc,omitempty"`
	SenderSubID    string            `json:"sender_sub_id,omitempty"`
	OrderFlags     string            `json:"oflags,omitempty"`
	Trades         []string          `json:"trades,omitempty"`
}

type ClosedOrder struct {
	CloseTm *kraken.Money `json:"closetm,omitempty"`
	Reason  string        `json:"reason,omitempty"`
	*Order
}

type OrderRequest struct {
	UserRef        int    `json:"userref,omitempty"`
	ClOrdId        string `json:"cl_ord_id,omitempty"`
	OrderType      string `json:"ordertype,omitempty"`
	Type           string `json:"type,omitempty"`
	Volume         string `json:"volume,omitempty"`
	DisplayVol     string `json:"displayvol,omitempty"`
	Price          string `json:"price,omitempty"`
	SecondaryPrice string `json:"price2,omitempty"`
	Trigger        string `json:"trigger,omitempty"`
	Leverage       string `json:"leverage,omitempty"`
	ReduceOnly     bool   `json:"reduce_only,omitempty"`
	StpType        string `json:"stptype,omitempty"`
	OrderFlags     string `json:"oflags,omitempty"`
	TimeInForce    string `json:"timeinforce,omitempty"`
	StartTm        string `json:"starttm,omitempty"`
	ExpireTm       string `json:"expiretm,omitempty"`
}

type OrderPlacementSingle struct {
	Descr OrderDescriptionInfoWithClose `json:"descr,omitempty"`
	ID    []string                      `json:"txid,omitempty"`
}

type OrderPlacementBatch struct {
	Descr OrderDescriptionInfo `json:"descr,omitempty"`
	Error string               `json:"error,omitempty"`
	ID    string               `json:"txid,omitempty"`
}

type AssetInfo struct {
	AssetClass      string        `json:"aclass,omitempty"`
	AltName         string        `json:"altname,omitempty"`
	Decimals        int           `json:"decimals,omitempty"`
	DisplayDecimals int           `json:"display_decimals,omitempty"`
	CollateralValue *kraken.Money `json:"collateral_value,omitempty"`
	Status          string        `json:"status,omitempty"`
}

type AssetPair struct {
	AltName            string            `json:"altname,omitempty"`
	WSName             string            `json:"wsname,omitempty"`
	BaseAssetClass     string            `json:"aclass_base,omitempty"`
	Base               string            `json:"base,omitempty"`
	QuoteAssetClass    string            `json:"aclass_quote,omitempty"`
	Quote              string            `json:"quote,omitempty"`
	PairDecimals       int               `json:"pair_decimals,omitempty"`
	CostDecimals       int               `json:"cost_decimals,omitempty"`
	LotDecimals        int               `json:"lot_decimals,omitempty"`
	LotMultiplier      int               `json:"lot_multiplier,omitempty"`
	BuyLeverage        []int             `json:"leverage_buy,omitempty"`
	SellLeverage       []int             `json:"leverage_sell,omitempty"`
	Fees               [][]*kraken.Money `json:"fees,omitempty"`
	FeesMaker          [][]*kraken.Money `json:"fees_maker,omitempty"`
	FeeVolumeCurrency  string            `json:"fee_volume_currency,omitempty"`
	MarginCall         int               `json:"margin_call,omitempty"`
	MarginStop         int               `json:"margin_stop,omitempty"`
	OrderMinimum       *kraken.Money     `json:"ordermin,omitempty"`
	CostMinimum        *kraken.Money     `json:"costmin,omitempty"`
	TickSize           *kraken.Money     `json:"tick_size,omitempty"`
	Status             string            `json:"status,omitempty"`
	LongPositionLimit  int               `json:"long_position_limit,omitempty"`
	ShortPositionLimit int               `json:"short_position_limit,omitempty"`
}

type AssetTickerInfo struct {
	Ask    []*kraken.Money `json:"a,omitempty"`
	Bid    []*kraken.Money `json:"b,omitempty"`
	Close  []*kraken.Money `json:"c,omitempty"`
	Volume []*kraken.Money `json:"v,omitempty"`
	VWAP   []*kraken.Money `json:"p,omitempty"`
	Trades []int           `json:"t,omitempty"`
	Low    []*kraken.Money `json:"l,omitempty"`
	High   []*kraken.Money `json:"h,omitempty"`
	Open   *kraken.Money   `json:"o,omitempty"`
}

type OrderBook struct {
	Asks []*kraken.Money `json:"asks,omitempty"`
	Bids []*kraken.Money `json:"bids,omitempty"`
}
