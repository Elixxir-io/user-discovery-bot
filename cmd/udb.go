////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// User Discovery Bot main functions (Start Bot and register)
// This file covers all of the glue code necessary to run the bot. All of the
// interesting code is in the udb module.

package cmd

import (
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"os"
	"strings"
	"time"
)

// RateLimit for messages in ms (100 = 10 msg per second)
const RateLimit = 100

// UDBSessionFileName used by UDB
var UDBSessionFileName string

var clientObj *api.Client

// StartBot starts the user discovery bot:
//  - Set up global variables
//  - Log into the server
//  - Start the main loop
func StartBot(sess string, def *ndf.NetworkDefinition) {
	udb.Log.DEBUG.Printf("Starting User Discovery Bot...")

	UDBSessionFileName = sess

	// Initialize the client
	regCode := udb.UDB_USERID.RegistrationCode()
	Init(UDBSessionFileName, regCode, def)

	udb.Log.INFO.Printf("Logging in")

	// Log into the server with a blank password
	_, err := clientObj.Login("")

	if err != nil {
		udb.Log.FATAL.Panicf("Could not login: %s", err)
	}

	// get the newest message ID on the reception gateway to stop the UDB from
	// replaying old messages in the event of a redeploy where the session file
	// is lost
	lastMessageID := getLatestMessageID()
	clientObj.GetSession().SetLastMessageID(lastMessageID)

	// Register the listeners with the user discovery bot
	udb.RegisterListeners(clientObj)

	udb.Log.INFO.Printf("Starting UDB")

	// starting the reception thread
	err = clientObj.StartMessageReceiver()
	if err != nil {
		udb.Log.FATAL.Panicf("Could not start message recievers:  %+v", err)
	}

	// Block forever as a keepalive
	select {}
}

// Initialize a session using the given session file and other info
func Init(sessionFile string, regCode string, def *ndf.NetworkDefinition) *id.User {
	userID := udb.UDB_USERID

	// We only register when the session file does not exist
	// FIXME: this is super weird -- why have to check for a file,
	// then init that file, then register optionally based on that check?
	_, err := os.Stat(sessionFile)
	// Get new client. Setting storage to nil internally creates a
	// default storage
	var initErr error

	dummyConnectionStatusHandler := func(status uint32, timeout int) {
		globals.Log.INFO.Printf("Network status: %+v, %+v", status, timeout)
	}
	secondarySessionFile := sessionFile + "-2"
	clientObj, initErr = api.NewClient(nil, sessionFile, secondarySessionFile, def, dummyConnectionStatusHandler)
	if initErr != nil {
		udb.Log.FATAL.Panicf("Could not initialize: %v", initErr)
	}

	if noTLS {
		clientObj.DisableTLS()
	}

	// API Settings (hard coded)
	clientObj.DisableBlockingTransmission() // Deprecated
	// Up to 10 messages per second
	clientObj.SetRateLimiting(uint32(RateLimit))

	// connect udb to gateways
	for {
		err = clientObj.Connect()
		if err == nil {
			break
		}
		time.Sleep(10 * time.Second)
		udb.Log.ERROR.Printf("UDB could not connect to gateways, "+
			"reconnecting: %+v", err)
	}

	// SB: Trying to always register.
	// I think it's needed for some things to work correctly.
	// Need a more accurate descriptor of what the method actually does than
	// Register, or to remove the things that aren't actually used for
	// registration.
	userID, err = clientObj.RegisterWithPermissioning(true, regCode, "",
		"", "", nil)
	if err != nil {
		udb.Log.FATAL.Panicf("Could not register with Permissioning: %v", err)
	}

	return userID
}

// getLatestMessageID gets the newest message ID on the reception gateway, used
// to stop the UDB from replaying old messages in the event of a redeploy where
// the session file is lost
func getLatestMessageID() string {
	//get the newest message id to
	msg := &mixmessages.ClientRequest{
		UserID:        udb.UDB_USERID.Bytes(),
		LastMessageID: "",
	}

	receiveGateway := id.NewNodeFromBytes(clientObj.GetNDF().Nodes[len(clientObj.GetNDF().Gateways)-1].ID).NewGateway()

	var idList *mixmessages.IDList

	for {
		var err error
		idList, err = clientObj.GetCommManager().Comms.SendCheckMessages(receiveGateway, msg)

		if err != nil {
			globals.Log.WARN.Printf("Failed to get the latest message "+
				"IDs from the reception gateway: %s", err.Error())
			if strings.Contains(err.Error(),
				"Could not find any message IDs for this user") {
				break
			}
		} else {
			break
		}

		time.Sleep(2 * time.Second)
	}

	lastMessage := ""

	if idList != nil && idList.IDs != nil && len(idList.IDs) != 0 {
		lastMessage = idList.IDs[len(idList.IDs)-1]
	}

	globals.Log.INFO.Printf("Discarding messages before ID `%s`", lastMessage)

	return lastMessage
}
