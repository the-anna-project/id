// Package id provides a simple ID generating service using pseudo random
// strings.
package id

import (
	"github.com/the-anna-project/random"
)

const (
	// Hex128 creates a new hexa decimal encoded, pseudo random, 128 bit hash.
	Hex128 int = 16
	// Hex512 creates a new hexa decimal encoded, pseudo random, 512 bit hash.
	Hex512 int = 64
	// Hex1024 creates a new hexa decimal encoded, pseudo random, 1024 bit hash.
	Hex1024 int = 128
	// Hex2048 creates a new hexa decimal encoded, pseudo random, 2048 bit hash.
	Hex2048 int = 256
	// Hex4096 creates a new hexa decimal encoded, pseudo random, 4096 bit hash.
	Hex4096 int = 512
)

// Config represents the configuration used to create a new ID service.
type Config struct {
	// Dependencies.

	// RandomService represents a factory returning random numbers.
	RandomService random.Service

	// Settings.

	// HashChars represents the characters used to create hashes.
	HashChars string
	// Length defines the ID bit size.
	Length int
}

// DefaultConfig provides a default configuration to create a new ID service
// by best effort.
func DefaultConfig() Config {
	var err error

	var randomService random.Service
	{
		randomConfig := random.DefaultConfig()
		randomService, err = random.New(randomConfig)
		if err != nil {
			panic(err)
		}
	}

	newConfig := Config{
		// Dependencies.
		RandomService: randomService,

		// Settings.
		HashChars: "abcdef0123456789",
		Length:    Hex128,
	}

	return newConfig
}

// New creates a new configured ID service.
func New(config Config) (Service, error) {
	// Dependencies.
	if config.RandomService == nil {
		return nil, maskAnyf(invalidConfigError, "random service must not be empty")
	}

	// Settings.
	if config.HashChars == "" {
		return nil, maskAnyf(invalidConfigError, "hash characters must not be empty")
	}
	if config.Length == 0 {
		return nil, maskAnyf(invalidConfigError, "length must not be empty")
	}

	newService := &service{
		// Dependencies.
		randomService: config.RandomService,

		// Settings.
		hashChars: config.HashChars,
		length:    config.Length,
	}

	return newService, nil
}

type service struct {
	// Dependencies.
	randomService random.Service

	// Settings.
	hashChars string
	length    int
}

func (s *service) New() (string, error) {
	ID, err := s.WithType(s.length)
	if err != nil {
		return "", maskAny(err)
	}

	return ID, nil
}

func (s *service) WithType(length int) (string, error) {
	n := int(length)
	max := len(s.hashChars)

	newRandomNumbers, err := s.randomService.CreateNMax(n, max)
	if err != nil {
		return "", maskAny(err)
	}

	b := make([]byte, n)

	for i, r := range newRandomNumbers {
		b[i] = s.hashChars[r]
	}

	return string(b), nil
}
