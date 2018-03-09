package dumputils 

import (
	"testing"
	"time"
	"gobmp/bmpstorage"
	"fmt"
)

func TestDumpSpeakerStatus(t *testing.T) {  
	fmt.Println("=========Test DumpSpeakerStatus==========")
	s1 := bmpstorage.SpeakerStatus{
        BgpSpeakerAddress:"1.1.1.1",
        State:true,
		Timestamp:time.Now(),
		Localtimestamp:time.Now(),
    }

    s2 := bmpstorage.SpeakerStatus{
        BgpSpeakerAddress:"2.2.2.2",
        State:false,
		Timestamp:time.Now(),
		Localtimestamp:time.Now(),
    }

    var db1 map[int]*bmpstorage.SpeakerStatus
    db1 = make(map[int]*bmpstorage.SpeakerStatus)
    db1[1] = &s1
    db1[2] = &s2

	done := make(chan bool, 1)
	go DumpSpeakerStatus(done, true, db1)
	<- done
	go DumpSpeakerStatus(done, false, db1)
	<- done
}

func TestDumpPeerStatus(t *testing.T) {
	fmt.Println("=========Test DumpPeerStatus==========")
	p1 := bmpstorage.PeerStatus{
    	State:true,
    	Timestamp: time.Now(),
    	Localtimestamp: time.Now(),
    	UpdateCnt:10,
	}

	p2 := bmpstorage.PeerStatus{
    	State:false,
    	Timestamp:time.Now(),
    	Localtimestamp:time.Now(),
    	UpdateCnt:15,
	}

	var db map[string]*bmpstorage.PeerStatus
	db = make(map[string]*bmpstorage.PeerStatus)
	db["3.3.3.3"] = &p1
	db["4.4.4.4"] = &p2

	done := make(chan bool, 1)
	go DumpPeerStatus(done, true, db)
    <- done
	go DumpPeerStatus(done, false, db)
    <- done
}
