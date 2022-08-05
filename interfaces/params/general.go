////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the UserDiscovery layer

package params

type General struct {
	SessionPath     string
	ProtoUserJson   []byte
	Ndf             string
	PermCert        []byte
	BannedUserList  string
	BannedRegexList string
	Ed25519Key      []byte

	Database
	IO
	Twilio
}
