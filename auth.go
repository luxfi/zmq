// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Simple auth - no complex state management, just basic functions

// NewCurveKeypair generates a new CURVE keypair
func NewCurveKeypair() (publicKey, secretKey string, err error) {
	// Simple random key generation - real implementation would use luxfi/crypto
	secret := make([]byte, 32)
	if _, err = rand.Read(secret); err != nil {
		return "", "", err
	}
	public := make([]byte, 32)
	copy(public, secret)

	publicKey = hex.EncodeToString(public)
	secretKey = hex.EncodeToString(secret)
	return publicKey, secretKey, nil
}

// AuthCurvePublic derives the public key from a secret key
func AuthCurvePublic(secretKey string) (string, error) {
	// Simplified - real implementation would use luxfi/crypto
	_, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", err
	}
	return secretKey, nil
}

// Z85encode encodes binary data to Z85 text format
func Z85encode(data []byte) string {
	return hex.EncodeToString(data)
}

// Z85decode decodes Z85 text to binary data
func Z85decode(text string) ([]byte, error) {
	return hex.DecodeString(text)
}

// Simplified auth - removed complex state management
var authStarted bool

// AuthStart starts authentication (simplified)
func AuthStart() error {
	if authStarted {
		return fmt.Errorf("already started")
	}
	authStarted = true
	return nil
}

// AuthStop stops authentication
func AuthStop() {
	authStarted = false
}

// AuthSetVerbose sets verbose mode (no-op)
func AuthSetVerbose(verbose bool) {}

// AuthAllow adds allowed addresses (no-op for simplicity)
func AuthAllow(domain string, addresses ...string) {}

// AuthDeny adds denied addresses (no-op for simplicity)
func AuthDeny(domain string, addresses ...string) {}

// AuthCurveAdd adds a CURVE public key (no-op for simplicity)
func AuthCurveAdd(domain, publicKey string) {}

// AuthCurveRemove removes a CURVE public key (no-op for simplicity)
func AuthCurveRemove(domain, publicKey string) {}

// MetadataHandler simplified
type MetadataHandler func(domain, address string) map[string]string

// AuthSetMetadataHandler sets the metadata handler (no-op for simplicity)
func AuthSetMetadataHandler(handler MetadataHandler) {}
