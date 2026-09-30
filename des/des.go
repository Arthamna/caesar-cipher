package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
)

const desBlockSize = 8

func permute(input uint64, inputSize int, table []int) uint64 {
	var output uint64
	for _, position := range table {
		output = (output << 1) | ((input >> (inputSize - position)) & 1)
	}
	return output
}

func generateRoundKeys(key []byte) ([16]uint64, error) {
	var roundKeys [16]uint64
	if len(key) != desBlockSize {
		return roundKeys, fmt.Errorf("key harus tepat %d byte", desBlockSize)
	}

	for i, b := range key {
		if bits.OnesCount8(b)%2 != 1 {
			return roundKeys, fmt.Errorf("byte key ke-%d tidak memiliki odd parity", i+1)
		}
	}

	permuted := permute(binary.BigEndian.Uint64(key), 64, permutedChoice1)
	c := uint32((permuted >> 28) & 0x0fffffff)
	d := uint32(permuted & 0x0fffffff)

	for i, shift := range keyShiftSchedule {
		c = ((c << shift) | (c >> (28 - shift))) & 0x0fffffff
		d = ((d << shift) | (d >> (28 - shift))) & 0x0fffffff
		combined := (uint64(c) << 28) | uint64(d)
		roundKeys[i] = permute(combined, 56, permutedChoice2)
	}

	return roundKeys, nil
}

func feistel(right uint32, roundKey uint64) uint32 {
	expanded := permute(uint64(right), 32, expansionTable)
	mixed := expanded ^ roundKey

	var substituted uint32
	for i := 0; i < 8; i++ {
		sixBits := uint8((mixed >> (42 - 6*i)) & 0x3f)
		row := ((sixBits & 0x20) >> 4) | (sixBits & 0x01)
		column := (sixBits >> 1) & 0x0f
		substituted = (substituted << 4) | uint32(sBoxes[i][row][column])
	}

	return uint32(permute(uint64(substituted), 32, straightPermutation))
}

func cryptBlock(block uint64, roundKeys [16]uint64, decrypt bool) uint64 {
	permuted := permute(block, 64, initialPermutation)
	left := uint32(permuted >> 32)
	right := uint32(permuted)

	for round := 0; round < 16; round++ {
		keyIndex := round
		if decrypt {
			keyIndex = 15 - round
		}
		left, right = right, left^feistel(right, roundKeys[keyIndex])
	}

	combined := (uint64(right) << 32) | uint64(left)
	return permute(combined, 64, inverseInitialPermutation)
}

func addPKCS7Padding(data []byte) []byte {
	paddingLength := desBlockSize - len(data)%desBlockSize
	padded := make([]byte, len(data)+paddingLength)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(paddingLength)
	}
	return padded
}

func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%desBlockSize != 0 {
		return nil, errors.New("panjang data padding tidak valid")
	}

	paddingLength := int(data[len(data)-1])
	if paddingLength < 1 || paddingLength > desBlockSize {
		return nil, errors.New("nilai padding tidak valid")
	}
	for _, value := range data[len(data)-paddingLength:] {
		if int(value) != paddingLength {
			return nil, errors.New("isi padding tidak valid")
		}
	}
	return data[:len(data)-paddingLength], nil
}

func encryptMessage(plaintext, key, iv []byte, blockMode string) ([]byte, error) {
	roundKeys, err := generateRoundKeys(key)
	if err != nil {
		return nil, err
	}
	if blockMode == "cbc" && len(iv) != desBlockSize {
		return nil, fmt.Errorf("IV CBC harus tepat %d byte", desBlockSize)
	}

	padded := addPKCS7Padding(plaintext)
	ciphertext := make([]byte, len(padded))
	var previous uint64
	if blockMode == "cbc" {
		previous = binary.BigEndian.Uint64(iv)
	}

	for offset := 0; offset < len(padded); offset += desBlockSize {
		block := binary.BigEndian.Uint64(padded[offset : offset+desBlockSize])
		if blockMode == "cbc" {
			block ^= previous
		}
		encrypted := cryptBlock(block, roundKeys, false)
		binary.BigEndian.PutUint64(ciphertext[offset:offset+desBlockSize], encrypted)
		previous = encrypted
	}
	return ciphertext, nil
}

func decryptMessage(ciphertext, key, iv []byte, blockMode string) ([]byte, error) {
	if len(ciphertext) == 0 || len(ciphertext)%desBlockSize != 0 {
		return nil, errors.New("ciphertext harus berupa blok-blok 8 byte")
	}
d
	roundKeys, err := generateRoundKeys(key)
	if err != nil {
		return nil, err
	}
	if blockMode == "cbc" && len(iv) != desBlockSize {
		return nil, fmt.Errorf("IV CBC harus tepat %d byte", desBlockSize)
	}

	plaintext := make([]byte, len(ciphertext))
	var previous uint64
	if blockMode == "cbc" {
		previous = binary.BigEndian.Uint64(iv)
	}

	for offset := 0; offset < len(ciphertext); offset += desBlockSize {
		cipherBlock := binary.BigEndian.Uint64(ciphertext[offset : offset+desBlockSize])
		decrypted := cryptBlock(cipherBlock, roundKeys, true)
		if blockMode == "cbc" {
			decrypted ^= previous
			previous = cipherBlock
		}
		binary.BigEndian.PutUint64(plaintext[offset:offset+desBlockSize], decrypted)
	}

	return removePKCS7Padding(plaintext)
}
