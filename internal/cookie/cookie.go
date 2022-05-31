package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

type CookieName string

const Name CookieName = "userID"

const (
	UserIDSize = 36

	secretKey = "wow-so-secret-key"
)

func CheckUserID(ticket string) (userID, authTicket string) {
	oldHMAC, err := hex.DecodeString(ticket[UserIDSize:])
	if err != nil {
		return GetSignedUserID()
	}

	id := ticket[:UserIDSize]
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(id))

	newHMAC := h.Sum(nil)

	if hmac.Equal(newHMAC, oldHMAC) {
		return id, ticket
	}
	return GetSignedUserID()
}

func GetSignedUserID() (userID, authTicket string) {
	userID = uuid.NewString()
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(userID))
	return userID, userID + hex.EncodeToString(h.Sum(nil))
}
