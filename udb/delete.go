package udb

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
)

func DeleteFact(msg *pb.FactRemovalRequest, client *api.Client, store storage.Storage, auth *connect.Auth) (*messages.Ack, error) {
	// Nil checks
	// Can we have a blank fact?
	if msg == nil || msg.RemovalData == nil || msg.UID == nil {
		return &messages.Ack{}, errors.New("Unable to parse required " +
			"fields in registration message")
	}

	// Ensure client is properly authenticated
	if !auth.IsAuthenticated || auth.Sender.IsDynamicHost() {
		return &messages.Ack{}, connect.AuthError(auth.Sender.GetId())
	}

	// Generate the hash function and hash the fact
	h, err := hash.NewCMixHash()
	if err != nil {
		return &messages.Ack{}, errors.New("Removal could not be " +
			"completed do to internal error, please try again later")
	}
	h.Write(msg.RemovalData.Digest())
	hashFact := h.Sum(nil)

	// Get the user who owns the fact
	users := store.Search([][]byte{hashFact})
	if len(users) != 1 {
		return &messages.Ack{}, errors.New("Removal could not be " +
			"completed do to internal error, please try again later")
	}
	// Unmarshal the owner ID
	uid, err := id.Unmarshal(users[0].Id)
	if err != nil {
		return &messages.Ack{}, errors.New("Removal could not be " +
			"completed do to internal error, please try again later")
	}
	// Check the owner ID matches the sender ID
	if uid != auth.Sender.GetId() {
		return &messages.Ack{}, errors.New("Removal could not be " +
			"completed because you do not own this fact.")
	}

	// Delete the fact
	err = store.DeleteFact(hashFact)
	if err != nil {
		return &messages.Ack{}, errors.New("Removal could not be " +
			"completed do to internal error, please try again later")
	}

	return &messages.Ack{}, nil
}
