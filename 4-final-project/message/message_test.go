package message

import (
	"testing"

	"bitcoin_miner/hash"
)

func TestMarshallReqNoErr(t *testing.T) {
	msg := NewRequest("foo", 0, 1000)
	_, err := msg.ToJSON()
	if err != nil {
		t.Error(err)
	}
}

func TestMarshallResultNoErr(t *testing.T) {
	msg := NewResult(hash.Hash("foo", 500), 500, 0, 1000)
	_, err := msg.ToJSON()
	if err != nil {
		t.Error(err)
	}
}

func TestMarshallUnmarshallRequest(t *testing.T) {
	msg1 := NewRequest("foo", 0, 1000)
	jmsg, _ := msg1.ToJSON()
	msg2, err := FromJSON(jmsg)
	if err != nil {
		t.Error(err)
	}
	if *msg1 != *msg2 {
		t.Errorf("Request %v and %v are not the same.", msg1, msg2)
	}
}

func TestMarshallUnmarshallResult(t *testing.T) {
	msg1 := NewResult(hash.Hash("foo", 500), 500, 0, 1000)
	jmsg, _ := msg1.ToJSON()
	msg2, err := FromJSON(jmsg)
	if err != nil {
		t.Error(err)
	}
	if *msg1 != *msg2 {
		t.Errorf("Message %v and %v are not the same.", msg1, msg2)
	}
}
