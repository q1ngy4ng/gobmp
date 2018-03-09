package dumputils

import (
	"fmt"
	"gobmp/bmpstorage"
	"net"
	"testing"
	"time"
)

func TestDumpSpeakerStatus(t *testing.T) {
	fmt.Println("=========Test DumpSpeakerStatus==========")
	s1 := bmpstorage.SpeakerStatus{
		BgpSpeakerAddress: "1.1.1.1",
		State:             true,
		Timestamp:         time.Now(),
		Localtimestamp:    time.Now(),
	}

	s2 := bmpstorage.SpeakerStatus{
		BgpSpeakerAddress: "2.2.2.2",
		State:             false,
		Timestamp:         time.Now(),
		Localtimestamp:    time.Now(),
	}

	var db1 map[int]*bmpstorage.SpeakerStatus
	db1 = make(map[int]*bmpstorage.SpeakerStatus)
	db1[1] = &s1
	db1[2] = &s2

	done := make(chan bool, 1)
	go DumpSpeakerStatus(done, true, db1)
	<-done
	go DumpSpeakerStatus(done, false, db1)
	<-done
}

func TestDumpPeerStatus(t *testing.T) {
	fmt.Println("=========Test DumpPeerStatus==========")
	p1 := bmpstorage.PeerStatus{
		State:          true,
		Timestamp:      time.Now(),
		Localtimestamp: time.Now(),
		UpdateCnt:      10,
	}

	p2 := bmpstorage.PeerStatus{
		State:          false,
		Timestamp:      time.Now(),
		Localtimestamp: time.Now(),
		UpdateCnt:      15,
	}

	var db map[string]*bmpstorage.PeerStatus
	db = make(map[string]*bmpstorage.PeerStatus)
	db["3.3.3.3"] = &p1
	db["4.4.4.4"] = &p2

	pdb := bmpstorage.PeerDB{
		Peer:db,
	}

	done := make(chan bool, 1)
	go DumpPeerDB(done, true, &pdb)
	<-done
	go DumpPeerDB(done, false, &pdb)
	<-done
}

func TestDumpPrefixDB(t *testing.T) {
	fmt.Println("=========Test DumpPrefixDB==========")
	var a [3]int
	a[0] = 1
	a[1] = 2
	a[2] = 3
	pathattr := bmpstorage.PathAttribute{
		NextHop:    net.IPv4(1, 1, 2, 2),
		Origin:     33,
		Med:        44,
		LocalPref:  55,
		AsPathLen:  3,
		AsPathData: []uint8{},
	}
	prefixattr := bmpstorage.PrefixAttr{
		PathAttribute:  &pathattr,
		Timestamp:      time.Now(),
		Localtimestamp: time.Now(),
		UpdateCnt:      15,
	}
	var prefixAttrmap map[string]*bmpstorage.PrefixAttr
	prefixAttrmap = make(map[string]*bmpstorage.PrefixAttr)
	prefixAttrmap["1.1.1.1"] = &prefixattr

	prefixDB := bmpstorage.PrefixDB{
		PrefixAttr: prefixAttrmap,
	}
	done := make(chan bool, 1)
	go DumpPrefixDB(done, true, &prefixDB)
	<-done
}

func TestDumpPeerPrefixDB(t *testing.T) {
	fmt.Println("=========Test DumpPeerPrefixDB==========")
	var a [3]int
	a[0] = 1
	a[1] = 2
	a[2] = 3
	pathattr := bmpstorage.PathAttribute{
		NextHop:    net.IPv4(1, 1, 2, 2),
		Origin:     33,
		Med:        44,
		LocalPref:  55,
		AsPathLen:  3,
		AsPathData: []uint8{},
	}
	prefixattr := bmpstorage.PrefixAttr{
		PathAttribute:  &pathattr,
		Timestamp:      time.Now(),
		Localtimestamp: time.Now(),
		UpdateCnt:      15,
	}
	var prefixAttrmap map[string]*bmpstorage.PrefixAttr
	prefixAttrmap = make(map[string]*bmpstorage.PrefixAttr)
	prefixAttrmap["1.1.1.1"] = &prefixattr
	prefixAttrmap["2.2.2.2"] = &prefixattr

	prefixDB := bmpstorage.PrefixDB{
		PrefixAttr: prefixAttrmap,
	}

	var prefixDBMap map[string]*bmpstorage.PrefixDB
	prefixDBMap = make(map[string]*bmpstorage.PrefixDB)
	prefixDBMap["3.3.3.3"] = &prefixDB
	prefixDBMap["4.4.4.4"] = &prefixDB

	peerPrefixDB := bmpstorage.PeerPrefixDB{
		PrefixDB: prefixDBMap,
	}

	done := make(chan bool, 1)
	go DumpPeerPrefixDB(done, true, &peerPrefixDB)

	<-done
}

/*
func TestAllDB(t *testing.T) {
	fmt.Println("=========Test DumpAllDB=============")
	pathattr := bmpstorage.PathAttribute{
		NextHop:    net.IPv4(1, 1, 2, 2),
		Origin:     33,
		Med:        44,
		LocalPref:  55,
		AsPathLen:  3,
		AsPathData: []uint8{},
	}
	prefixattr := bmpstorage.PrefixAttr{
		PathAttribute:  &pathattr,
		Timestamp:      time.Now(),
		Localtimestamp: time.Now(),
		UpdateCnt:      15,
	}
	var prefixAttrmap map[string]*bmpstorage.PrefixAttr
	prefixAttrmap = make(map[string]*bmpstorage.PrefixAttr)
	prefixAttrmap["1.1.1.1"] = &prefixattr

	prefixDB := bmpstorage.PrefixDB{
		PrefixAttr: prefixAttrmap,
	}
} */
