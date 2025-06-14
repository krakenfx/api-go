package spot

import (
	"github.com/krakenfx/api-go/pkg/kraken"
)

type FullName struct {
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
}

type TaxID struct {
	ID             string `json:"id,omitempty"`
	IssuingCountry string `json:"issuing_country,omitempty"`
}

type Residence struct {
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Province   string `json:"province,omitempty"`
	Country    string `json:"country,omitempty"`
}

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

type ActionDetailsType struct {
	Type    string `json:"type,omitempty"`
	Version int    `json:"version,omitempty"`
}

type UserRequiredAction struct {
	ActionType       string             `json:"action_type,omitempty"`
	VerificationType string             `json:"verification_type,omitempty"`
	WaitReasonCode   string             `json:"wait_reason_code,omitempty"`
	DetailsType      *ActionDetailsType `json:"details_type,omitempty"`
	Reasons          []string           `json:"reasons,omitempty"`
	Deadline         *string            `json:"deadline,omitempty"`
}

type UserStatus struct {
	State           string                `json:"state,omitempty"`
	Reasons         []string              `json:"reasons,omitempty"`
	RequiredActions []*UserRequiredAction `json:"required_actions,omitempty"`
}

type Identity struct {
	FullName    *FullName `json:"full_name,omitempty"`
	DateOfBirth string    `json:"date_of_birth,omitempty"`
}

type Watchlist struct {
	Status                 string `json:"status,omitempty"`
	Verifier               string `json:"verifier,omitempty"`
	VerifiedAt             string `json:"verified_at,omitempty"`
	VerifierResponse       any    `json:"verifier_response,omitempty"`
	ExternalVerificationID string `json:"external_verification_id,omitempty"`
	ExpirationDate         string `json:"expiration_date,omitempty"`
}

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

type DepositMethod struct {
	Method        string        `json:"method,omitempty"`
	Limit         *kraken.Money `json:"limit,omitempty"`
	Fee           *kraken.Money `json:"fee,omitempty"`
	FeePercentage *kraken.Money `json:"fee-percentage,omitempty"`
	GenAddress    bool          `json:"gen-address,omitempty"`
	Minimum       *kraken.Money `json:"minimum,omitempty"`
}

type DepositAddress struct {
	Address  string `json:"address,omitempty"`
	ExpireTM string `json:"expiretm,omitempty"`
	New      bool   `json:"new,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

type DepositStatus struct {
	Method      string   `json:"method,omitempty"`
	Aclass      string   `json:"aclass,omitempty"`
	Asset       string   `json:"asset,omitempty"`
	Refid       string   `json:"refid,omitempty"`
	Txid        string   `json:"txid,omitempty"`
	Info        string   `json:"info,omitempty"`
	Amount      string   `json:"amount,omitempty"`
	Fee         string   `json:"fee,omitempty"`
	Time        int64    `json:"time,omitempty"`
	Status      string   `json:"status,omitempty"`
	StatusProp  string   `json:"status-prop,omitempty"`
	Originators []string `json:"originators,omitempty"`
}

type WithdrawMethod struct {
	Asset   string        `json:"asset,omitempty"`
	Method  string        `json:"method,omitempty"`
	Network string        `json:"network,omitempty"`
	Minimum *kraken.Money `json:"minimum,omitempty"`
}

type WithdrawAddress struct {
	Address  string `json:"address,omitempty"`
	Asset    string `json:"asset,omitempty"`
	Method   string `json:"method,omitempty"`
	Key      string `json:"key,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Verified bool   `json:"verified,omitempty"`
}

type WithdrawInfo struct {
	Method string        `json:"method,omitempty"`
	Limit  *kraken.Money `json:"limit,omitempty"`
	Amount *kraken.Money `json:"amount,omitempty"`
	Fee    *kraken.Money `json:"fee,omitempty"`
}

type WithdrawStatus struct {
	Method     string `json:"method,omitempty"`
	Network    string `json:"network,omitempty"`
	Aclass     string `json:"aclass,omitempty"`
	Asset      string `json:"asset,omitempty"`
	Refid      string `json:"refid,omitempty"`
	Txid       string `json:"txid,omitempty"`
	Info       string `json:"info,omitempty"`
	Amount     string `json:"amount,omitempty"`
	Fee        string `json:"fee,omitempty"`
	Time       int64  `json:"time,omitempty"`
	Status     string `json:"status,omitempty"`
	StatusProp string `json:"status-prop,omitempty"`
	Key        string `json:"key,omitempty"`
}
type EarnStrategyAPR struct {
	Low  string `json:"low,omitempty"`
	High string `json:"high,omitempty"`
}
type EarnStrategy struct {
	ID                        string           `json:"id,omitempty"`
	Asset                     string           `json:"asset,omitempty"`
	LockType                  any              `json:"lock_type,omitempty"`
	AprEstimate               *EarnStrategyAPR `json:"apr_estimate,omitempty"`
	UserMinAllocation         string           `json:"user_min_allocation,omitempty"`
	AllocationFee             any              `json:"allocation_fee,omitempty"`
	DeallocationFee           any              `json:"deallocation_fee,omitempty"`
	AutoCompound              any              `json:"auto_compound,omitempty"`
	YieldSource               any              `json:"yield_source,omitempty"`
	CanAllocate               bool             `json:"can_allocate,omitempty"`
	CanDeallocate             bool             `json:"can_deallocate,omitempty"`
	AllocationRestrictionInfo []string         `json:"allocation_restriction_info,omitempty"`
}

type EarnAllocationReward struct {
	Native    string `json:"native"`
	Converted string `json:"converted"`
}
type EarnAllocationAmountState struct {
	Native          string                       `json:"native,omitempty"`
	Converted       string                       `json:"converted,omitempty"`
	AllocationCount int                          `json:"allocation_count,omitempty"`
	Allocations     []*EarnAllocationStateDetail `json:"allocations,omitempty"`
}
type EarnAllocationStateDetail struct {
	Native    string `json:"native,omitempty"`
	Converted string `json:"converted,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Expires   string `json:"expires,omitempty"`
}
type EarnAllocationAmount struct {
	Bonding   *EarnAllocationAmountState `json:"bonding,omitempty"`
	ExitQueue *EarnAllocationAmountState `json:"exit_queue,omitempty"`
	Pending   *EarnAllocationAmountState `json:"pending,omitempty"`
	Total     *EarnAllocationAmountState `json:"total,omitempty"`
	Unbonding *EarnAllocationAmountState `json:"unbonding,omitempty"`
}
type EarnAllocationPayout struct {
	AccumulatedReward *EarnAllocationReward `json:"accumulated_reward,omitempty"`
	EstimatedReward   *EarnAllocationReward `json:"estimated_reward,omitempty"`
	PeriodStart       string                `json:"period_start,omitempty"`
	PeriodEnd         string                `json:"period_end,omitempty"`
}
type EarnAllocation struct {
	StrategyID      string                `json:"strategy_id,omitempty"`
	NativeAsset     string                `json:"native_asset,omitempty"`
	AmountAllocated EarnAllocationAmount  `json:"amount_allocated,omitempty"`
	TotalRewarded   EarnAllocationReward  `json:"total_rewarded,omitempty"`
	Payout          *EarnAllocationPayout `json:"payout,omitempty"`
}
