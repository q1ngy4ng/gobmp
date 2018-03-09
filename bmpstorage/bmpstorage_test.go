// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package bmpstorage

import (
	"fmt"
	"testing"
	"time"
)

func TestBmpstorage(t *testing.T) {
	fmt.Printf("TestBmpstorage Ok \n")
	db := GetBmpDB()
	now := time.Now()

	db.UpdateSpeaker(123, "2.2.2.21", true, now)
	if db.Speaker[123].BgpSpeakerAddress != "2.2.2.21" {
		t.Log("db.Speaker[123].BgpSpeakerAddress != '2.2.2.21'")
		t.Fail()
	}
	if db.Speaker[123].State != true {
		t.Log("db.Speaker[123].State != true")
		t.Fail()
	}
	if db.Speaker[123].UpdateCnt != 1 {
		t.Log("db.Speaker[123].UpdateCnt != 1")
		t.Fail()
	}
	if db.Speaker[123].Timestamp != now {
		t.Log("db.Speaker[123].Timestamp != now")
		t.Fail()
	}

	db.UpdatePeer(123, "2.2.2.2", true, now)
	if db.PeerDB[123].Peer["2.2.2.2"].State != true {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].State != true")
		t.Fail()
	}
	if db.PeerDB[123].Peer["2.2.2.2"].UpdateCnt != 1 {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].UpdateCnt != 1")
		t.Fail()
	}

	var pathAttr PathAttr
	db.UpdateRoute(123, "2.2.2.2", "3.3.3.0/24", &pathAttr, now)
	prefixAttr := db.PeerPrefixDB[123].PrefixDB["2.2.2.2"].PrefixAttr["3.3.3.0/24"]
	if prefixAttr.PathAttr != &pathAttr {
		t.Log("prefixAttr.PathAttr != &pathAttr")
		t.Fail()
	}
	if prefixAttr.Timestamp != now {
		t.Log("prefixAttr.Timestamp != now")
		t.Fail()
	}
	if prefixAttr.UpdateCnt != 1 {
		t.Log("prefixAttr.UpdateCnt != 1")
		t.Fail()
	}
	if prefixAttr.Localtimestamp.Before(now) {
		t.Log("prefixAttr.Localtimestamp.Before( now )")
		t.Fail()
	}
	db.UpdateRoute(123, "2.2.2.2", "3.3.30.0/24", &pathAttr, now)
	prefixAttr = db.PeerPrefixDB[123].PrefixDB["2.2.2.2"].PrefixAttr["3.3.3.0/24"]
	if prefixAttr.PathAttr != &pathAttr {
		t.Log("prefixAttr.PathAttr != &pathAttr")
		t.Fail()
	}
	prefixAttr = db.PeerPrefixDB[123].PrefixDB["2.2.2.2"].PrefixAttr["3.3.30.0/24"]
	if prefixAttr.PathAttr != &pathAttr {
		t.Log("prefixAttr.PathAttr != &pathAttr")
		t.Fail()
	}
	if prefixAttr.Timestamp != now {
		t.Log("prefixAttr.Timestamp != now")
		t.Fail()
	}
	if prefixAttr.UpdateCnt != 1 {
		t.Log("prefixAttr.UpdateCnt != 1")
		t.Fail()
	}
	if prefixAttr.Localtimestamp.Before(now) {
		t.Log("prefixAttr.Localtimestamp.Before( now )")
		t.Fail()
	}

	db = GetBmpDB()
	prefixAttr = db.PeerPrefixDB[123].PrefixDB["2.2.2.2"].PrefixAttr["3.3.3.0/24"]
	if prefixAttr.PathAttr != &pathAttr {
		t.Log("prefixAttr.PathAttr != &pathAttr")
		t.Fail()
	}

	var modifiedPathAttr PathAttr
	now = time.Now()
	db.UpdateRoute(123, "2.2.2.2", "3.3.3.0/24", &modifiedPathAttr, now)
	prefixAttr = db.PeerPrefixDB[123].PrefixDB["2.2.2.2"].PrefixAttr["3.3.3.0/24"]
	if prefixAttr.PathAttr != &modifiedPathAttr {
		t.Log("prefixAttr.PathAttr != &modifiedPathAttr")
		t.Fail()
	}
	if prefixAttr.Timestamp != now {
		t.Log("prefixAttr.Timestamp != now")
		t.Fail()
	}
	if prefixAttr.UpdateCnt != 2 {
		t.Log("prefixAttr.UpdateCnt != 1")
		t.Fail()
	}
	if prefixAttr.Localtimestamp.Before(now) {
		t.Log("prefixAttr.Localtimestamp.Before( now )")
		t.Fail()
	}

	if db.PeerDB[123].Peer["2.2.2.2"].UpdateCnt != 1 {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].UpdateCnt != 2",
			db.PeerDB[123].Peer["2.2.2.2"].UpdateCnt)
		t.Fail()
	}
	if db.PeerDB[123].Peer["2.2.2.2"].State != true {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].State != true")
		t.Fail()
	}
	db.UpdatePeer(123, "2.2.2.2", false, now)
	if db.PeerDB[123].Peer["2.2.2.2"].State != false {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].State != false")
		t.Fail()
	}
	if db.PeerDB[123].Peer["2.2.2.2"].UpdateCnt != 2 {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].UpdateCnt != 2")
		t.Fail()
	}

	prefixAttr = db.PeerPrefixDB[123].PrefixDB["2.2.2.2"].PrefixAttr["3.3.3.0/24"]
	if prefixAttr.PathAttr != &modifiedPathAttr {
		t.Log("prefixAttr.PathAttr != &modifiedPathAttr")
		t.Fail()
	}

	db.UpdatePeer(123, "2.2.2.2", true, now)
	if _, ok := db.PeerPrefixDB[123].PrefixDB["2.2.2.2"]; ok {
		t.Log("PrefixDB is not cleaned up when the peer 2.2.2.2 is up")
		t.Fail()
	}
	if db.PeerDB[123].Peer["2.2.2.2"].UpdateCnt != 3 {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].UpdateCnt != 3")
		t.Fail()
	}

	db.UpdateSpeaker(123, "2.2.2.21", false, now)
	if db.Speaker[123].UpdateCnt != 2 {
		t.Log("db.Speaker[123].UpdateCnt != 2")
		t.Fail()
	}
	if db.PeerDB[123].Peer["2.2.2.2"].UpdateCnt != 3 {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].UpdateCnt != 3")
		t.Fail()
	}
	db.UpdateSpeaker(123, "2.2.2.21", true, now)
	if db.Speaker[123].UpdateCnt != 3 {
		t.Log("db.Speaker[123].UpdateCnt != 3")
		t.Fail()
	}
	if _, ok := db.PeerDB[123]; ok {
		t.Log("db.PeerDB[123].Peer is not clean up when the speaker is up")
		t.Fail()
	}
	if _, ok := db.PeerPrefixDB[123]; ok {
		t.Log("db.PeerPrefixDB[123] is not clean up when the speaker is up")
		t.Fail()
	}
	fmt.Printf("TestBmpstorage Ok0 \n")
}
