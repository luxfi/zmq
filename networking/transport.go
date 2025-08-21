// Copyright (C) 2020-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package networking provides high-level ZMQ4 networking primitives
// for distributed systems communication.
package networking

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/luxfi/zmq/v4"
)

// Transport provides high-performance message passing using ZMQ4
type Transport struct {
	nodeID  string
	ctx     context.Context
	cancel  context.CancelFunc
	pub     zmq4.Socket
	sub     zmq4.Socket
	router  zmq4.Socket
	dealers map[string]zmq4.Socket
	config  Config

	mu       sync.RWMutex
	handlers map[string]MessageHandler
	peers    []string

	// Metrics
	msgSent     atomic.Uint64
	msgReceived atomic.Uint64
	msgDropped  atomic.Uint64

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// Config holds transport configuration
type Config struct {
	NodeID      string
	BasePort    int
	BindAddress string        // Default: "127.0.0.1"
	MaxRetries  int           // Default: 3
	RetryDelay  time.Duration // Default: 100ms
	BufferSize  int           // Default: 1000
}

// Message represents a network message
type Message struct {
	Type      string          `json:"type"`
	From      string          `json:"from"`
	To        string          `json:"to,omitempty"`
	SessionID []byte          `json:"session_id,omitempty"`
	Height    uint64          `json:"height,omitempty"`
	Round     uint32          `json:"round,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp int64           `json:"timestamp"`
}

// MessageHandler processes incoming messages
type MessageHandler func(msg *Message)

// DefaultConfig returns a default configuration
func DefaultConfig(nodeID string, basePort int) Config {
	return Config{
		NodeID:      nodeID,
		BasePort:    basePort,
		BindAddress: "127.0.0.1",
		MaxRetries:  3,
		RetryDelay:  100 * time.Millisecond,
		BufferSize:  1000,
	}
}

// New creates a new ZMQ4 transport
func New(ctx context.Context, config Config) *Transport {
	// Apply defaults
	if config.BindAddress == "" {
		config.BindAddress = "127.0.0.1"
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 100 * time.Millisecond
	}
	if config.BufferSize == 0 {
		config.BufferSize = 1000
	}

	tCtx, cancel := context.WithCancel(ctx)
	return &Transport{
		nodeID:   config.NodeID,
		ctx:      tCtx,
		cancel:   cancel,
		config:   config,
		handlers: make(map[string]MessageHandler),
		dealers:  make(map[string]zmq4.Socket),
		stopCh:   make(chan struct{}),
	}
}

// Start initializes the transport
func (t *Transport) Start() error {
	// PUB socket for broadcasting
	t.pub = zmq4.NewPub(t.ctx)
	pubAddr := fmt.Sprintf("tcp://%s:%d", t.config.BindAddress, t.config.BasePort)
	if err := t.pub.Listen(pubAddr); err != nil {
		return fmt.Errorf("failed to bind pub socket on %s: %w", pubAddr, err)
	}

	// SUB socket for receiving broadcasts
	t.sub = zmq4.NewSub(t.ctx)
	t.sub.SetOption(zmq4.OptionSubscribe, "")

	// ROUTER socket for direct messages
	t.router = zmq4.NewRouter(t.ctx)
	routerAddr := fmt.Sprintf("tcp://%s:%d", t.config.BindAddress, t.config.BasePort+1000)
	if err := t.router.Listen(routerAddr); err != nil {
		return fmt.Errorf("failed to bind router socket on %s: %w", routerAddr, err)
	}

	t.wg.Add(2)
	go t.subLoop()
	go t.routerLoop()

	return nil
}

// Stop gracefully shuts down the transport
func (t *Transport) Stop() {
	close(t.stopCh)
	t.cancel()
	t.wg.Wait()

	// Close sockets
	if t.pub != nil {
		t.pub.Close()
	}
	if t.sub != nil {
		t.sub.Close()
	}
	if t.router != nil {
		t.router.Close()
	}

	// Close dealer sockets
	t.mu.Lock()
	for _, dealer := range t.dealers {
		dealer.Close()
	}
	t.mu.Unlock()
}

// ConnectPeer establishes a connection to a peer
func (t *Transport) ConnectPeer(peerID string, port int) error {
	return t.ConnectPeerWithAddress(peerID, t.config.BindAddress, port)
}

// ConnectPeerWithAddress establishes a connection to a peer at a specific address
func (t *Transport) ConnectPeerWithAddress(peerID, address string, port int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check if already connected
	for _, p := range t.peers {
		if p == peerID {
			return nil // Already connected
		}
	}

	// Subscribe to peer's broadcasts
	subAddr := fmt.Sprintf("tcp://%s:%d", address, port)
	if err := t.sub.Dial(subAddr); err != nil {
		return fmt.Errorf("failed to connect sub to %s at %s: %w", peerID, subAddr, err)
	}

	// Create dealer for direct messages
	dealer := zmq4.NewDealer(t.ctx, zmq4.WithID(zmq4.SocketIdentity(t.nodeID)))

	routerAddr := fmt.Sprintf("tcp://%s:%d", address, port+1000)
	if err := dealer.Dial(routerAddr); err != nil {
		return fmt.Errorf("failed to connect dealer to %s at %s: %w", peerID, routerAddr, err)
	}

	t.dealers[peerID] = dealer
	t.peers = append(t.peers, peerID)

	return nil
}

// DisconnectPeer removes a peer connection
func (t *Transport) DisconnectPeer(peerID string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Close and remove dealer
	if dealer, ok := t.dealers[peerID]; ok {
		dealer.Close()
		delete(t.dealers, peerID)
	}

	// Remove from peers list
	newPeers := make([]string, 0, len(t.peers)-1)
	for _, p := range t.peers {
		if p != peerID {
			newPeers = append(newPeers, p)
		}
	}
	t.peers = newPeers
}

// Broadcast sends a message to all connected peers
func (t *Transport) Broadcast(msg *Message) error {
	msg.From = t.nodeID
	msg.Timestamp = time.Now().UnixNano()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	t.msgSent.Add(1)
	return t.pub.Send(zmq4.NewMsg(data))
}

// Send sends a direct message to a specific peer
func (t *Transport) Send(peerID string, msg *Message) error {
	t.mu.RLock()
	dealer, ok := t.dealers[peerID]
	t.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no connection to peer %s", peerID)
	}

	msg.From = t.nodeID
	msg.To = peerID
	msg.Timestamp = time.Now().UnixNano()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	t.msgSent.Add(1)
	return dealer.Send(zmq4.NewMsg(data))
}

// SendWithRetry sends a message with retry logic
func (t *Transport) SendWithRetry(peerID string, msg *Message) error {
	var lastErr error

	for i := 0; i < t.config.MaxRetries; i++ {
		if err := t.Send(peerID, msg); err == nil {
			return nil
		} else {
			lastErr = err
			if i < t.config.MaxRetries-1 {
				time.Sleep(time.Duration(i+1) * t.config.RetryDelay)
			}
		}
	}

	return fmt.Errorf("failed after %d retries: %w", t.config.MaxRetries, lastErr)
}

// BroadcastWithRetry broadcasts a message with retry logic
func (t *Transport) BroadcastWithRetry(msg *Message) error {
	var lastErr error

	for i := 0; i < t.config.MaxRetries; i++ {
		if err := t.Broadcast(msg); err == nil {
			return nil
		} else {
			lastErr = err
			if i < t.config.MaxRetries-1 {
				time.Sleep(time.Duration(i+1) * t.config.RetryDelay)
			}
		}
	}

	return fmt.Errorf("failed after %d retries: %w", t.config.MaxRetries, lastErr)
}

// RegisterHandler registers a message handler for a specific type
func (t *Transport) RegisterHandler(msgType string, handler MessageHandler) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.handlers[msgType] = handler
}

// UnregisterHandler removes a message handler
func (t *Transport) UnregisterHandler(msgType string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.handlers, msgType)
}

// GetPeers returns the list of connected peers
func (t *Transport) GetPeers() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	peers := make([]string, len(t.peers))
	copy(peers, t.peers)
	return peers
}

// GetNodeID returns this node's ID
func (t *Transport) GetNodeID() string {
	return t.nodeID
}

// GetMetrics returns transport metrics
func (t *Transport) GetMetrics() (sent, received, dropped uint64) {
	return t.msgSent.Load(), t.msgReceived.Load(), t.msgDropped.Load()
}

// subLoop processes broadcast messages
func (t *Transport) subLoop() {
	defer t.wg.Done()

	for {
		select {
		case <-t.stopCh:
			return
		case <-t.ctx.Done():
			return
		default:
			msg, err := t.sub.Recv()
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					return
				}
				// Transient error, continue
				continue
			}

			t.processMessage(msg.Bytes())
		}
	}
}

// routerLoop processes direct messages
func (t *Transport) routerLoop() {
	defer t.wg.Done()

	for {
		select {
		case <-t.stopCh:
			return
		case <-t.ctx.Done():
			return
		default:
			msg, err := t.router.Recv()
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					return
				}
				// Transient error, continue
				continue
			}

			// Router socket includes identity frame
			if msg.Frames != nil && len(msg.Frames) > 0 {
				// Last frame contains the actual message
				t.processMessage(msg.Frames[len(msg.Frames)-1])
			}
		}
	}
}

// processMessage handles incoming messages
func (t *Transport) processMessage(data []byte) {
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		// Silently drop malformed messages
		t.msgDropped.Add(1)
		return
	}

	// Skip our own broadcasts
	if message.From == t.nodeID && message.To == "" {
		return
	}

	t.msgReceived.Add(1)

	// Route to appropriate handler
	t.mu.RLock()
	handler, ok := t.handlers[message.Type]
	t.mu.RUnlock()

	if ok && handler != nil {
		// Call handler in goroutine to avoid blocking
		go handler(&message)
	}
}
