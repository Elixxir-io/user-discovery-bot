////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package udb

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/client/api"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

// Endpoint which handles a users attempt to register
func RegisterUser(msg *pb.UDBUserRegistration, client *api.Client,
	store storage.Storage) (*messages.Ack, error) {

	// Nil checks
	if msg == nil || msg.Frs == nil || msg.Frs.Fact == nil ||
		msg.IdentityRegistration == nil {
		return &messages.Ack{}, errors.New("Unable to parse required " +
			"fields in registration message")
	}


	// Parse the username and UserID
	username := msg.IdentityRegistration.Username
	uid, err := id.Unmarshal(msg.UID)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not parse UID sent over. " +
			"Please try again")
	}

	// Check if username is taken
	err = store.CheckUser(username, uid, msg.RSAPublicPem)
	if err != nil {
		return &messages.Ack{}, errors.Errorf("Username %s is already taken. " +
			"Please try again", username)
	}

	// Pull the public key out of the permissioning cert
	permCert := client.GetNDF().Registration.TlsCertificate
	permPubKey, err := LoadPermissioningPubKey(permCert)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify signature due " +
			"to internal error. Please try again later")
	}

	// Hash the permissioning signature
	h, err := hash.NewCMixHash()
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify signature due " +
			"to internal error. Please try again later")
	}
	h.Write(msg.PermissioningSignature)
	hashedPermSig := h.Sum(nil)

	// Verify the Permissioning signature provided
	err = rsa.Verify(permPubKey, hash.CMixHash, hashedPermSig, msg.PermissioningSignature, nil)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify permissioning signature")
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(msg.RSAPublicPem))
	if err != nil {
		return &messages.Ack{}, errors.New("Could not parse key passed in")
	}

	// Verify the signed identity data
	hashedIdentity := msg.IdentityRegistration.Digest()
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedIdentity, msg.IdentitySignature, nil)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify identity signature")
	}

	// Verify the signed fact
	hashedFact := msg.Frs.Fact.Digest()
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, msg.Frs.FactSig, nil)

	// Create fact off of username
	f := storage.Fact{
		FactHash:  hashedFact,
		UserId:    msg.UID,
		Fact:      msg.Frs.Fact.Fact,
		FactType:  uint8(msg.Frs.Fact.FactType),
		Signature: msg.Frs.FactSig,
		Verified:  true,
		Timestamp: time.Now(),
	}

	// Create the user to insert into the database
	u := &storage.User{
		Id:        msg.UID,
		RsaPub:    msg.RSAPublicPem,
		DhPub:     msg.IdentityRegistration.DhPubKey,
		Salt:      msg.IdentityRegistration.Salt,
		Signature: msg.PermissioningSignature,
		Facts:     []storage.Fact{f},
	}

	// Insert the user into the database
	err = storage.UserDiscoveryDB.InsertUser(u)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not register username due " +
			"to internal error. Please try again later")

	}

	return &messages.Ack{}, nil
}

// Loads permissioning public key from the certificate
func LoadPermissioningPubKey(cert string) (*rsa.PublicKey, error) {
	permCert, err := tls.LoadCertificate(cert)
	if err != nil {
		return nil, errors.Errorf("Could not decode permissioning tls cert file "+
			"into a tls cert: %v", err)
	}


	return tls.ExtractPublicKey(permCert)
}