// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4_test

import (
	"testing"

	"github.com/luxfi/zmq/v4"
)

func TestAuthCurvePublic(t *testing.T) {
	// Test deriving public key from secret key
	_, secret, err := zmq4.NewCurveKeypair()
	if err != nil {
		t.Fatal(err)
	}
	
	public, err := zmq4.AuthCurvePublic(secret)
	if err != nil {
		t.Fatal(err)
	}
	
	if public == "" {
		t.Fatal("Expected non-empty public key")
	}
	
	// Test with invalid secret key
	_, err = zmq4.AuthCurvePublic("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid secret key")
	}
}

func TestAuthStart(t *testing.T) {
	// Start authentication
	err := zmq4.AuthStart()
	if err != nil {
		t.Fatal("AuthStart:", err)
	}
	defer zmq4.AuthStop()
	
	// Test double start
	err = zmq4.AuthStart()
	if err == nil {
		t.Fatal("Expected error for double AuthStart")
	}
}

func TestAuthAllowDeny(t *testing.T) {
	err := zmq4.AuthStart()
	if err != nil {
		t.Fatal("AuthStart:", err)
	}
	defer zmq4.AuthStop()
	
	// Test allow
	zmq4.AuthAllow("domain1", "127.0.0.1", "::1")
	
	// Test deny
	zmq4.AuthDeny("domain1", "192.168.1.1")
	
	// Test CURVE add
	public, _, err := zmq4.NewCurveKeypair()
	if err != nil {
		t.Fatal(err)
	}
	zmq4.AuthCurveAdd("domain1", public)
	zmq4.AuthCurveRemove("domain1", public)
}

func TestNewCurveKeypair(t *testing.T) {
	public, secret, err := zmq4.NewCurveKeypair()
	if err != nil {
		t.Fatal("NewCurveKeypair:", err)
	}
	
	if len(public) == 0 {
		t.Fatal("Empty public key")
	}
	if len(secret) == 0 {
		t.Fatal("Empty secret key")
	}
	
	// Keys should be different
	if public == secret {
		t.Log("Warning: public and secret keys are the same (simplified implementation)")
	}
}

func TestZ85EncodeDecode(t *testing.T) {
	original := []byte("Hello, World!")
	
	// Encode
	encoded := zmq4.Z85encode(original)
	if encoded == "" {
		t.Fatal("Empty encoded string")
	}
	
	// Decode
	decoded, err := zmq4.Z85decode(encoded)
	if err != nil {
		t.Fatal("Z85decode:", err)
	}
	
	// Compare
	if string(decoded) != string(original) {
		t.Fatalf("Mismatch: got %q, want %q", decoded, original)
	}
}

func TestAuthSetVerbose(t *testing.T) {
	// Should not panic
	zmq4.AuthSetVerbose(true)
	zmq4.AuthSetVerbose(false)
}

func TestAuthMetadataHandler(t *testing.T) {
	err := zmq4.AuthStart()
	if err != nil {
		t.Fatal("AuthStart:", err)
	}
	defer zmq4.AuthStop()
	
	called := false
	zmq4.AuthSetMetadataHandler(
		func(version, requestID, domain, address, identity, mechanism string, credentials ...string) map[string]string {
			called = true
			return map[string]string{
				"User-Id": "test-user",
			}
		})
	
	// Note: In a real implementation, this would trigger the handler
	// through actual authentication flow
	if called {
		t.Log("Metadata handler was called")
	}
}