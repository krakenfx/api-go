package spot

import "github.com/krakenfx/api-go/v2/pkg/decimal"

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
	OrderID        string           `json:"ordertxid,omitempty"`
	PositionID     string           `json:"postxid,omitempty"`
	Pair           string           `json:"pair,omitempty"`
	Time           *decimal.Decimal `json:"time,omitempty"`
	Type           string           `json:"type,omitempty"`
	OrderType      string           `json:"ordertype,omitempty"`
	Price          *decimal.Decimal `json:"price,omitempty"`
	Cost           *decimal.Decimal `json:"cost,omitempty"`
	Fee            *decimal.Decimal `json:"fee,omitempty"`
	Volume         *decimal.Decimal `json:"vol,omitempty"`
	Margin         *decimal.Decimal `json:"margin,omitempty"`
	Leverage       *decimal.Decimal `json:"leverage,omitempty"`
	Misc           string           `json:"misc,omitempty"`
	Ledgers        []string         `json:"ledgers,omitempty"`
	TradeID        *decimal.Decimal `json:"trade_id,omitempty"`
	Maker          bool             `json:"maker,omitempty"`
	PositionStatus string           `json:"posstatus,omitempty"`
	CPrice         *decimal.Decimal `json:"cprice,omitempty"`
	CCost          *decimal.Decimal `json:"ccost,omitempty"`
	CFee           *decimal.Decimal `json:"cfee,omitempty"`
	CVol           *decimal.Decimal `json:"cvol,omitempty"`
	CMargin        *decimal.Decimal `json:"cmargin,omitempty"`
	Net            *decimal.Decimal `json:"net,omitempty"`
	Trades         []string         `json:"trades,omitempty"`
}

type OrderDescriptionInfo struct {
	Order string `json:"order,omitempty"`
}

type OrderDescriptionInfoWithClose struct {
	Close string `json:"close,omitempty"`
	OrderDescriptionInfo
}

type OrderDescription struct {
	Pair           string           `json:"pair,omitempty"`
	Type           string           `json:"type,omitempty"`
	OrderType      string           `json:"ordertype,omitempty"`
	Price          *decimal.Decimal `json:"price,omitempty"`
	SecondaryPrice *decimal.Decimal `json:"price2,omitempty"`
	Leverage       string           `json:"leverage,omitempty"`
	OrderDescriptionInfoWithClose
}

type Order struct {
	RefID          string            `json:"refid,omitempty"`
	UserRef        *decimal.Decimal  `json:"userref,omitempty"`
	ClOrdID        string            `json:"cl_ord_id,omitempty"`
	Status         string            `json:"status,omitempty"`
	OpenTm         *decimal.Decimal  `json:"opentm,omitempty"`
	StartTm        *decimal.Decimal  `json:"starttm,omitempty"`
	ExpireTm       *decimal.Decimal  `json:"expiretm,omitempty"`
	Description    *OrderDescription `json:"descr,omitempty"`
	Volume         *decimal.Decimal  `json:"vol,omitempty"`
	VolumeExecuted *decimal.Decimal  `json:"vol_exec,omitempty"`
	Cost           *decimal.Decimal  `json:"cost,omitempty"`
	Fee            *decimal.Decimal  `json:"fee,omitempty"`
	Price          *decimal.Decimal  `json:"price,omitempty"`
	StopPrice      *decimal.Decimal  `json:"stopprice,omitempty"`
	LimitPrice     *decimal.Decimal  `json:"limitprice,omitempty"`
	Trigger        string            `json:"trigger,omitempty"`
	Margin         bool              `json:"margin,omitempty"`
	Misc           string            `json:"misc,omitempty"`
	SenderSubID    string            `json:"sender_sub_id,omitempty"`
	OrderFlags     string            `json:"oflags,omitempty"`
	Trades         []string          `json:"trades,omitempty"`
}

type ClosedOrder struct {
	CloseTm *decimal.Decimal `json:"closetm,omitempty"`
	Reason  string           `json:"reason,omitempty"`
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
	AssetClass      string           `json:"aclass,omitempty"`
	AltName         string           `json:"altname,omitempty"`
	Decimals        int              `json:"decimals,omitempty"`
	DisplayDecimals int              `json:"display_decimals,omitempty"`
	CollateralValue *decimal.Decimal `json:"collateral_value,omitempty"`
	Status          string           `json:"status,omitempty"`
}

