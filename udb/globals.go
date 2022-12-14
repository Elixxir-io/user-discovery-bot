////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/primitives/id"
	"io/ioutil"
	"log"
	"os"
)

// The User Discovery Bot's user ID and registration code
// (this is global in cMix systems)
var UDB_USERID *id.User = id.NewUserFromUints(&[4]uint64{0, 0, 0, 3})

var Log = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, os.Stdout,
	ioutil.Discard, "CLIENT", log.Ldate|log.Ltime)

func init() {
	globals.Log = Log
}
