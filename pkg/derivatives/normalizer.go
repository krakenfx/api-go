package derivatives

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/krakenfx/api-go/pkg/decimal"
)

// Normalizer provides helper methods for a group of [Instrument] objects.
type Normalizer struct {
	instruments map[string]Instrument
	mux         sync.RWMutex
}

// NewNormalizer constructs a new [Normalizer] object.
// The map will need to be initialized with [Normalizer.Use] or [Normalizer.Update].
func NewNormalizer() *Normalizer {
	return &Normalizer{
		instruments: make(map[string]Instrument),
	}
}

// Use retrieves instrument specifications using the specified [REST] structure.
func (m *Normalizer) Use(r *REST) error {
	resp, err := r.Instruments()
	if err != nil {
		return err
	}
	m.Update(resp.Result.Instruments)
	return nil
}

// Update creates and stores a map of [Instrument] objects.
func (m *Normalizer) Update(update []Instrument) {
	instruments := make(map[string]Instrument)
	for _, i := range update {
		instruments[strings.ToUpper(i.Symbol)] = i
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	m.instruments = instruments
}

// PairInfo returns the [AssetPair] struct corresponding to the symbol.
func (m *Normalizer) Info(symbol string) (*Instrument, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	instrument, ok := m.instruments[strings.ToUpper(symbol)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &instrument, nil
}

// FormatPrice sets the decimals to the tick size of the contract.
func (m *Normalizer) FormatPrice(symbol string, v *decimal.Decimal) (*decimal.Decimal, error) {
	info, err := m.Info(symbol)
	if err != nil {
		return nil, err
	}
	return v.SetSize(info.TickSize), nil
}

// FormatSize sets the decimals to the lot size of the contract.
func (m *Normalizer) FormatSize(symbol string, v *decimal.Decimal) (*decimal.Decimal, error) {
	info, err := m.Info(symbol)
	if err != nil {
		return nil, err
	}
	if info.ContractValueTradePrecision.Sign() == -1 {
		v = v.SetIncrement(decimal.NewFromInt64(10).Pow(info.ContractValueTradePrecision.Abs()).Int64())
	}
	return v.SetScale(int64(math.Max(info.ContractValueTradePrecision.Float64(), 0))), nil
}
