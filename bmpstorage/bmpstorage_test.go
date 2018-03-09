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
	db := new(BmpDB)
	now := time.Now()
	db.UpdatePeer(123, "2.2.2.2", true, now)
	if db.PeerDB[123].Peer["2.2.2.2"].State != true {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].State != true")
		t.Fail()
	}
	if db.PeerDB[123].Peer["2.2.2.2"].UpdateCnt != 1 {
		t.Log("db.PeerDB[123].Peer['2.2.2.2'].UpdateCnt != 1")
		t.Fail()
	}

	db.UpdateSpeaker(123, "2.2.2.2", true, now)
	if db.Speaker[123].BgpSpeakerAddress != "2.2.2.2" {
		t.Log("db.Speaker[123].BgpSpeakerAddress != '2.2.2.2'")
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

	fmt.Printf("TestBmpstorage Ok0 \n")
}
