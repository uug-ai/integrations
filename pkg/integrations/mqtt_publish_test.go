package integrations

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

// fakeBroker accepts MQTT 3.1.1 connections, acknowledges CONNECT and QoS
// 1 PUBLISH packets, and records what was published.
type fakeBroker struct {
	listener net.Listener
	mu       sync.Mutex
	topic    string
	payload  []byte
	qos      byte
	retain   bool
	user     string
}

func newFakeBroker(t *testing.T) *fakeBroker {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	b := &fakeBroker{listener: listener}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go b.serve(conn)
		}
	}()
	t.Cleanup(func() { listener.Close() })
	return b
}

func (b *fakeBroker) uri() string { return "tcp://" + b.listener.Addr().String() }

func (b *fakeBroker) serve(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		header, err := r.ReadByte()
		if err != nil {
			return
		}
		length, multiplier := 0, 1
		for {
			digit, err := r.ReadByte()
			if err != nil {
				return
			}
			length += int(digit&127) * multiplier
			if digit&128 == 0 {
				break
			}
			multiplier *= 128
		}
		body := make([]byte, length)
		if _, err := io.ReadFull(r, body); err != nil {
			return
		}
		switch header >> 4 {
		case 1: // CONNECT: protocol name, level, flags, keep alive, client id, then user/password.
			nameLen := int(binary.BigEndian.Uint16(body))
			flags := body[2+nameLen+1]
			rest := body[2+nameLen+4:]
			idLen := int(binary.BigEndian.Uint16(rest))
			rest = rest[2+idLen:]
			if flags&0x80 != 0 {
				userLen := int(binary.BigEndian.Uint16(rest))
				b.mu.Lock()
				b.user = string(rest[2 : 2+userLen])
				b.mu.Unlock()
			}
			conn.Write([]byte{0x20, 0x02, 0x00, 0x00})
		case 3: // PUBLISH
			topicLen := int(binary.BigEndian.Uint16(body))
			topic := string(body[2 : 2+topicLen])
			qos := (header >> 1) & 0x03
			rest := body[2+topicLen:]
			var id []byte
			if qos > 0 {
				id, rest = rest[:2], rest[2:]
			}
			b.mu.Lock()
			b.topic, b.payload, b.qos, b.retain = topic, append([]byte(nil), rest...), qos, header&0x01 == 1
			b.mu.Unlock()
			if qos == 1 {
				conn.Write([]byte{0x40, 0x02, id[0], id[1]})
			}
		case 12: // PINGREQ
			conn.Write([]byte{0xD0, 0x00})
		case 14: // DISCONNECT
			return
		}
	}
}

func TestMQTTPublishSendsAcknowledgedQoS1Message(t *testing.T) {
	broker := newFakeBroker(t)
	sender, err := NewMQTT(NewMQTTOptions().SetURI(broker.uri()).SetUsername("hub").SetPassword("secret").Build())
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Publish("kerberos/vault/AKIA1", []byte(`{"action":"forward-recordings"}`)); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	broker.mu.Lock()
	defer broker.mu.Unlock()
	if broker.topic != "kerberos/vault/AKIA1" || string(broker.payload) != `{"action":"forward-recordings"}` {
		t.Fatalf("broker got %q %q", broker.topic, broker.payload)
	}
	if broker.qos != 1 || broker.retain {
		t.Fatalf("qos %d retain %v, want 1 and false", broker.qos, broker.retain)
	}
	if broker.user != "hub" {
		t.Fatalf("broker user %q, want hub", broker.user)
	}
}

func TestMQTTPublishReportsAnUnreachableBroker(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	uri := "tcp://" + listener.Addr().String()
	listener.Close()

	previous := mqttPublishTimeout
	mqttPublishTimeout = 2 * time.Second
	t.Cleanup(func() { mqttPublishTimeout = previous })

	sender, err := NewMQTT(NewMQTTOptions().SetURI(uri).Build())
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Publish("kerberos/vault/AKIA1", []byte("x")); err == nil {
		t.Fatal("Publish() to a closed port returned no error")
	}
}

func TestMQTTPublishRejectsAnEmptyTopic(t *testing.T) {
	sender, err := NewMQTT(NewMQTTOptions().SetURI("tcp://127.0.0.1:1").Build())
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Publish(" ", []byte("x")); err == nil {
		t.Fatal("Publish() accepted an empty topic")
	}
}
