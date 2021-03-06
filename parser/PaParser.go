package parser

import (
	"encoding/binary"
	"gobmp/bmpstorage"
	"net"
)

// the following parse the input byte array and returns the path attribute
// struct along with the consumed length
func parsePathAttribute(inArray []uint8) (uint16, bmpstorage.PathAttribute) {
	var totalLen uint16 = 0
	var attrLen uint16
	var typeCode uint8
	var flag uint8
	var pa bmpstorage.PathAttribute

	var index uint16
	index = 0

	totalLen = binary.BigEndian.Uint16(inArray[index : index+2])
	index += 2

	for i := index; i < totalLen; i++ {
		flag = inArray[index]
		index += 1
		attrLen = uint16(inArray[index])
		index += 1
		if (flag & 0x10) != 0 { // extended length
			attrLen = (attrLen << 8) | uint16(inArray[index])
			index += 1
		}
		typeCode = inArray[index]
		index += 1

		if attrLen == 0 {
			panic("incorrect total length")
		}

		// at this point we can allocate a new struct
		pa := new(bmpstorage.PathAttribute)

		switch typeCode {
		case 1: // origin
			if attrLen != 4 {
				panic("incorrect origin length")
			}
			pa.Origin = binary.BigEndian.Uint32(inArray[index : index+4])
			index += 4
		case 2: // as path
			pa.AsPathLen = attrLen
			pa.AsPathData = make([]uint8, attrLen)
			copy(pa.AsPathData, inArray[index:index+attrLen])
			index += attrLen
		case 3: // next hop
			if attrLen != 4 {
				panic("incorrect nexthop length")
			}
			pa.NextHop = net.IPv4(inArray[index], inArray[index+1], inArray[index+2], inArray[index+3])
			index += 4
		case 4: // med
			if attrLen != 4 {
				panic("incorrect med length")
			}
			pa.Med = binary.BigEndian.Uint32(inArray[index : index+4])
			index += 4
		case 5: // local pref
			if attrLen != 4 {
				panic("incorrect local pref length")
			}
			pa.LocalPref = binary.BigEndian.Uint32(inArray[index : index+4])
			index += 4

		case 6: // atomic aggregate - we don't care for now
		case 7: // aggregator - we don't care
		}
	}
	return index, pa
}
