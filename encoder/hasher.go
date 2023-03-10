package encoder

import (
	"fmt"
	"github.com/GiladLeef/SharkCrack/interfaces"
	"github.com/GiladLeef/SharkCrack/models"
	"hash"
	"io"
	"sync"
)

type Hasher struct {
	config          *models.Config
	logger          interfaces.Logger
	mux             *sync.Mutex
	requestQueue    interfaces.RequestQueue
	stopQueue       interfaces.ClientStopQueue
	submissionQueue interfaces.SubmissionQueue
	supportedHashes map[string]hash.Hash
	waiter          interfaces.Waiter
}

func NewHasher(c *models.Config, l interfaces.Logger, r interfaces.RequestQueue, s interfaces.SubmissionQueue, cl interfaces.ClientStopQueue, w interfaces.Waiter, m *sync.Mutex) interfaces.Encoder {
	return &Hasher{
		config:          c,
		logger:          l,
		mux:             m,
		requestQueue:    r,
		stopQueue:       cl,
		submissionQueue: s,
		supportedHashes: models.GetSupportedHashFunctions(),
		waiter:          w,
	}
}

func (e *Hasher) Start() error {
	e.logger.LogMessage("Starting hasher...")
	for {
		err := e.processOrSleep()
		if err != nil {
			e.updateStopQueue(err)
			return err
		}

		stopReason, err := e.stopQueue.Get()
		if err == nil {
			err = fmt.Errorf("Hasher observed updateStopQueue reason:\n\t%+v", stopReason)
			return err
		}
	}
}

func (e *Hasher) processOrSleep() error {
	hashingRequest, err := e.requestQueue.Get()
	if err != nil {
		e.waiter.Wait()
	} else {
		err = e.handleHashingRequest(hashingRequest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Hasher) handleHashingRequest(hashingRequest models.HashingRequest) error {
	if e.requestRequiresInflation(hashingRequest) {
		e.inflateHashingRequest(&hashingRequest)
	}
	hashSubmission := e.getHashSubmission(hashingRequest)
	if e.config.Verbose {
		numResults := len(hashSubmission.Results)
		logMessage := fmt.Sprintf("Hasher has created hash submission with hash type: %s and %d results", hashSubmission.HashType, numResults)
		e.logger.LogMessage(logMessage)
	}

	err := e.submissionQueue.Put(hashSubmission)
	for err != nil {
		return err
	}

	return nil
}

func (e *Hasher) requestRequiresInflation(hashingRequest models.HashingRequest) bool {
	return hashingRequest.Hash == nil
}

func (e *Hasher) inflateHashingRequest(hashingRequest *models.HashingRequest) {
	hashingRequest.Hash = e.supportedHashes[hashingRequest.HashName]
}

func (e *Hasher) getHashSubmission(hashingRequest models.HashingRequest) models.HashSubmission {
	passwordHashes := e.getPasswordHashes(hashingRequest.Hash, hashingRequest.Passwords)
	return models.HashSubmission{hashingRequest.HashName, passwordHashes}
}

func (e *Hasher) updateStopQueue(err error) {
	stopReason := models.ClientStopReason{
		Requester: "",
		Encoder:   err.Error(),
		Submitter: "",
	}

	var i uint16
	for i = 0; i < e.config.Threads - 1; i++ {
		e.stopQueue.Put(stopReason)
	}
}

func (e *Hasher) getPasswordHashes(hash hash.Hash, passwords []string) []string {
	var passwordHashes []string
	for _, password := range passwords {
		e.mux.Lock()
		passwordHash := e.getPasswordHash(hash, password)
		e.mux.Unlock()
		passwordHashes = append(passwordHashes, passwordHash)
	}

	return passwordHashes
}

func (e *Hasher) getPasswordHash(hash hash.Hash, password string) string {
	io.WriteString(hash, password)
	humanReadableHash := fmt.Sprintf("%x", hash.Sum(nil))
	hash.Reset()

	return password + ":" + humanReadableHash
}
