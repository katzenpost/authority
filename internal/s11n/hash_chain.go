// hash_chain.go - authority document header hash chain.
// Copyright (C) 2019  David Stainton.
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
	"bytes"
	"errors"

	"github.com/katzenpost/core/epochtime"
	"github.com/ugorji/go/codec"
	"golang.org/x/crypto/blake2b"
)

const HashChainHeaderVersion = "0.0.0"

var (
	cborHandle *codec.CborHandle

	// ErrGenesisMismatch is an error message indicating genesis mismatch.
	ErrGenesisMismatch = errors.New("genesis mismatch error")

	// ErrHashMismatch is an error whcih indicates a link in the hash chain is broken.
	ErrHashMismatch = errors.New("hash mismatch error")
)

type Header struct {
	Version     string
	Epoch       uint64
	PrevHash    []byte
	PayloadHash []byte
}

func (h *Header) ToBytes() ([]byte, error) {
	out := []byte{}
	enc := codec.NewEncoderBytes(&out, cborHandle)
	err := enc.Encode(h)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (h *Header) hash() ([]byte, error) {
	headerBytes, err := h.ToBytes()
	if err != nil {
		return nil, err
	}
	hash := blake2b.Sum256(headerBytes)
	return hash[:], nil
}

func makeGenesisHeader() *Header {
	current, _, _ := epochtime.Now()
	return &Header{
		Version:     HashChainHeaderVersion,
		Epoch:       current - 1,
		PrevHash:    []byte{},
		PayloadHash: []byte{},
	}
}

func NewChain() []*Header {
	genesis := makeGenesisHeader()
	return []*Header{genesis}
}

func VerifyChain(chain []*Header, genesis *Header) (bool, error) {
	genesisBytes, err := genesis.ToBytes()
	if err != nil {
		return false, err
	}
	genesisChainBytes, err := chain[0].ToBytes()
	if err != nil {
		return false, err
	}
	if !bytes.Equal(genesisBytes, genesisChainBytes) {
		return false, ErrGenesisMismatch
	}

	for i := len(chain) - 1; i >= 0; i-- {
		prevHash := []byte{}
		if i != 0 {
			prevHash, err = chain[i-1].hash()
		}
		if err != nil {
			return false, err
		}
		if !bytes.Equal(prevHash, chain[i].PrevHash) {
			return false, ErrHashMismatch
		}
	}
	return true, nil
}

func AppendChain(chain []*Header, payload []byte) ([]*Header, error) {
	payloadHash := blake2b.Sum256(payload)
	epoch, _, _ := epochtime.Now()
	prevHash, err := chain[len(chain)-1].hash()
	if err != nil {
		return nil, err
	}
	link := Header{
		Version:     HashChainHeaderVersion,
		Epoch:       epoch,
		PrevHash:    prevHash,
		PayloadHash: payloadHash[:],
	}
	chain = append(chain, &link)
	return chain, nil
}

func init() {
	cborHandle = new(codec.CborHandle)
	cborHandle.Canonical = true
}
