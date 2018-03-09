package dumputils

import (
	"encoding/json"
	"fmt"
	"gobmp/bmpstorage"
	"log"
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
		fmt.Println("SpeakerId	", "Address	", "State	",
			"TimeStamp	", "LocalTimeStamp	")
		fmt.Println("---------	", "-------	", "-----	",
			"---------	", "--------------	")
		for key, value := range db {
			fmt.Println(key, "		",
				value.BgpSpeakerAddress, "	",
				value.State, "	",
				value.Timestamp, "	",
				value.Localtimestamp, "	")
		}
	}
	done <- true
}

func DumpAllSpeakerStatusDB(done chan bool, isJson bool,
	allSpeakerStatusdb map[int]*bmpstorage.SpeakerStatus) {

	if isJson {
		m, err := json.Marshal(allSpeakerStatusdb)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", m)
	} else {
		fmt.Println("All Speaker All Peer DB:")
	}
	done <- true
}

func DumpPeerDB(done chan bool, isJson bool,
	peerdb *bmpstorage.PeerDB) {
	if isJson {
		m, err := json.Marshal(peerdb.Peer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", m)
	} else {
		fmt.Println("Peer Status:")
		fmt.Println("PeerAddr	", "State	", "UpdateCnt	",
			"TimeStamp	", "LocalTimeStamp	")
		fmt.Println("---------	", "-----	", "--------		",
			"---------	", "--------------	")
		for key, value := range peerdb.Peer {
			fmt.Println(key, "		",
				value.State, "	",
				value.UpdateCnt, "	",
				value.Timestamp, "	",
				value.Localtimestamp, "	")
		}
	}
	done <- true
}

func DumpAllSpeakerPeerDB(done chan bool, isJson bool,
	allSpeakerPeerdb map[int]*bmpstorage.PeerDB) {

	if isJson {
		m, err := json.Marshal(allSpeakerPeerdb)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", m)
	} else {
		fmt.Println("All Speaker All Peer DB:")
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
	prefixdb *bmpstorage.PrefixDB) {

	if isJson {
		m, err := json.Marshal(prefixdb.PrefixAttr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", m)
	} else {
		fmt.Println("One Peer Prefixes:")
	}
	done <- true
}

func DumpPeerPrefixDB(done chan bool, isJson bool,
	peerPrefixdb *bmpstorage.PeerPrefixDB) {

	if isJson {
		m, err := json.Marshal(peerPrefixdb.PrefixDB)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", m)
	} else {
		fmt.Println("All Peer Prefixes:")
	}
	done <- true
}

func DumpAllSpeakerPeerPrefixDB(done chan bool, isJson bool,
	allSpeakerPrefixDB map[int]*bmpstorage.PeerPrefixDB) {

	if isJson {
		m, err := json.Marshal(allSpeakerPrefixDB)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", m)
	} else {
		fmt.Println("All Speaker All Peer Prefixes:")
	}
	done <- true
}

func DumpAllDB(done chan bool, isJson bool,
	allDB *bmpstorage.BmpDB) {

	d := make(chan bool, 1)
	DumpAllSpeakerStatusDB(d, true, allDB.Speaker)
	<-d
	DumpAllSpeakerPeerPrefixDB(d, true, allDB.PeerPrefixDB)
	<-d
	DumpAllSpeakerPeerDB(d, true, allDB.PeerDB)
}
