// client.go - Katzenpost non-voting authority client.
// Copyright (C) 2018  David Stainton
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package client implements the Katzenpost non-voting authority client.
package client

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/katzenpost/authority/voting/server/config"
	"github.com/katzenpost/core/crypto/eddsa"
	"github.com/katzenpost/core/crypto/rand"
	"github.com/katzenpost/core/epochtime"
	"github.com/katzenpost/core/log"
	"github.com/stretchr/testify/assert"
)

type mockDialer struct {
	serverConn net.Conn
	clientConn net.Conn
	dialCh     chan interface{}
}

func newMockDialer() *mockDialer {
	d := new(mockDialer)
	d.serverConn, d.clientConn = net.Pipe()
	d.dialCh = make(chan interface{}, 0)
	return d
}

func (d *mockDialer) dial(context.Context, string, string) (net.Conn, error) {
	return d.clientConn, nil
}

func (d *mockDialer) waitUntilDialed() {
	<-d.dialCh
}

func TestClient(t *testing.T) {
	assert := assert.New(t)

	dialer := newMockDialer()

	logBackend, err := log.New("", "DEBUG", false)
	assert.NoError(err, "wtf")

	peer1IdentityPrivateKey, err := eddsa.NewKeypair(rand.Reader)
	assert.NoError(err, "wtf")
	peer1LinkPrivateKey := peer1IdentityPrivateKey.ToECDH()

	cfg := &Config{
		LogBackend: logBackend,
		Authorities: []*config.AuthorityPeer{
			&config.AuthorityPeer{
				IdentityPublicKey: peer1IdentityPrivateKey.PublicKey(),
				LinkPublicKey:     peer1LinkPrivateKey.PublicKey(),
				Addresses:         []string{"127.0.0.1:1234"},
			},
		},
		DialContextFn: dialer.dial,
	}
	client, err := New(cfg)
	assert.NoError(err, "wtf")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	epoch, _, _ := epochtime.Now()
	doc, rawDoc, err := client.Get(ctx, epoch)
	assert.NoError(err, "wtf")
	assert.Equal(doc.Epoch, epoch)
	t.Logf("rawDoc size is %d", len(rawDoc))
}
