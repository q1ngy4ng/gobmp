// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package bmpstorage

import "fmt"
import "time"

type PathAttr struct {
	attr string
}
type PrefixAttr struct {
	PathAttr       *PathAttr
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
	Speaker  map[int]*SpeakerStatus
	PeerDB   map[int]*PeerDB
	PrefixDB map[int]*PeerPrefixDB
}

func (db *BmpDB) UpdateRoute(speakerId int, peerAddress string,
	prefix string, pathAttr *PathAttr, timestamp time.Time) {
	fmt.Printf("RouteUpdate %d %s %s %p %s %s\n", speakerId, peerAddress,
		prefix, pathAttr,
		timestamp.Format(time.RFC850),
		time.Now().Format(time.RFC850))

}

func (db *BmpDB) UpdateSpeaker(speakerId int, speakerAddress string,
	initialize bool, timestamp time.Time) {
	fmt.Printf("UpdateSpeaker %d %s %v %s %s\n",
		speakerId, speakerAddress, initialize,
		timestamp.Format(time.RFC850),
		time.Now().Format(time.RFC850))

}

func (db *BmpDB) UpdatePeer(speakerId int, peerAddress string,
	up bool, timestamp time.Time) {
	fmt.Printf("UpdatePeer %d %s %v %s %s\n",
		speakerId, peerAddress, up,
		timestamp.Format(time.RFC850),
		time.Now().Format(time.RFC850))
}
