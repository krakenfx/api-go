package spot

import (
	"fmt"
	"maps"
	"strings"
	"sync"

	"github.com/krakenfx/api-go/pkg/kraken"
	"golang.org/x/sync/errgroup"
)

// AssetManager provides helper methods for a group of [AssetInfo] objects.
type AssetManager struct {
	aliases map[string]*AssetName
	assets  map[*AssetName]AssetInfo
	pairs   map[*AssetName]map[*AssetName]AssetPair
	mux     sync.RWMutex
}

// NewAssetManager constructs a new [AssetManager] object.
// The map will need to be initialized with [AssetManager.Use] or [AssetManager.Update].
func NewAssetManager() *AssetManager {
	return &AssetManager{
		aliases: make(map[string]*AssetName),
		assets:  make(map[*AssetName]AssetInfo),
		pairs:   make(map[*AssetName]map[*AssetName]AssetPair),
	}
}

// AssetName contains all possible names for an asset.
type AssetName struct {
	Name    string
	AltName string
	OldName string
}

// Use retrieves asset specifications using the specified [REST] structure.
func (m *AssetManager) Use(r *REST) error {
	update := &AssetsManagerUpdate{}
	var wg errgroup.Group
	wg.Go(func() error {
		resp, err := r.Assets(nil)
		if err != nil {
			return fmt.Errorf("old assets: %w", err)
		}
		update.OldAssets = resp.Result
		return nil
	})
	wg.Go(func() error {
		resp, err := r.Assets(&AssetsRequest{
			AssetVersion: 1,
		})
		if err != nil {
			return fmt.Errorf("new assets: %w", err)
		}
		update.NewAssets = resp.Result
		return nil
	})
	wg.Go(func() error {
		resp, err := r.AssetPairs(nil)
		if err != nil {
			return fmt.Errorf("old pairs: %w", err)
		}
		update.OldPairs = resp.Result
		return nil
	})
	wg.Go(func() error {
		resp, err := r.AssetPairs(&AssetPairsRequest{
			AssetVersion: 1,
		})
		if err != nil {
			return fmt.Errorf("new pairs: %w", err)
		}
		update.NewPairs = resp.Result
		return nil
	})
	if err := wg.Wait(); err != nil {
		return err
	}
	m.Update(update)
	return nil
}

type AssetsManagerUpdate struct {
	OldAssets map[string]AssetInfo
	NewAssets map[string]AssetInfo
	OldPairs  map[string]AssetPair
	NewPairs  map[string]AssetPair
}

