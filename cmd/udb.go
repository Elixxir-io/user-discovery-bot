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
	"crypto/rand"
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"io/ioutil"
	"os"
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

	// Use RAM storage for now
	udb.DataStore = storage.NewRamStorage()

	UDBSessionFileName = sess

	// Initialize the client
	regCode := udb.UDB_USERID.RegistrationCode()
	userID := Init(UDBSessionFileName, regCode, def)

	// Get the default parameters and generate a public key from it
	dsaParams := signature.GetDefaultDSAParams()
	publicKey := dsaParams.PrivateKeyGen(rand.Reader).PublicKeyGen()

	// Save DSA public key and user ID to JSON file
	outputDsaPubKeyToJSON(publicKey, udb.UDB_USERID, ".elixxir",
		"udb_info.json")

	// API Settings (hard coded)
	clientObj.DisableBlockingTransmission() // Deprecated
	// Up to 10 messages per second
	clientObj.SetRateLimiting(uint32(RateLimit))

	// Log into the server
	Login(userID)

	// Register the listeners with the user discovery bot
	udb.RegisterListeners(clientObj)

	// TEMPORARILY try starting the reception thread here instead-it seems to
	// not be starting?
	//go io.Messaging.MessageReceiver(time.Second)

	// Block forever as a keepalive
	quit := make(chan bool)
	<-quit
}

// Init -ialize a session using the given session file and other info
func Init(sessionFile string, regCode string, def *ndf.NetworkDefinition) *id.User {
	userID := udb.UDB_USERID

	// We only register when the session file does not exist
	// FIXME: this is super weird -- why have to check for a file,
	// then init that file, then register optionally based on that check?
	_, err := os.Stat(sessionFile)
	// Get new client. Setting storage to nil internally creates a
	// default storage
	var initErr error
	clientObj, initErr = api.NewClient(nil, sessionFile, def)
	if initErr != nil {
		udb.Log.FATAL.Panicf("Could not initialize: %v", initErr)
	}

	//connect udb to gateways
	err = clientObj.Connect()
	if err != nil {
		udb.Log.FATAL.Printf("UDB could not connect to gateways: %+v",
			err)
	}

	// SB: Trying to always register.
	// I think it's needed for some things to work correctly.
	// Need a more accurate descriptor of what the method actually does than
	// Register, or to remove the things that aren't actually used for
	// registration.

	userID, err = clientObj.Register(true, regCode, "",
		"")
	if err != nil {
		udb.Log.FATAL.Panicf("Could not register: %v", err)
	}

	return userID
}

// Login to the server using the user id generated from Init
func Login(userID *id.User) {
	_, err := clientObj.Login(userID)

	if err != nil {
		udb.Log.FATAL.Panicf("Could not log into the server: %s", err)
	}
}

// outputDsaPubKeyToJSON encodes the DSA public key and user ID to JSON and
// outputs it to the specified directory with the specified file name.
func outputDsaPubKeyToJSON(publicKey *signature.DSAPublicKey, userID *id.User,
	dir, fileName string) {
	// Encode the public key for the pem format
	encodedKey, err := publicKey.PemEncode()
	if err != nil {
		jww.ERROR.Printf("Error Pem encoding public key: %s", err)
	}

	// Setup struct that will dictate the JSON structure
	jsonStruct := struct {
		Id             *id.User
		Dsa_public_key string
	}{
		Id:             userID,
		Dsa_public_key: string(encodedKey),
	}

	// Generate JSON from structure
	data, err := json.MarshalIndent(jsonStruct, "", "\t")
	if err != nil {
		jww.ERROR.Printf("Error encoding structure to JSON: %s", err)
	}

	// Get the user's home directory
	homeDir, err := homedir.Dir()
	if err != nil {
		jww.ERROR.Printf("Unable to retrieve user's home directory: %s", err)
	}

	// Write JSON to file
	err = ioutil.WriteFile(homeDir+"/"+dir+"/"+fileName, data, 0644)
	if err != nil {
		jww.ERROR.Printf("Error writing JSON file: %s", err)
	}
}
