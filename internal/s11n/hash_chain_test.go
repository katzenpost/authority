// hash_chain_test.go - Hash chain tests.
// Copyright (C) 2019  David Stainton
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

package s11n

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashChain(t *testing.T) {
	require := require.New(t)

	chain := NewChain()
	chain, err := AppendChain(chain, []byte("yo1"))
	require.NoError(err)
	chain, err = AppendChain(chain, []byte("yo2"))
	require.NoError(err)
	chain, err = AppendChain(chain, []byte("yo3"))
	require.NoError(err)

	ok, err := VerifyChain(chain, chain[0])
	require.NoError(err)
	require.True(ok)

	for _, link := range chain {
		bytes, err := link.ToBytes()
		require.NoError(err)
		t.Logf("element size %d", len(bytes))
	}
}
