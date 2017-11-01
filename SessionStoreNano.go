package nanostore

import (
	"errors"
	"unsafe"

	"github.com/aerogo/nano"
	"github.com/aerogo/session"
)

// SessionStoreNano is a store saving sessions in a nano database.
type SessionStoreNano struct {
	collection *nano.Collection
}

// interfaceStruct reflects Go's internal interface{} structure.
type interfaceStruct struct {
	Type unsafe.Pointer
	Data unsafe.Pointer
}

// New creates a session store using an Aerospike database.
func New(collection *nano.Collection) *SessionStoreNano {
	return &SessionStoreNano{
		collection: collection,
	}
}

// Get loads the initial session values from the database.
func (store *SessionStoreNano) Get(sid string) (*session.Session, error) {
	record, err := store.collection.Get(sid)

	if err != nil {
		return nil, err
	}

	interfaceContainer := *(*interfaceStruct)(unsafe.Pointer(&record))
	data := *(*map[string]interface{})(interfaceContainer.Data)
	return session.New(sid, data), nil
}

// Set updates the session values in the database.
func (store *SessionStoreNano) Set(sid string, session *session.Session) error {
	sessionData := session.Data()

	// Set with nil as data means we should delete the session.
	if sessionData == nil {
		existed := store.collection.Delete(sid)

		if !existed {
			return errors.New("Session doesn't exist")
		}

		return nil
	}

	store.collection.Set(sid, &sessionData)
	return nil
}
