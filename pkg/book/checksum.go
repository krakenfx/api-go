package book

import (
	"fmt"
	"hash/crc32"
	"strings"
)

// ChecksumPart contains information regarding a price level or order.
type ChecksumPart struct {
	Level        *Level `json:"-"`
	Order        *Order `json:"order,omitempty"`
	Price        string `json:"price,omitempty"`
	Quantity     string `json:"quantity,omitempty"`
	Concatenated string `json:"concatenated,omitempty"`
}

// ChecksumResult contains the result of the book checksum validation.
type ChecksumResult struct {
	Level          int             `json:"level,omitempty"`
	ServerChecksum string          `json:"serverChecksum,omitempty"`
	LocalChecksum  string          `json:"localChecksum,omitempty"`
	Match          bool            `json:"match,omitempty"`
	AskParts       []*ChecksumPart `json:"askParts,omitempty"`
	BidParts       []*ChecksumPart `json:"bidParts,omitempty"`
	Asks           string          `json:"asks,omitempty"`
	Bids           string          `json:"bids,omitempty"`
}

// L2Checksum verifies that the L2 book is synchronized with the exchange.
//
// https://docs.kraken.com/api/docs/guides/spot-ws-book-v2
func (b *Book) L2Checksum(checksum string) *ChecksumResult {
	result := &ChecksumResult{
		Level:          2,
		ServerChecksum: checksum,
	}
	cursor := b.BestAsk()
	for range 10 {
		if cursor == nil {
			break
		}
		price := strings.TrimLeft(strings.ReplaceAll(cursor.GetPriceString(), ".", ""), "0")
		quantity := strings.TrimLeft(strings.ReplaceAll(cursor.GetQuantityString(), ".", ""), "0")
		concatenated := price + quantity
		result.AskParts = append(result.AskParts, &ChecksumPart{
			Level:        cursor,
			Price:        price,
			Quantity:     quantity,
			Concatenated: concatenated,
		})
		result.Asks += concatenated
		cursor = cursor.Higher
	}
	cursor = b.BestBid()
	for range 10 {
		if cursor == nil {
			break
		}
		price := strings.TrimLeft(strings.ReplaceAll(cursor.GetPriceString(), ".", ""), "0")
		quantity := strings.TrimLeft(strings.ReplaceAll(cursor.GetQuantityString(), ".", ""), "0")
		concatenated := price + quantity
		result.BidParts = append(result.BidParts, &ChecksumPart{
			Level:        cursor,
			Price:        price,
			Quantity:     quantity,
			Concatenated: concatenated,
		})
		result.Bids += concatenated
		cursor = cursor.Lower
	}
	result.LocalChecksum = fmt.Sprint(crc32.Checksum([]byte(result.Asks+result.Bids), crc32.IEEETable))
	if result.LocalChecksum == result.ServerChecksum {
		result.Match = true
	}
	b.OnChecksummed.Call(result)
	return result
}

// L3Checksum verifies that the L3 book is synchronized with the exchange.
//
// https://docs.kraken.com/api/docs/guides/spot-ws-l3-v2
func (b *Book) L3Checksum(checksum string) *ChecksumResult {
	result := &ChecksumResult{
		Level:          3,
		ServerChecksum: checksum,
	}
	cursor := b.BestAsk()
	for range 10 {
		if cursor == nil {
			break
		}
		for _, order := range cursor.Queue() {
			price := strings.TrimLeft(strings.ReplaceAll(order.LimitPrice.String(), ".", ""), "0")
			quantity := strings.TrimLeft(strings.ReplaceAll(order.Quantity.String(), ".", ""), "0")
			concatenated := price + quantity
			result.AskParts = append(result.AskParts, &ChecksumPart{
				Level:        cursor,
				Order:        order,
				Price:        price,
				Quantity:     quantity,
				Concatenated: concatenated,
			})
			result.Asks += concatenated
		}
		cursor = cursor.Higher
	}
	cursor = b.BestBid()
	for range 10 {
		if cursor == nil {
			break
		}
		for _, order := range cursor.Queue() {
			price := strings.TrimLeft(strings.ReplaceAll(order.LimitPrice.String(), ".", ""), "0")
			quantity := strings.TrimLeft(strings.ReplaceAll(order.Quantity.String(), ".", ""), "0")
			concatenated := price + quantity
			result.BidParts = append(result.BidParts, &ChecksumPart{
				Level:        cursor,
				Order:        order,
				Price:        price,
				Quantity:     quantity,
				Concatenated: concatenated,
			})
			result.Bids += concatenated
		}
		cursor = cursor.Lower
	}
	result.LocalChecksum = fmt.Sprint(crc32.Checksum([]byte(result.Asks+result.Bids), crc32.IEEETable))
	if result.LocalChecksum == result.ServerChecksum {
		result.Match = true
	}
	b.OnChecksummed.Call(result)
	return result
}