type JSONAssetPair struct {
	AltName            string               `json:"altname,omitempty"`
	WSName             string               `json:"wsname,omitempty"`
	BaseAssetClass     string               `json:"aclass_base,omitempty"`
	Base               string               `json:"base,omitempty"`
	QuoteAssetClass    string               `json:"aclass_quote,omitempty"`
	Quote              string               `json:"quote,omitempty"`
	PairDecimals       int                  `json:"pair_decimals,omitempty"`
	CostDecimals       int                  `json:"cost_decimals,omitempty"`
	LotDecimals        int                  `json:"lot_decimals,omitempty"`
	LotMultiplier      int                  `json:"lot_multiplier,omitempty"`
	BuyLeverage        []int                `json:"leverage_buy,omitempty"`
	SellLeverage       []int                `json:"leverage_sell,omitempty"`
	Fees               [][]*decimal.Decimal `json:"fees,omitempty"`
	FeesMaker          [][]*decimal.Decimal `json:"fees_maker,omitempty"`
	FeeVolumeCurrency  string               `json:"fee_volume_currency,omitempty"`
	MarginCall         int                  `json:"margin_call,omitempty"`
	MarginStop         int                  `json:"margin_stop,omitempty"`
	OrderMinimum       *decimal.Decimal     `json:"ordermin,omitempty"`
	CostMinimum        *decimal.Decimal     `json:"costmin,omitempty"`
	TickSize           *decimal.Decimal     `json:"tick_size,omitempty"`
	Status             string               `json:"status,omitempty"`
	LongPositionLimit  int                  `json:"long_position_limit,omitempty"`
	ShortPositionLimit int                  `json:"short_position_limit,omitempty"`
}

func (jap JSONAssetPair) AssetPair() AssetPair {
	fees := make([]Fee, len(jap.Fees))
	for i, fee := range jap.Fees {
		fees[i] = Fee{
			Volume:  fee[0],
			Percent: fee[1],
		}
	}
	feesMaker := make([]Fee, len(jap.FeesMaker))
	for i, fee := range jap.FeesMaker {
		feesMaker[i] = Fee{
			Volume:  fee[0],
			Percent: fee[1],
		}
	}
	return AssetPair{
		AltName:            jap.AltName,
		WSName:             jap.WSName,
		BaseAssetClass:     jap.BaseAssetClass,
		Base:               jap.Base,
		QuoteAssetClass:    jap.QuoteAssetClass,
		Quote:              jap.Quote,
		PairDecimals:       jap.PairDecimals,
		CostDecimals:       jap.CostDecimals,
		LotDecimals:        jap.LotDecimals,
		LotMultiplier:      jap.LotMultiplier,
		BuyLeverage:        jap.BuyLeverage,
		SellLeverage:       jap.SellLeverage,
		Fees:               fees,
		FeesMaker:          feesMaker,
		FeeVolumeCurrency:  jap.FeeVolumeCurrency,
		MarginCall:         jap.MarginCall,
		MarginStop:         jap.MarginStop,
		OrderMinimum:       jap.OrderMinimum,
		CostMinimum:        jap.CostMinimum,
		TickSize:           jap.TickSize,
		Status:             jap.Status,
		LongPositionLimit:  jap.LongPositionLimit,
		ShortPositionLimit: jap.ShortPositionLimit,
	}
}

