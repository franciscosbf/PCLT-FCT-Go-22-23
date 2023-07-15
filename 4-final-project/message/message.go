// DO NOT MODIFY THIS FILE!

package message

import (
	"encoding/json"
	"fmt"
)

type MsgType int

const (
	Request MsgType = iota
	Result
)

// Message represents a message that can be sent between the
// server and a client of the bitcoin project. The messages
// effectively are of two types: a client sends Requests to
// to the server, generally of the form:
// Message{Type:Request,Data:msg,Lower:0,Upper:N,Hash:0,Nonce:0}
// hoping to receive back a message of the form:
// Message{Type:Result,Data:nil,Lower:0,Upper:N,Hash:value,Nonce:nonce}
// where Nonce corresponds to a value between Lower and Upper such that
// the hash of the original Data ++ Nonce is Hash and Hash is the maximum
// such hash.
type Message struct {
	Type         MsgType
	Data         string
	Lower, Upper uint64
	Hash, Nonce  uint64
}

// ToJson() method serializes a Message into a JSON byte slice, ready
// to be sent over the wire.
func (msg *Message) ToJSON() ([]byte, error) {
	return json.Marshal(msg)
}

// FromJson() function takes a byte slice containing a JSON representation of
// a Message and unmarshals it back into a Go value.
func FromJSON(data []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return &m, err
}

// NewResult(hash,nonce,low,up) creates a result message for the mining problem
// whose solution is the hash nonce pair, with lower bound low and upper bound up.
func NewResult(hash, nonce, low, up uint64) *Message {
	return &Message{
		Type:  Result,
		Hash:  hash,
		Nonce: nonce,
		Lower: low,
		Upper: up,
	}
}

// NewRequest(data,low,up) creates a request message for data and interval
// given by [low,up].
func NewRequest(data string, low uint64, up uint64) *Message {
	return &Message{
		Type:  Request,
		Data:  data,
		Lower: low,
		Upper: up,
	}
}

func (msg *Message) String() string {
	if msg.Type == Request {
		return fmt.Sprintf("MinerMsg[Request \"%s\" %d %d]", msg.Data, msg.Lower, msg.Upper)
	} else {
		return fmt.Sprintf("MinerMsg[Result %d %d %d %d]", msg.Hash, msg.Nonce, msg.Lower, msg.Upper)
	}
}
