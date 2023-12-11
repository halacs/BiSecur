package payload

import (
	"bisecure/sdk/payload/hcp"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type HmGetTransitionResponse struct {
	Payload
	StateInPercent        int // 100 is OPEN, 0 = CLOSED, 200 = UNLOCKED, 0 = LOCKED????
	DesiredStateInPercent int // 100 is OPEN, 0 = CLOSED
	Error                 bool
	AutoClose             bool
	DriveTime             int
	Gk                    uint16
	Hcp                   *hcp.Hcp
	Exst                  []byte
	Time                  time.Time
	//IgnoreRetries         bool
}

func DecodeHmGetTransitionResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	hcpPayloadBytes := payloadBytes[6:8]
	h := hcp.DecodeHcp(hcpPayloadBytes)

	hmgtr := &HmGetTransitionResponse{
		Payload: Payload{
			data:       payloadBytes,
			dataLength: byte(len(payloadBytes)),
		},
		StateInPercent:        int(payloadBytes[0]) / 2,
		DesiredStateInPercent: int(payloadBytes[1]) / 2,
		Error:                 hcp.BitToBool(uint16(payloadBytes[2]), 7),
		AutoClose:             hcp.BitToBool(uint16(payloadBytes[2]), 6),
		DriveTime:             int(payloadBytes[3]), // TODO: from kotlin code: "TODO: clear 6th and 7th bit from byte3 and shift add it here"
		Gk:                    binary.BigEndian.Uint16(payloadBytes[4:6]),
		Hcp:                   h,
		Exst:                  payloadBytes[8:16],
		Time:                  time.Now(),
		//IgnoreRetries:         true, // TODO: what does this field means? Why it is hardcoded in the kotlin code?
	}

	return hmgtr, nil
}

func (hgt *HmGetTransitionResponse) Encode() []byte {
	return []byte(hex.EncodeToString(hgt.data))
}

func (hgt *HmGetTransitionResponse) String() string {
	return fmt.Sprintf("HmGetTransitionResponse[StateInPercent: %v, DesiredStateInPerced: %v, Error: %v, AutoClose: %v, DriveTime: %v, Gk: %v, Hcp: %v, Exst: %v, Time: %v]", hgt.StateInPercent,
		hgt.DesiredStateInPercent,
		hgt.Error,
		hgt.AutoClose,
		hgt.DriveTime,
		hgt.Gk,
		hgt.Hcp,
		hgt.Exst,
		hgt.Time,
		//	hgt.IgnoreRetries,
	)
}

func (hgtr *HmGetTransitionResponse) Equal(o *HmGetTransitionResponse) bool {
	if hgtr.DesiredStateInPercent != o.DesiredStateInPercent {
		return false
	}

	if hgtr.Error != o.Error {
		return false
	}

	if hgtr.AutoClose != o.AutoClose {
		return false
	}

	if hgtr.DriveTime != o.DriveTime {
		return false
	}

	if hgtr.Gk != o.Gk {
		return false
	}

	if !bytes.Equal(hgtr.Exst, o.Exst) {
		return false
	}

	/*
		if hgtr.Time != o.Time {
			return false
		}
	*/

	if !hgtr.Hcp.Equal(o.Hcp) {
		return false
	}

	return true
}
