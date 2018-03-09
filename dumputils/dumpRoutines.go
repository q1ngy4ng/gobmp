package dumputils

import (
	"fmt"
	"encoding/json"
	"log"
	"gobmp/bmpstorage"
)

// Go routine to dump summary
func DumpSummary(done chan bool, m map[string]int) {
	fmt.Println("Show Summary:")
	fmt.Println("Prefix		", "NumRcvd")
	fmt.Println("------		", "---------------")
	for key, value := range m {
		fmt.Println(key, "	", value)
	}
	done <- true
}

func DumpSpeakerStatus(done chan bool, isJson bool,
					   db map[int]*bmpstorage.SpeakerStatus) {
	if isJson {
		m, err := json.Marshal(db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", m)
	} else {
		fmt.Println("Speaker Status:")
		fmt.Println("SpeakerId	","Address	","State	",
					"TimeStamp	","LocalTimeStamp	")
		fmt.Println("---------	","-------	","-----	",
					"---------	","--------------	")
		for key, value := range db {
			fmt.Println(key,"		",
						value.BgpSpeakerAddress,"	",
						value.State, "	",
						value.Timestamp, "	",
						value.Localtimestamp, "	")
		}
	}
	done <- true
}

func DumpPeerStatus(done chan bool, isJson bool,
					peerdb map[string]*bmpstorage.PeerStatus) {
    if isJson {
        m, err := json.Marshal(peerdb)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s\n", m)
    } else {
		fmt.Println("Peer Status:")
		fmt.Println("PeerAddr	","State	","UpdateCnt	",
					"TimeStamp	","LocalTimeStamp	")
		fmt.Println("---------	","-----	","--------		",
					"---------	","--------------	")
		for key, value := range peerdb {
			fmt.Println(key,"		",
						value.State,"	",
						value.UpdateCnt, "	",
						value.Timestamp, "	",
						value.Localtimestamp, "	")
		}
	}
	done <- true
}

/* For reference
type PrefixDB struct {
    // key: prefix
    PrefixAttr map[string]*PrefixAttr
}
type PeerPrefixDB struct {
    // key: peer_address
    PrefixDB map[string]*PrefixDB
}
*/

func DumpPrefixDB(done chan bool, isJson bool,
				  peerdb map[string]*bmpstorage.PrefixDB) {

    if isJson {
        m, err := json.Marshal(peerdb)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s\n", m)
    } else {
		fmt.Println("Prefixes:")
	}
	done <- true
}

