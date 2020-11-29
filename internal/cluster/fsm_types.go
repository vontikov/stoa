package cluster

import (
	"encoding/binary"
)

// entry is an data object stored in the FSM
type entry struct {
	Value     []byte // Value is the object's value
	TTLMillis int64  // TTLMillis is the object's TTL
}

type entryPtr = *entry

// MarshalBinary implements encoding.BinaryMarshaler.MarshalBinary
// (https://golang.org/pkg/encoding/#BinaryMarshaler).
func (e *entry) MarshalBinary() ([]byte, error) {
	data := make([]byte, 4+len(e.Value)+8)
	binary.LittleEndian.PutUint32(data, uint32(len(e.Value)))
	n := copy(data[4:], e.Value)
	binary.LittleEndian.PutUint64(data[4+n:], uint64(e.TTLMillis))
	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.UnmarshalBinary
// (https://golang.org/pkg/encoding/#BinaryUnmarshaler).
func (e *entry) UnmarshalBinary(data []byte) error {
	e.Value = nil
	sz := int(binary.LittleEndian.Uint32(data))
	if sz > 0 {
		e.Value = append(e.Value, data[4:4+sz]...)
	}
	e.TTLMillis = int64(binary.LittleEndian.Uint64(data[4+sz:]))
	return nil
}
