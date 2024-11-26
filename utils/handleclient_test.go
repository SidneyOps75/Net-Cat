package utils

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

// mockConn implements net.Conn interface for testing
type mockConn struct {
	readData  *bytes.Buffer
	writeData *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (n int, err error)   { return m.readData.Read(b) }
func (m *mockConn) Write(b []byte) (n int, err error)  { return m.writeData.Write(b) }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic client creation",
			input:    "testuser\n",
			expected: "testuser",
		},
		{
			name:     "client name with whitespace",
			input:    "  test user  \n",
			expected: "test user",
		},
		{
			name:     "empty name",
			input:    "\n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock connection
			mockConn := &mockConn{
				readData:  bytes.NewBuffer([]byte(tt.input)),
				writeData: &bytes.Buffer{},
			}

			server := &Server{}
			client := NewClient(mockConn, server)

			if client == nil {
				t.Fatal("expected client to not be nil")
			}

			if client.Name != tt.expected {
				t.Errorf("expected client name to be %q, got %q", tt.expected, client.Name)
			}

			if client.Conn != mockConn {
				t.Error("expected client connection to match mock connection")
			}

			if client.Messages == nil {
				t.Error("expected messages channel to be initialized")
			}

			if cap(client.Messages) != 10 {
				t.Errorf("expected messages channel capacity to be 10, got %d", cap(client.Messages))
			}
		})
	}
}

func TestClientSendMessages(t *testing.T) {
	mockConn := &mockConn{
		readData:  &bytes.Buffer{},
		writeData: &bytes.Buffer{},
	}

	client := &Client{
		Conn:     mockConn,
		Name:     "testuser",
		Messages: make(chan string, 10),
	}

	// Start sending messages in a goroutine
	go client.SendMessages()

	// Test messages
	messages := []string{
		"Hello, world!\n",
		"This is a test message\n",
		"Testing 123\n",
	}

	// Send messages through the channel
	for _, msg := range messages {
		client.Messages <- msg
	}
	close(client.Messages)

	// Wait a brief moment for messages to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify written data
	written := mockConn.writeData.String()
	for _, msg := range messages {
		if !strings.Contains(written, msg) {
			t.Errorf("expected message %q to be written to connection", msg)
		}
	}
}

func TestClientListen(t *testing.T) {
	// Test messages to simulate reading from connection
	testMessages := []string{
		"Hello\n",
		"Test message\n",
		"", // Empty message should be skipped
		"Goodbye\n",
	}

	mockConn := &mockConn{
		readData:  bytes.NewBuffer([]byte(strings.Join(testMessages, ""))),
		writeData: &bytes.Buffer{},
	}

	client := &Client{
		Conn:     mockConn,
		Name:     "testuser",
		Messages: make(chan string, 10),
	}

	// Create broadcast channel
	broadcast := make(chan string, 10)

	// Start listening in a goroutine
	go client.Listen(broadcast)

	// Collect messages from broadcast channel
	var receivedMsgs []string
	timeout := time.After(100 * time.Millisecond)

messageLoop:
	for {
		select {
		case msg := <-broadcast:
			receivedMsgs = append(receivedMsgs, msg)
		case <-timeout:
			break messageLoop
		}
	}

	// Verify received messages
	for _, msg := range testMessages {
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		found := false
		timeFormat := time.Now().Format("2006-01-02 15:04:05")
		expectedFormat := fmt.Sprintf("[%s][testuser]: %s\n", timeFormat, msg)

		for _, received := range receivedMsgs {
			// Compare the message part only, ignoring timestamp
			if strings.Contains(received, msg) &&
				strings.Contains(received, "[testuser]") {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected to receive message containing %q with format similar to %q",
				msg, expectedFormat)
		}
	}
}
