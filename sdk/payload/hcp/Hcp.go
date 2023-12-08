package hcp

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

type Hcp struct {
	PositionOpen     bool
	PositionClose    bool
	OptionRelais     bool
	LightBarrier     bool
	Error            bool
	DrivingToClose   bool
	Driving          bool
	HalfOpened       bool
	ForecastLeadTime bool
	Learned          bool
	NotReferenced    bool
}

func DecodeHcp(payloadBytes []byte) *Hcp {
	bitset := binary.LittleEndian.Uint16(payloadBytes)
	return &Hcp{
		PositionOpen:     BitToBool(bitset, 0),
		PositionClose:    BitToBool(bitset, 1),
		OptionRelais:     BitToBool(bitset, 2),
		LightBarrier:     BitToBool(bitset, 3),
		Error:            BitToBool(bitset, 4),
		DrivingToClose:   BitToBool(bitset, 5),
		Driving:          BitToBool(bitset, 6),
		HalfOpened:       BitToBool(bitset, 7),
		ForecastLeadTime: BitToBool(bitset, 8),
		Learned:          BitToBool(bitset, 9),
		NotReferenced:    BitToBool(bitset, 10),
	}
}

func BitToBool(bitset uint16, index int) bool {
	bitMask := uint16(1 << index)
	masked := bitset & bitMask
	return masked == bitMask
}

func (h *Hcp) String() string {
	return fmt.Sprintf("HCP[PositionOpen: %v, PositionClose: %v, OptionRelais: %v, LightBarrier: %v, Error: %v, DrivingToClose: %v, Driving: %v, HalfOpened: %v, ForecastLeadTime: %v, Learned: %v, NotReferenced: %v]",
		h.PositionOpen,
		h.PositionClose,
		h.OptionRelais,
		h.LightBarrier,
		h.Error,
		h.DrivingToClose,
		h.Driving,
		h.HalfOpened,
		h.ForecastLeadTime,
		h.Learned,
		h.NotReferenced)
}

func (h *Hcp) Equal(o *Hcp) bool {
	return reflect.DeepEqual(h, o)
}
