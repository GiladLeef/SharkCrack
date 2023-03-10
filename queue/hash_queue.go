package queue

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/GiladLeef/SharkCrack/interfaces"
	"github.com/GiladLeef/SharkCrack/models"
	"os"
)

type HashQueue struct {
	hashes chan string
	config models.Config
}

func NewHashQueue(p interfaces.ConfigProvider) interfaces.FlushingQueue {
	config := p.GetConfig()
	hashes := make(chan string, config.HashQueueBuffer)
	return &HashQueue{hashes, *config}
}

func (q HashQueue) Size() int {
	return len(q.hashes)
}

func (q HashQueue) Get() (string, error) {
	for {
		select {
		case hash := <-q.hashes:
			return hash, nil
		default:
			err := errors.New("No hashes in queue.")
			return "", err
		}
	}
}

func (q HashQueue) Put(hash string) error {
	select {
	case q.hashes <- hash:
		return nil
	default:
		err := fmt.Errorf("No room in buffer. Discarding hash: %+v\n", hash)
		return err
	}
}

func (q HashQueue) Flush() error {
	if q.config.FlushToFile {
		return q.flushToFile()
	}
	_, err := q.emptyChannel()

	return err
}

func (q HashQueue) flushToFile() error {
	file, err := os.OpenFile(q.config.ComputedHashOverflowPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	hashes, err := q.emptyChannel()
	if err != nil {
		return err
	}

	for _, hash := range hashes {
		fmt.Fprintln(writer, hash)
	}

	return writer.Flush()
}

func (q HashQueue) emptyChannel() ([]string, error) {
	initialSize := len(q.hashes)
	var hashes []string
	for i := 0; i < initialSize; i++ {
		hash, err := q.Get()
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, hash)
	}

	return hashes, nil
}
