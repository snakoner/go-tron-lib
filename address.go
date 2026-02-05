package tron

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
var b58Indexes [256]int

func init() {
	for i := range b58Indexes {
		b58Indexes[i] = -1
	}
	for i, c := range b58Alphabet {
		b58Indexes[c] = i
	}
}

func TronBase58ToHex(b58 string) (string, error) {
	raw, err := base58Decode(strings.TrimSpace(b58))
	if err != nil {
		return "", err
	}
	if len(raw) < 5 {
		return "", errors.New("invalid base58: too short")
	}

	payload := raw[:len(raw)-4]
	checksum := raw[len(raw)-4:]
	want := checksum4(payload)
	if !bytes.Equal(checksum, want) {
		return "", errors.New("invalid base58 checksum")
	}

	if len(payload) != 21 {
		return "", errors.New("invalid tron address length")
	}
	return hex.EncodeToString(payload), nil
}

func TronHexToBase58(hexAddr string) (string, error) {
	hexAddr = strings.TrimSpace(strings.TrimPrefix(hexAddr, "0x"))
	b, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", err
	}
	if len(b) != 21 {
		return "", errors.New("invalid tron hex address length (want 21 bytes: 0x41 + 20 bytes)")
	}
	payload := b
	sum := checksum4(payload)
	full := append(payload, sum...)
	return base58Encode(full), nil
}

func checksum4(b []byte) []byte {
	h1 := sha256.Sum256(b)
	h2 := sha256.Sum256(h1[:])
	return h2[:4]
}

func base58Encode(b []byte) string {
	zeros := 0
	for zeros < len(b) && b[zeros] == 0 {
		zeros++
	}

	var encoded []byte
	x := make([]byte, len(b))
	copy(x, b)

	for len(x) > 0 && !(len(x) == 1 && x[0] == 0) {
		mod := divmod58(x)
		encoded = append(encoded, b58Alphabet[mod])
		x = trimLeadingZeros(x)
	}

	for i := 0; i < zeros; i++ {
		encoded = append(encoded, b58Alphabet[0])
	}

	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}
	return string(encoded)
}

func base58Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty base58 string")
	}

	zeros := 0
	for zeros < len(s) && s[zeros] == '1' {
		zeros++
	}

	out := []byte{0}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		val := -1
		if ch < 255 {
			val = b58Indexes[ch]
		}
		if val < 0 {
			return nil, errors.New("invalid base58 character")
		}
		out = mulAdd256(out, byte(val), 58)
	}

	out = trimLeadingZeros(out)

	if zeros > 0 {
		prefix := make([]byte, zeros)
		out = append(prefix, out...)
	}
	return out, nil
}

func mulAdd256(out []byte, add byte, base byte) []byte {
	carry := int(add)
	for i := len(out) - 1; i >= 0; i-- {
		v := int(out[i])*int(base) + carry
		out[i] = byte(v & 0xff)
		carry = v >> 8
	}
	for carry > 0 {
		out = append([]byte{byte(carry & 0xff)}, out...)
		carry >>= 8
	}
	return out
}

func divmod58(x []byte) byte {
	var rem int
	for i := 0; i < len(x); i++ {
		num := rem*256 + int(x[i])
		x[i] = byte(num / 58)
		rem = num % 58
	}
	return byte(rem)
}

func trimLeadingZeros(b []byte) []byte {
	i := 0
	for i < len(b) && b[i] == 0 {
		i++
	}
	if i == 0 {
		return b
	}
	if i == len(b) {
		return []byte{0}
	}
	return b[i:]
}

func PrivateKeyHexToAddressBase58(privateKeyHex string) (string, error) {
	privateKeyHex = strings.TrimSpace(privateKeyHex)
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0X")

	priv, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}

	pubBytes := crypto.FromECDSAPub(&priv.PublicKey)

	hash := crypto.Keccak256(pubBytes[1:])

	addr20 := hash[12:]

	payload := append([]byte{0x41}, addr20...)

	h1 := sha256.Sum256(payload)
	h2 := sha256.Sum256(h1[:])
	checksum := h2[:4]

	full := append(payload, checksum...)

	return base58Encode(full), nil
}
