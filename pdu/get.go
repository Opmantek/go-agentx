// Copyright 2018 The agentx authors
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package pdu

// Get defines the pdu get packet.
type Get struct {
	Oids []ObjectIdentifier
}

const (
	REQ_OID_LIST_PAD_BYTES = 4
)

// Type returns the pdu packet type.
func (g *Get) Type() Type {
	return TypeGet
}

// MarshalBinary returns the pdu packet as a slice of bytes.
func (g *Get) MarshalBinary() ([]byte, error) {
	byteArr := []byte{}

	for _, oid := range g.Oids {
		oidByteArr, err := oid.MarshalBinary()
		if err != nil {
			return []byte{}, err
		}
		byteArr = append(byteArr, oidByteArr...)
		byteArr = append(byteArr, []byte{0, 0, 0, 0}...)
	}

	return byteArr, nil
}

// UnmarshalBinary sets the packet structure from the provided slice of bytes.
func (g *Get) UnmarshalBinary(data []byte) error {

	maxLength := len(data)
	digested := 0

	for digested < maxLength {

		oid := &ObjectIdentifier{}
		err := oid.UnmarshalBinary(data[digested:])
		if err != nil {
			return err
		}

		g.Oids = append(g.Oids, *oid)

		digested += oid.ByteSize() + REQ_OID_LIST_PAD_BYTES
	}

	return nil
}
