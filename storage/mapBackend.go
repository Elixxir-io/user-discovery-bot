////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the Map backend for the user discovery bot

package storage

import (
	"fmt"
	//"gitlab.com/elixxir/primitives/id"
	//"gitlab.com/elixxir/user-discovery-bot/fingerprint"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
)

// Insert or Update a User into the database
func (m *MapImpl) UpsertUser(user *User) error {
	m.lock.Lock()
	if m.users[] != nil {
		usr := NewUser()
		m.users[user.Id].Id =
	}
	m.users[user.Id] = user
	m.lock.Unlock()
	return nil
}

// Fetch a User from the database
func (m *MapImpl) GetUser(user *User) (*User, error) {
	m.lock.Lock()
	var err error
	retUser, ok := m.users[user.Id]
	if ok {
		err = errors.New(fmt.Sprintf(
			"User %+v has not been added!", user))
	}
	m.lock.Unlock()
	return &retUser, err
}
/*
// AddKey - Add a key stream, return the fingerprint
func (rs RamStorage) AddKey(value []byte) (string, error) {
	keyFingerprint := fingerprint.Fingerprint(value)

	// Error out if the key exists already
	_, ok := rs.Keys[keyFingerprint]
	if ok {
		return "", fmt.Errorf("fingerprint already exists: %s", keyFingerprint)
	}

	rs.Keys[keyFingerprint] = value
	return keyFingerprint, nil
}

// GetKey - Get a key based on the key id (retval of AddKey)
func (rs RamStorage) GetKey(keyId string) ([]byte, bool) {
	publicKey, ok := rs.Keys[keyId]
	return publicKey, ok
}

// AddUserKey - Add a user id to keyId (not used in high security)
func (rs RamStorage) AddUserKey(userId *id.User, keyId string) error {
	_, ok := rs.Users[*userId]
	if ok {
		return fmt.Errorf("UserId already exists: %d", userId)
	}
	rs.Users[*userId] = keyId
	return nil
}

// GetUserKey - Get a user's keyId (not used in high security)
func (rs RamStorage) GetUserKey(userId *id.User) (string, bool) {
	keyId, ok := rs.Users[*userId]
	return keyId, ok
}

// AddUserID - Add an email to userID mapping
func (rs RamStorage) AddUserID(email string, userID *id.User) error {
	_, ok := rs.UserIDs[email]
	if ok {
		return fmt.Errorf("email already exists: %s", email)
	}
	rs.UserIDs[email] = *userID
	return nil
}

// GetUserID - Get a user's ID from registered email
func (rs RamStorage) GetUserID(email string) (id.User, bool) {
	userID, ok := rs.UserIDs[email]
	return userID, ok
}

// AddValue - Add a searchable value (e-mail, nickname, etc)
func (rs RamStorage) AddValue(value string, valType ValueType,
	keyId string) error {
	_, ok := rs.KeyVal[valType]
	if !ok {
		rs.KeyVal[valType] = make(map[string][]string)
	}
	_, ok = rs.KeyVal[valType][value]
	if !ok {
		rs.KeyVal[valType][value] = make([]string, 0)
	}
	keyIds, _ := rs.KeyVal[valType][value]
	keyIds = append(keyIds, keyId)
	rs.KeyVal[valType][value] = keyIds
	return nil
}

// GetKeys - Returns all values that match the search criteria
func (rs RamStorage) GetKeys(value string, valType ValueType) (
	[]string, bool) {
	_, ok := rs.KeyVal[valType]
	if ok {
		keyIds, ok := rs.KeyVal[valType][value]
		if ok {
			return keyIds, ok
		}
	}
	return nil, false
}*/
