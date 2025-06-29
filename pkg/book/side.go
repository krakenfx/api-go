package book

import (
	"github.com/krakenfx/api-go/v2/pkg/decimal"
)

// Side encompasses the price levels in one side of the book.
type Side struct {
	Direction BookDirection     `json:"direction,omitempty"`
	High      *Level            `json:"high,omitempty"`
	Low       *Level            `json:"low,omitempty"`
	Last      *Level            `json:"last,omitempty"`
	Levels    map[string]*Level `json:"levels,omitempty"`
}

// NewSide constructs a new [Side] with default values.
func NewSide() *Side {
	return &Side{
		Levels: make(map[string]*Level),
	}
}

// FindAdjacent finds the nearest price level close to the given price.
func (s *Side) FindAdjacent(price *decimal.Decimal) *Level {
	if s.High == nil || s.Low == nil {
		return nil
	}
	if price.Cmp(s.High.Price) > 0 {
		return s.High
	}
	if price.Cmp(s.Low.Price) < 0 {
		return s.Low
	}
	highDistance := s.High.Price.Sub(price)
	lowDistance := price.Sub(s.Low.Price)
	if highDistance.Cmp(lowDistance) == 1 {
		return s.FindAdjacentBelow(price)
	} else {
		return s.FindAdjacentAbove(price)
	}
}

// FindAdjacentBelow finds the nearest price level from below the given price.
func (s *Side) FindAdjacentBelow(price *decimal.Decimal) *Level {
	if s.Low == nil || price.Cmp(s.Low.Price) <= 0 {
		return nil
	}
	nearest := s.Low
	for {
		next := nearest.Higher
		if next == nil {
			break
		}
		nearestDiff := nearest.Price.Sub(price).Abs()
		nextDiff := next.Price.Sub(price).Abs()
		if nearestDiff.Cmp(nextDiff) < 0 || next.Price.Cmp(price) >= 0 {
			break
		}
		nearest = next
	}
	return nearest
}

// FindAdjacentAbove finds the nearest price level above the given price.
func (s *Side) FindAdjacentAbove(price *decimal.Decimal) *Level {
	if s.High == nil || price.Cmp(s.High.Price) >= 0 {
		return nil
	}
	nearest := s.High
	for {
		next := nearest.Lower
		if next == nil {
			break
		}
		nearestDiff := nearest.Price.Sub(price).Abs()
		nextDiff := next.Price.Sub(price).Abs()
		if nearestDiff.Cmp(nextDiff) < 0 || next.Price.Cmp(price) <= 0 {
			break
		}
		nearest = next
	}
	return nearest
}

func (s *Side) insert(opts *UpdateOptions) {
	level := NewLevel()
	level.update(opts)
	nearest := s.FindAdjacent(level.Price)
	if nearest == nil || s.High == nil || level.Price.Cmp(s.High.Price) > 0 {
		s.High = level
	}
	if nearest == nil || s.Low == nil || level.Price.Cmp(s.Low.Price) < 0 {
		s.Low = level
	}
	if nearest != nil {
		if level.Price.Cmp(nearest.Price) > 0 {
			level.Lower = nearest
			level.Higher = nearest.Higher
			nearest.Higher = level
			if level.Higher != nil {
				level.Higher.Lower = level
			}
		} else if level.Price.Cmp(nearest.Price) < 0 {
			level.Higher = nearest
			level.Lower = nearest.Lower
			nearest.Lower = level
			if level.Lower != nil {
				level.Lower.Higher = level
			}
		}
	}
	s.Levels[level.Price.String()] = level
}

// Update interprets a [UpdateOptions] message and decides if it should insert, update, or delete the price level.
func (s *Side) update(opts *UpdateOptions) {
	level, ok := s.Levels[opts.Price.String()]
	if !ok && opts.Quantity.Sign() == 1 {
		s.insert(opts)
	} else {
		level.update(opts)
	}
	if level != nil && level.Quantity.Sign() <= 0 {
		s.delete(level)
	}
}

func (s *Side) delete(level *Level) {
	if s.High == level {
		s.High = level.Lower
	}
	if s.Low == level {
		s.Low = level.Higher
	}
	if level.Lower != nil {
		level.Lower.Higher = level.Higher
	}
	if level.Higher != nil {
		level.Higher.Lower = level.Lower
	}
	delete(s.Levels, level.Price.String())
}
