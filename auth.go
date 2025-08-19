// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

// Auth provides ZeroMQ authentication functionality
type Auth struct {
	mu      sync.RWMutex
	verbose bool
	running bool
	
	// CURVE authentication
	curveKeys map[string]string // domain -> public key
	
	// Allowed/denied addresses
	allow map[string][]string // domain -> addresses
	deny  map[string][]string // domain -> addresses
}

var (
	authInstance *Auth
	authOnce     sync.Once
)

// getAuth returns the singleton auth instance
func getAuth() *Auth {
	authOnce.Do(func() {
		authInstance = &Auth{
			curveKeys: make(map[string]string),
			allow:     make(map[string][]string),
			deny:      make(map[string][]string),
		}
	})
	return authInstance
}

// AuthStart starts the authentication service
func AuthStart() error {
	auth := getAuth()
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	if auth.running {
		return fmt.Errorf("authentication already running")
	}
	
	auth.running = true
	return nil
}

// AuthStop stops the authentication service
func AuthStop() {
	auth := getAuth()
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	auth.running = false
	auth.curveKeys = make(map[string]string)
	auth.allow = make(map[string][]string)
	auth.deny = make(map[string][]string)
}

// AuthSetVerbose sets verbose mode for authentication
func AuthSetVerbose(verbose bool) {
	auth := getAuth()
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	auth.verbose = verbose
}

// AuthAllow adds allowed addresses for a domain
func AuthAllow(domain string, addresses ...string) {
	auth := getAuth()
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	if auth.allow[domain] == nil {
		auth.allow[domain] = []string{}
	}
	auth.allow[domain] = append(auth.allow[domain], addresses...)
}

// AuthDeny adds denied addresses for a domain
func AuthDeny(domain string, addresses ...string) {
	auth := getAuth()
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	if auth.deny[domain] == nil {
		auth.deny[domain] = []string{}
	}
	auth.deny[domain] = append(auth.deny[domain], addresses...)
}

// AuthCurveAdd adds a CURVE public key for a domain
func AuthCurveAdd(domain, publicKey string) {
	auth := getAuth()
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	auth.curveKeys[domain] = publicKey
}

// AuthCurveRemove removes a CURVE public key from a domain
func AuthCurveRemove(domain, publicKey string) {
	auth := getAuth()
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	if auth.curveKeys[domain] == publicKey {
		delete(auth.curveKeys, domain)
	}
}

// NewCurveKeypair generates a new CURVE keypair
func NewCurveKeypair() (publicKey, secretKey string, err error) {
	// Generate 32 bytes for secret key
	secret := make([]byte, 32)
	if _, err = rand.Read(secret); err != nil {
		return "", "", err
	}
	
	// For now, use the secret as the public key (simplified)
	// In a real implementation, this would use actual CURVE cryptography
	public := make([]byte, 32)
	copy(public, secret)
	
	// Convert to Z85 encoding (simplified hex for now)
	publicKey = hex.EncodeToString(public)
	secretKey = hex.EncodeToString(secret)
	
	return publicKey, secretKey, nil
}

// AuthCurvePublic derives the public key from a secret key
func AuthCurvePublic(secretKey string) (string, error) {
	// Simplified implementation - in real CURVE, this would derive the public key
	secret, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", err
	}
	if len(secret) != 32 {
		return "", fmt.Errorf("invalid secret key length")
	}
	
	// For now, just return the secret as public (simplified)
	return secretKey, nil
}

// Z85encode encodes binary data to Z85 text format
func Z85encode(data []byte) string {
	// Simplified implementation using hex encoding
	// Real Z85 encoding would use the Z85 alphabet
	return hex.EncodeToString(data)
}

// Z85decode decodes Z85 text to binary data
func Z85decode(text string) ([]byte, error) {
	// Simplified implementation using hex decoding
	// Real Z85 decoding would use the Z85 alphabet
	return hex.DecodeString(text)
}

// MetadataHandler is called for authentication metadata
type MetadataHandler func(version, requestID, domain, address, identity, mechanism string, credentials ...string) map[string]string

var metadataHandler MetadataHandler

// AuthSetMetadataHandler sets the metadata handler for authentication
func AuthSetMetadataHandler(handler MetadataHandler) {
	metadataHandler = handler
}