type AssetPair struct {
	AltName            string           `json:"altname,omitempty"`
	WSName             string           `json:"wsname,omitempty"`
	BaseAssetClass     string           `json:"aclass_base,omitempty"`
	Base               string           `json:"base,omitempty"`
	QuoteAssetClass    string           `json:"aclass_quote,omitempty"`
	Quote              string           `json:"quote,omitempty"`
	PairDecimals       int              `json:"pair_decimals,omitempty"`
	CostDecimals       int              `json:"cost_decimals,omitempty"`
	LotDecimals        int              `json:"lot_decimals,omitempty"`
	LotMultiplier      int              `json:"lot_multiplier,omitempty"`
	BuyLeverage        []int            `json:"leverage_buy,omitempty"`
	SellLeverage       []int            `json:"leverage_sell,omitempty"`
	Fees               []Fee            `json:"fees,omitempty"`
	FeesMaker          []Fee            `json:"fees_maker,omitempty"`
	FeeVolumeCurrency  string           `json:"fee_volume_currency,omitempty"`
	MarginCall         int              `json:"margin_call,omitempty"`
	MarginStop         int              `json:"margin_stop,omitempty"`
	OrderMinimum       *decimal.Decimal `json:"ordermin,omitempty"`
	CostMinimum        *decimal.Decimal `json:"costmin,omitempty"`
	TickSize           *decimal.Decimal `json:"tick_size,omitempty"`
	Status             string           `json:"status,omitempty"`
	LongPositionLimit  int              `json:"long_position_limit,omitempty"`
	ShortPositionLimit int              `json:"short_position_limit,omitempty"`
}

func (ap *AssetPair) UnmarshalJSON(data []byte) error {
	var v JSONAssetPair
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*ap = v.AssetPair()
	return nil
}

type Fee struct {
	Volume  *decimal.Decimal `json:"volume,omitempty"`
	Percent *decimal.Decimal `json:"percent_fee,omitempty"`
}

type AssetTickerInfo struct {
	Ask    []*decimal.Decimal `json:"a,omitempty"`
	Bid    []*decimal.Decimal `json:"b,omitempty"`
	Close  []*decimal.Decimal `json:"c,omitempty"`
	Volume []*decimal.Decimal `json:"v,omitempty"`
	VWAP   []*decimal.Decimal `json:"p,omitempty"`
	Trades []int              `json:"t,omitempty"`
	Low    []*decimal.Decimal `json:"l,omitempty"`
	High   []*decimal.Decimal `json:"h,omitempty"`
	Open   *decimal.Decimal   `json:"o,omitempty"`
}

type JSONOrderBook struct {
	Asks [][]*decimal.Decimal `json:"asks,omitempty"`
	Bids [][]*decimal.Decimal `json:"bids,omitempty"`
}

