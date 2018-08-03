////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Receive and parse user discovery bot messages, and run the appropriate
// command
package udb

import (
	"github.com/mattn/go-shellwords"
	"gitlab.com/privategrity/client/parse"
	"gitlab.com/privategrity/client/switchboard"
	"gitlab.com/privategrity/client/user"
	"gitlab.com/privategrity/crypto/cyclic"
)

type SearchListener struct{}
type RegisterListener struct{}
type PushKeyListener struct{}
type GetKeyListener struct{}

// Register the UDB listeners
func init() {
	switchboard.Listeners.Register(user.ID(0), parse.Type_UDB_SEARCH, SearchListener{})
	switchboard.Listeners.Register(user.ID(0), parse.Type_UDB_REGISTER, RegisterListener{})
	switchboard.Listeners.Register(user.ID(0), parse.Type_UDB_PUSH_KEY, PushKeyListener{})
	switchboard.Listeners.Register(user.ID(0), parse.Type_UDB_GET_KEY, GetKeyListener{})
}

// Listen for Search Messages
func (s SearchListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := cyclic.NewIntFromBytes(message.GetSender()).Uint64()
	args, err := shellwords.Parse(message.GetPayload())
	if err != nil {

	}
	Search(sender, args[1:])
}

// Listen for Register Messages
func (s RegisterListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := cyclic.NewIntFromBytes(message.GetSender()).Uint64()
	args, err := shellwords.Parse(message.GetPayload())
	if err != nil {

	}
	Register(sender, args[1:])
}

// Listen for PushKey Messages
func (s PushKeyListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := cyclic.NewIntFromBytes(message.GetSender()).Uint64()
	args, err := shellwords.Parse(message.GetPayload())
	if err != nil {

	}
	PushKey(sender, args[1:])
}

// Listen for GetKey Messages
func (s GetKeyListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := cyclic.NewIntFromBytes(message.GetSender()).Uint64()
	args, err := shellwords.Parse(message.GetPayload())
	if err != nil {

	}
	GetKey(sender, args[1:])
}
