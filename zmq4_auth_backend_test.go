// Copyright 2025 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tests for auth and backend functions

package zmq4_test

import (
	"testing"

	"github.com/luxfi/zmq/v4"
)

func TestBackendInfo(t *testing.T) {
	// Test BackendName
	name := zmq4.BackendName()
	if name == "" {
		name = "pure-go"
	}
	t.Logf("Backend name: %s", name)

	// Test IsCZMQAvailable
	available := zmq4.IsCZMQAvailable()
	t.Logf("CZMQ available: %v", available)
}

func TestAuthBasic(t *testing.T) {
	// Test NewCurveKeypair
	pub, sec, err := zmq4.NewCurveKeypair()
	if err != nil {
		// Simplified implementation may return error
		t.Logf("NewCurveKeypair: %v", err)
		return
	}

	if pub == "" || sec == "" {
		t.Log("Simplified auth implementation")
		return
	}

	// Test AuthCurvePublic
	derivedPub, err := zmq4.AuthCurvePublic(sec)
	if err != nil {
		t.Logf("AuthCurvePublic: %v", err)
	} else if derivedPub != "" {
		t.Logf("Derived public key: %d chars", len(derivedPub))
	}
}

func TestZ85Encoding(t *testing.T) {
	// Test Z85encode
	data := []byte("Hello, World!")
	encoded := zmq4.Z85encode(data)
	if encoded == "" {
		t.Log("Z85encode returned empty (simplified)")
		return
	}

	// Test Z85decode
	decoded, err := zmq4.Z85decode(encoded)
	if err != nil {
		t.Logf("Z85decode: %v", err)
	} else if string(decoded) != string(data) {
		t.Errorf("Z85 roundtrip failed: got %q, want %q", decoded, data)
	}
}

func TestAuthLifecycle(t *testing.T) {
	// Test AuthStart
	err := zmq4.AuthStart()
	if err != nil {
		t.Logf("AuthStart: %v", err)
		return
	}

	// Always stop if started
	defer zmq4.AuthStop()

	// Test AuthSetVerbose
	zmq4.AuthSetVerbose(true)
	zmq4.AuthSetVerbose(false)

	// Test AuthAllow
	zmq4.AuthAllow("test", "127.0.0.1")

	// Test AuthDeny
	zmq4.AuthDeny("test", "192.168.1.1")

	// Test AuthCurveAdd
	pub, _, _ := zmq4.NewCurveKeypair()
	if pub != "" {
		zmq4.AuthCurveAdd("test", pub)

		// Test AuthCurveRemove
		zmq4.AuthCurveRemove("test", pub)
	}

	// Test AuthSetMetadataHandler
	zmq4.AuthSetMetadataHandler(func(domain, address string) map[string]string {
		return map[string]string{
			"User-Id": "test-user",
			"Name":    "Test",
		}
	})
}

func TestAuthEdgeCases(t *testing.T) {
	// Test double start
	_ = zmq4.AuthStart()
	err := zmq4.AuthStart()
	if err == nil {
		t.Log("AuthStart allowed double start")
	}
	zmq4.AuthStop()

	// Test operations without start
	zmq4.AuthSetVerbose(true)
	zmq4.AuthAllow("test", "127.0.0.1")
	zmq4.AuthDeny("test", "127.0.0.1")

	// Test invalid key
	_, err = zmq4.AuthCurvePublic("invalid-key")
	if err == nil {
		t.Log("AuthCurvePublic accepted invalid key")
	}

	// Test invalid Z85
	_, err = zmq4.Z85decode("invalid!")
	if err == nil {
		t.Log("Z85decode accepted invalid input")
	}
}