// Update creates an alias map out of the old and new [AssetInfo] map.
func (m *AssetManager) Update(update *AssetsManagerUpdate) {
	aliases := make(map[string]*AssetName)
	info := make(map[*AssetName]AssetInfo)
	pairs := make(map[*AssetName]map[*AssetName]AssetPair)
	for name, assetInfo := range update.NewAssets {
		assetName := &AssetName{
			Name:    name,
			AltName: assetInfo.AltName,
		}
		aliases[assetName.Name] = assetName
		aliases[assetName.AltName] = assetName
		info[assetName] = assetInfo
	}
	for name, assetInfo := range update.OldAssets {
		assetName, ok := aliases[assetInfo.AltName]
		if !ok {
			continue
		}
		assetName.OldName = name
		aliases[name] = assetName
	}
	for _, pair := range update.NewPairs {
		baseAlt, quoteAlt, found := strings.Cut(pair.WSName, "/")
		if !found {
			continue
		}
		baseName, ok := aliases[baseAlt]
		if !ok {
			baseName = &AssetName{
				Name:    pair.Base,
				AltName: baseAlt,
			}
			aliases[baseAlt] = baseName
		}
		quoteName, ok := aliases[quoteAlt]
		if !ok {
			quoteName = &AssetName{
				Name:    pair.Quote,
				AltName: quoteAlt,
			}
			aliases[quoteAlt] = quoteName
		}
		baseMap, ok := pairs[baseName]
		if !ok {
			baseMap = make(map[*AssetName]AssetPair)
			pairs[baseName] = baseMap
		}
		baseMap[quoteName] = pair
	}
	for _, pair := range update.OldPairs {
		baseAlt, quoteAlt, found := strings.Cut(pair.WSName, "/")
		if !found {
			continue
		}
		baseName, ok := aliases[baseAlt]
		if !ok {
			continue
		}
		baseName.OldName = pair.Base
		quoteName, ok := aliases[quoteAlt]
		if !ok {
			continue
		}
		quoteName.OldName = pair.Quote
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	m.aliases = aliases
	m.assets = info
	m.pairs = pairs
}

// Map returns a group of [AssetInfo] structs mapped to their [AssetName].
func (m *AssetManager) Map() map[*AssetName]AssetInfo {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return maps.Clone(m.assets)
}

// AssetName returns the matching [AssetName]. If not found, returns false.
func (m *AssetManager) AssetName(name string) (*AssetName, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	assetName, ok := m.aliases[strings.ToUpper(name)]
	if !ok {
		return nil, false
	}
	return assetName, true
}

// PairName returns the two [AssetName] structures corresponding to the name. If not found, returns false.
func (m *AssetManager) PairName(name string) (*AssetName, *AssetName, bool) {
	base, quote, found := strings.Cut(name, "/")
	if found {
		baseName, found := m.AssetName(base)
		if !found {
			return nil, nil, false
		}
		quoteName, found := m.AssetName(quote)
		if !found {
			return nil, nil, false
		}
		return baseName, quoteName, true
	} else {
		m.mux.RLock()
		defer m.mux.RUnlock()
		for baseAlias, baseName := range m.aliases {
			for quoteAlias, quoteName := range m.aliases {
				if strings.EqualFold(name, baseAlias+quoteAlias) {
					return baseName, quoteName, true
				}
			}
		}
	}
	return nil, nil, false
}

// Name returns the standard asset or pair name.
func (m *AssetManager) Name(name string) string {
	if asset, found := m.AssetName(name); found {
		return asset.Name
	} else if base, quote, found := m.PairName(name); found {
		return base.Name + "/" + quote.Name
	} else {
		return strings.ToUpper(name)
	}
}

// AssetInfo returns the [AssetInfo] struct corresponding to the name.
func (m *AssetManager) AssetInfo(name string) (*AssetInfo, error) {
	asset, found := m.AssetName(name)
	if !found {
		return nil, fmt.Errorf("not found")
	}
	m.mux.RLock()
	defer m.mux.RUnlock()
	info := m.assets[asset]
	return &info, nil
}

// PairInfo returns the [AssetPair] struct corresponding to the symbol.
func (m *AssetManager) PairInfo(symbol string) (*AssetPair, error) {
	base, quote, found := m.PairName(symbol)
	if !found {
		return nil, fmt.Errorf("name not found")
	}
	m.mux.RLock()
	defer m.mux.RUnlock()
	baseMapped, ok := m.pairs[base]
	if !ok {
		return nil, fmt.Errorf("base not found")
	}
	pair, ok := baseMapped[quote]
	if !ok {
		return nil, fmt.Errorf("quote not found")
	}
	return &pair, nil
}

// FormatDecimals sets the decimals by the decimal property of the symbol.
func (m *AssetManager) FormatDecimals(symbol string, v *kraken.Money) (*kraken.Money, error) {
	info, err := m.AssetInfo(symbol)
	if err != nil {
		return nil, err
	}
	return v.SetDecimals(int64(info.Decimals)), nil
}

// FormatDisplayDecimals sets the decimals by the display decimals property of the symbol.
func (m *AssetManager) FormatDisplayDecimals(symbol string, v *kraken.Money) (*kraken.Money, error) {
	info, err := m.AssetInfo(symbol)
	if err != nil {
		return nil, err
	}
	return v.SetDecimals(int64(info.DisplayDecimals)), nil
}

// FormatPrice sets the decimals to the tick size of the asset pair.
func (m *AssetManager) FormatPrice(pair string, v *kraken.Money) (*kraken.Money, error) {
	info, err := m.PairInfo(pair)
	if err != nil {
		return nil, err
	}
	return v.
		SetSize(info.TickSize).
		SetDecimals(int64(info.PairDecimals)), nil
}

// FormatSize sets the decimals to the lot size of the asset pair.
func (m *AssetManager) FormatSize(pair string, v *kraken.Money) (*kraken.Money, error) {
	info, err := m.PairInfo(pair)
	if err != nil {
		return nil, err
	}
	return v.
		SetDecimals(int64(info.LotDecimals)).
		SetGranularity(int64(info.LotMultiplier)), nil
}
