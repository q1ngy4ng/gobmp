package main

import "fmt"

type AsPathAttr struct {
   len U32
   data U8[int]
}

type PathAttribute struct {
   nextHop U32  // let’s just do v4 for now (v6 would be in a different TLV anyway)
   origin U8
   //pathFlags U8
   //originatorId U32 // this is not interesting and I will just remove it
   //aspType U32  ← just for the record, we don’t seem to need this in bgpSmash!
   med U32
   localPref U32
   asPathLen U32
   asPathData *U8
   //commList CommList
   //extCommListId : ExtCommListId;
}
     
// the following parse the input byte array and returns the path attribute 
// struct along with the consumed length
func parsePathAttribute(inArray []byte) {
	var totalLen int = 0
        var attrLen int
        var typeCode int
        var flag int
 
	var index int = 0

        totalLen = (inArray[index++] << 8) | inArray[index]

        for i := index; i < totalLen; i++ {
		flag = inArray[index++]
                attrLen = inArray[index++]
                if (flag & 0x10) {  // extended length
                     attrLen = (attrLen << 8) | inArray[index++]
                }
                typeCode = inArray[index++]
		
		if (attrLen == 0) {
			panic("incorrect total length")
		} 

		// at this point we can allocate a new struct 
                pa := new(PathAttribute)

        	switch typeCode {
		case 1:  // origin
			if (attrLen != 4) {
				panic("incorrect origin length") 
			}
			pa.origin = (inArray[index++] << 24) | (inArray[index++] << 16) | (inArray[index++] << 8) |(inArray[index++]) 
                case 2:  // as path
 			pa.asPathLen = attrLen
                        pa.asPathData = make([]byte, attrLen)
			copy(pa.asPathData, inArray[index:index+attrLen])
			index += attrLen                       
                case 3:  // next hop
               		if (attrLen != 4) {
				panic("incorrect nexthop length")
			}
			pa.nexthop = (inArray[index++] << 24) | (inArray[index++] << 16) | (inArray[index++] << 8) |(inArray[index++])
                case 4: // med
			if (attrLen != 4) {
                                panic("incorrect med length")
                        }
                        pa.med = (inArray[index++] << 24) | (inArray[index++] << 16) | (inArray[index++] << 8) |(inArray[index++])
                case 5: // local pref
                        if (attrLen != 4) {
                                panic("incorrect local pref length")
                        }
                        pa.localPref = (inArray[index++] << 24) | (inArray[index++] << 16) | (inArray[index++] << 8) |(inArray[index++])
	
                case 6: // atomic aggregate - we don't care for now
                case 7: // aggregator - we don't care
                } 
 	} 
        return index, pa      
}

// the following dumps the as path into a human readable string
func dumpAsPath(aspath AsPathAttr) {
} 
