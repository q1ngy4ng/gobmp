// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package bmpstorage

import (
	"fmt"
	"sync"
	"time"
	"net"
)

type PathAttribute struct {
   NextHop net.IP  // let’s just do v4 for now (v6 would be in a different TLV)
   Origin uint32
   //PathFlags uint8
   //OriginatorId uint32 // this is not interesting and I will just remove it
   //AspType uint32  ← just for the record, we don’t seem to need this in bgpSmaash!
   Med uint32
   LocalPref uint32
   AsPathLen uint16
   AsPathData []uint8
   //CommList CommList
   //ExtCommListId : ExtCommListId;
}

type PrefixAttr struct {
	PathAttribute       *PathAttribute
	Timestamp      time.Time
	Localtimestamp time.Time
	UpdateCnt      uint32
}
type PrefixDB struct {
	// key: prefix
	PrefixAttr map[string]*PrefixAttr
}
type PeerPrefixDB struct {
	// key: peer_address
	PrefixDB map[string]*PrefixDB
}
type SpeakerStatus struct {
	BgpSpeakerAddress string
	State             bool
	Timestamp         time.Time
	Localtimestamp    time.Time
	UpdateCnt         uint32
}
type PeerStatus struct {
	State          bool
	Timestamp      time.Time
	Localtimestamp time.Time
	UpdateCnt      uint32
}
type PeerDB struct {
	// key : peer_address
	Peer map[string]*PeerStatus
}
type BmpDB struct {
	// key : bgp_speaker_id
	Speaker      map[int]*SpeakerStatus
	PeerDB       map[int]*PeerDB
	PeerPrefixDB map[int]*PeerPrefixDB
	mutex        sync.Mutex
}

func (db *BmpDB) UpdateRoute(speakerId int, peerAddress string,
	prefix string, pathAttr *PathAttribute, timestamp time.Time) {
	fmt.Printf("UpdateRoute %d %s %s %p %s %s\n", speakerId, peerAddress,
		prefix, pathAttr,
		timestamp.Format(time.RFC850),
		time.Now().Format(time.RFC850))
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if db.PeerPrefixDB == nil {
		db.PeerPrefixDB = map[int]*PeerPrefixDB{}
	}
	if _, ok := db.PeerPrefixDB[speakerId]; !ok {
		db.PeerPrefixDB[speakerId] = new(PeerPrefixDB)
	}
	peerPrefixDB := db.PeerPrefixDB[speakerId]
	if peerPrefixDB.PrefixDB == nil {
		peerPrefixDB.PrefixDB = map[string]*PrefixDB{}
	}
	if _, ok := peerPrefixDB.PrefixDB[peerAddress]; !ok {
		peerPrefixDB.PrefixDB[peerAddress] = new(PrefixDB)
	}
	prefixDB := peerPrefixDB.PrefixDB[peerAddress]
	if prefixDB.PrefixAttr == nil {
		prefixDB.PrefixAttr = map[string]*PrefixAttr{}
	}
	if _, ok := prefixDB.PrefixAttr[prefix]; !ok {
		prefixDB.PrefixAttr[prefix] = new(PrefixAttr)
		prefixDB.PrefixAttr[prefix].UpdateCnt = 0
	}
	prefixDB.PrefixAttr[prefix].UpdateCnt += 1
	prefixDB.PrefixAttr[prefix].Timestamp = timestamp
	prefixDB.PrefixAttr[prefix].Localtimestamp = time.Now()
	prefixDB.PrefixAttr[prefix].PathAttribute = pathAttr
}

func (db *BmpDB) UpdateSpeaker(speakerId int, speakerAddress string,
	initialize bool, timestamp time.Time) {
	fmt.Printf("UpdateSpeaker %d %s %v %s %s\n",
		speakerId, speakerAddress, initialize,
		timestamp.Format(time.RFC850),
		time.Now().Format(time.RFC850))
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if db.Speaker == nil {
		db.Speaker = map[int]*SpeakerStatus{}
	}
	if _, ok := db.Speaker[speakerId]; !ok {
		db.Speaker[speakerId] = new(SpeakerStatus)
		db.Speaker[speakerId].UpdateCnt = 0
	}
	db.Speaker[speakerId].UpdateCnt += 1
	db.Speaker[speakerId].BgpSpeakerAddress = speakerAddress
	db.Speaker[speakerId].State = initialize
	db.Speaker[speakerId].Timestamp = timestamp
	db.Speaker[speakerId].Localtimestamp = time.Now()
}

func (db *BmpDB) UpdatePeer(speakerId int, peerAddress string,
	up bool, timestamp time.Time) {
	fmt.Printf("UpdatePeer %d %s %v %s %s\n",
		speakerId, peerAddress, up,
		timestamp.Format(time.RFC850),
		time.Now().Format(time.RFC850))
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if db.PeerDB == nil {
		db.PeerDB = map[int]*PeerDB{}
	}
	if _, ok := db.PeerDB[speakerId]; !ok {
		db.PeerDB[speakerId] = new(PeerDB)
	}
	peerDB := db.PeerDB[speakerId]
	if peerDB.Peer == nil {
		peerDB.Peer = map[string]*PeerStatus{}
	}
	if _, ok := peerDB.Peer[peerAddress]; !ok {

		peerDB.Peer[peerAddress] = new(PeerStatus)
		peerDB.Peer[peerAddress].UpdateCnt = 0
	}
	peerDB.Peer[peerAddress] = new(PeerStatus)
	peerDB.Peer[peerAddress].State = up
	peerDB.Peer[peerAddress].Timestamp = timestamp
	peerDB.Peer[peerAddress].Localtimestamp = time.Now()
	peerDB.Peer[peerAddress].UpdateCnt += 1
}

var (
	bmpDB *BmpDB
	once  sync.Once
)

func GetBmpDB() *BmpDB {
	once.Do(func() {
		bmpDB = new(BmpDB)
	})
	return bmpDB
}