func (job JSONOrderBook) OrderBook() (book OrderBook) {
	book.Asks = make([]PriceLevel, len(job.Asks))
	for i, ask := range job.Asks {
		book.Asks[i] = PriceLevel{
			Price:     ask[0],
			Volume:    ask[1],
			Timestamp: time.Unix(ask[2].Int64(), 0),
		}
	}
	book.Bids = make([]PriceLevel, len(job.Bids))
	for i, bid := range job.Bids {
		book.Bids[i] = PriceLevel{
			Price:     bid[0],
			Volume:    bid[1],
			Timestamp: time.Unix(bid[2].Int64(), 0),
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
	Price     *decimal.Decimal `json:"price,omitempty"`
	Volume    *decimal.Decimal `json:"volume,omitempty"`
	Timestamp time.Time        `json:"timestamp,omitempty"`
}

type DepositMethod struct {
	Method        string           `json:"method,omitempty"`
	Limit         *decimal.Decimal `json:"limit,omitempty"`
	Fee           *decimal.Decimal `json:"fee,omitempty"`
	FeePercentage *decimal.Decimal `json:"fee-percentage,omitempty"`
	GenAddress    bool             `json:"gen-address,omitempty"`
	Minimum       *decimal.Decimal `json:"minimum,omitempty"`
}

type DepositAddress struct {
	Address  string `json:"address,omitempty"`
	ExpireTM string `json:"expiretm,omitempty"`
	New      bool   `json:"new,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

type JSONDepositStatus struct {
	Method      string           `json:"method,omitempty"`
	Aclass      string           `json:"aclass,omitempty"`
	Asset       string           `json:"asset,omitempty"`
	RefID       string           `json:"refid,omitempty"`
	TxID        string           `json:"txid,omitempty"`
	Info        string           `json:"info,omitempty"`
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	Fee         *decimal.Decimal `json:"fee,omitempty"`
	Time        int              `json:"time,omitempty"`
	Status      string           `json:"status,omitempty"`
	StatusProp  string           `json:"status-prop,omitempty"`
	Originators []string         `json:"originators,omitempty"`
}

func (jds JSONDepositStatus) DepositStatus() DepositStatus {
	return DepositStatus{
		Method:      jds.Method,
		Aclass:      jds.Aclass,
		Asset:       jds.Asset,
		RefID:       jds.RefID,
		TxID:        jds.TxID,
		Info:        jds.Info,
		Amount:      jds.Amount,
		Fee:         jds.Fee,
		Time:        time.Unix(int64(jds.Time), 0),
		Status:      jds.Status,
		StatusProp:  jds.StatusProp,
		Originators: jds.Originators,
	}
}

type DepositStatus struct {
	Method      string           `json:"method,omitempty"`
	Aclass      string           `json:"aclass,omitempty"`
	Asset       string           `json:"asset,omitempty"`
	RefID       string           `json:"refid,omitempty"`
	TxID        string           `json:"txid,omitempty"`
	Info        string           `json:"info,omitempty"`
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	Fee         *decimal.Decimal `json:"fee,omitempty"`
	Time        time.Time        `json:"time,omitempty"`
	Status      string           `json:"status,omitempty"`
	StatusProp  string           `json:"status-prop,omitempty"`
	Originators []string         `json:"originators,omitempty"`
}

func (ds *DepositStatus) UnmarshalJSON(data []byte) error {
	var v JSONDepositStatus
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*ds = v.DepositStatus()
	return nil
}

type WithdrawMethod struct {
	Asset   string           `json:"asset,omitempty"`
	Method  string           `json:"method,omitempty"`
	Network string           `json:"network,omitempty"`
	Minimum *decimal.Decimal `json:"minimum,omitempty"`
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
	Method string           `json:"method,omitempty"`
	Limit  *decimal.Decimal `json:"limit,omitempty"`
	Amount *decimal.Decimal `json:"amount,omitempty"`
	Fee    *decimal.Decimal `json:"fee,omitempty"`
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
type EarnStrategyLockType struct {
	Type                    string           `json:"type,omitempty"`
	BondingPeriod           *decimal.Decimal `json:"bonding_period,omitempty"`
	BondingPeriodVariable   bool             `json:"bonding_period_variable,omitempty"`
	BondingRewards          bool             `json:"bonding_rewards,omitempty"`
	ExitQueuePeriod         *decimal.Decimal `json:"exit_queue_period"`
	PayoutFrequency         *decimal.Decimal `json:"payout_frequency,omitempty"`
	UnbondingPeriod         *decimal.Decimal `json:"unbonding_period,omitempty"`
	UnbondingPeriodVariable bool             `json:"unbonding_period_variable,omitempty"`
	UnbondingRewards        bool             `json:"unbonding_rewards,omitempty"`
}

type EarnStrategyYieldSource struct {
	Type string `json:"type,omitempty"`
}

type EarnStrategyAutoCompound struct {
	Type    string `json:"type,omitempty"`
	Default bool   `json:"default,omitempty"`
}

type EarnStrategy struct {
	ID                        string                   `json:"id,omitempty"`
	Asset                     string                   `json:"asset,omitempty"`
	LockType                  EarnStrategyLockType     `json:"lock_type,omitempty"`
	AprEstimate               *EarnStrategyAPR         `json:"apr_estimate,omitempty"`
	UserMinAllocation         string                   `json:"user_min_allocation,omitempty"`
	AllocationFee             *decimal.Decimal         `json:"allocation_fee,omitempty"`
	DeallocationFee           *decimal.Decimal         `json:"deallocation_fee,omitempty"`
	AutoCompound              EarnStrategyAutoCompound `json:"auto_compound,omitempty"`
	YieldSource               EarnStrategyYieldSource  `json:"yield_source,omitempty"`
	CanAllocate               bool                     `json:"can_allocate,omitempty"`
	CanDeallocate             bool                     `json:"can_deallocate,omitempty"`
	AllocationRestrictionInfo []string                 `json:"allocation_restriction_info,omitempty"`
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
