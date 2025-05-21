package derivatives

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// AssetManager provides helper methods for a group of [Instrument] objects.
type AssetManager struct {
	instruments map[string]Instrument
	mux         sync.RWMutex
}

// NewAssetManager constructs a new [AssetManager] object.
// The map will need to be initialized with [AssetManager.Use] or [AssetManager.Update].
func NewAssetManager() *AssetManager {
	return &AssetManager{
		instruments: make(map[string]Instrument),
	}
}

// Use retrieves instrument specifications using the specified [REST] structure.
func (m *AssetManager) Use(r *REST) error {
	resp, err := r.Instruments()
	if err != nil {
		return err
	}
	m.Update(resp.Instruments)
	return nil
}

// Update creates and stores a map of [Instrument] objects.
func (m *AssetManager) Update(update []Instrument) {
	instruments := make(map[string]Instrument)
	for _, i := range update {
		instruments[strings.ToUpper(i.Symbol)] = i
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	m.instruments = instruments
}

// PairInfo returns the [AssetPair] struct corresponding to the symbol.
func (m *AssetManager) Info(symbol string) (*Instrument, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	instrument, ok := m.instruments[strings.ToUpper(symbol)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &instrument, nil
}

// FormatPrice sets the decimals to the tick size of the contract.
func (m *AssetManager) FormatPrice(symbol string, v *kraken.Money) (*kraken.Money, error) {
	info, err := m.Info(symbol)
	if err != nil {
		return nil, err
	}
	return v.SetSize(info.TickSize), nil
}

// FormatSize sets the decimals to the lot size of the contract.
func (m *AssetManager) FormatSize(symbol string, v *kraken.Money) (*kraken.Money, error) {
	info, err := m.Info(symbol)
	if err != nil {
		return nil, err
	}
	if info.ContractValueTradePrecision.Sign() == -1 {
		v = v.SetGranularity(kraken.NewMoneyFromInt64(10).Exp(info.ContractValueTradePrecision.Abs()).Int64())
	}
	return v.SetDecimals(int64(math.Max(info.ContractValueTradePrecision.Float64(), 0))), nil
}
