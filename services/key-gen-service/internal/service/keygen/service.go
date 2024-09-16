package keygensvc

import (
	"errors"
)

var ErrKeysBufferClosed = errors.New("keys buffer has already closed")

type KeyGenerator interface {
	Generate() string
}

type keyGenService struct {
	generator KeyGenerator
	keysBuf   chan string

	closech chan struct{}
}

func New(gen KeyGenerator, keysBufCap int) *keyGenService {
	svc := &keyGenService{
		generator: gen,
		keysBuf:   make(chan string, keysBufCap),
		closech:   make(chan struct{}),
	}

	go svc.backgroundGenerate()

	return svc
}

func (kg *keyGenService) backgroundGenerate() {
	go func() {
		for {
			select {
			case <-kg.closech:
				close(kg.keysBuf)
				return
			case kg.keysBuf <- kg.generator.Generate():
			}
		}
	}()
}

func (kg *keyGenService) Stop() {
	close(kg.closech)
}

func (kg *keyGenService) GetKey() (string, error) {
	if key, ok := <-kg.keysBuf; ok {
		return key, nil
	}

	return "", ErrKeysBufferClosed
}
