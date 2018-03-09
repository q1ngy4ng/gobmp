package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"gobmp/parser"
)

// the following dumps the as path into a human readable string
//func dumpAsPath(aspath AsPathAttr) {
//}
func dumpAsPath(len uint16, data []uint8) {
	var segType uint8
	var segLen uint8
	var as uint32

        for i:= uint16(0); i < len; {
		segType = data[i]
		i += 1
		segLen = data[i]
		i += 1
		var ob, cb string
		switch segType {
		case 1: // as set
			ob = "("
			cb = ")"
		case 2: // as sequence
			ob = ""
			cb = ""
		case 3: // confed set
			ob = "["
			cb = "]"
		case 4: // confed seq
			ob = "{"
			cb = "}"		
		}
		fmt.Println(ob)
		for j := segLen; j > 0; j-- {
			as = binary.BigEndian.Uint32(data[i:i+4])
			i += 4
			fmt.Printf("%d", as)
			if (j > 1) {
				fmt.Println(" ")
			}
		}
		fmt.Println(cb) 
	}
}

func dumpPathAttribute(pa PathAttribute, indent string) {
	fmt.Printf("%s", indent)

	fmt.Printf("Path Attribute:\n")

        fmt.Printf("\t%s", indent)
        fmt.Printf("Origin: ")
        switch pa.origin {
	case 0: // igp
		fmt.Printf("i")
	case 1: // egp
		fmt.Printf("e")
	case 2: // unknown
		fmt.Printf("?")
	}
	fmt.Printf("\n")

	fmt.Printf("\t%s", indent)
	fmt.Printf("Next Hop: ")
	fmt.Println(pa.nextHop)
	fmt.Printf("\n")

	if (pa.med > 0) {
		fmt.Printf("\t%s", indent)
		fmt.Printf("Multi Exit Discriminator: ")
		fmt.Printf("%d", pa.med)
		fmt.Printf("\n")
	}

	if (pa.localPref > 0) {  // this is probably wrong as local pref could be 0
		fmt.Printf("\t%s", indent)
		fmt.Printf("Local Preference: ")
		fmt.Printf("%d", pa.localPref)
		fmt.Printf("\n")
	}

	fmt.Printf("\t%s", indent)
	fmt.Printf("As Path: ")
	dumpAsPath(pa.asPathLen, pa.asPathData)
	fmt.Printf("\n")	
	

} 
