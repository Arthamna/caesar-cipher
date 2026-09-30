package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type config struct {
	role      string
	address   string
	blockMode string
	key       []byte
	iv        []byte
}

func parseConfig() (config, error) {
	role := flag.String("mode", "", "peran program: sender atau receiver")
	address := flag.String("addr", ":8080", "alamat TCP")
	keyHex := flag.String("key", "", "key DES dalam 16 digit hexadecimal")
	blockMode := flag.String("block-mode", "ecb", "mode blok: ecb atau cbc")
	padding := flag.String("padding", "pkcs7", "padding: pkcs7")
	ivHex := flag.String("iv", "", "IV CBC dalam 16 digit hexadecimal")
	flag.Parse()

	cfg := config{
		role:      strings.ToLower(*role),
		address:   *address,
		blockMode: strings.ToLower(*blockMode),
	}

	if cfg.role != "sender" && cfg.role != "receiver" {
		return cfg, errors.New("-mode harus sender atau receiver")
	}
	if cfg.blockMode != "ecb" && cfg.blockMode != "cbc" {
		return cfg, errors.New("-block-mode harus ecb atau cbc")
	}
	if strings.ToLower(*padding) != "pkcs7" {
		return cfg, errors.New("-padding yang didukung hanya pkcs7")
	}

	key, err := decodeHexParameter("key", *keyHex)
	if err != nil {
		return cfg, err
	}
	if _, err := generateRoundKeys(key); err != nil {
		return cfg, err
	}
	cfg.key = key

	if cfg.blockMode == "cbc" {
		iv, err := decodeHexParameter("IV", *ivHex)
		if err != nil {
			return cfg, err
		}
		cfg.iv = iv
	}

	return cfg, nil
}

func decodeHexParameter(name, value string) ([]byte, error) {
	if len(value) != desBlockSize*2 {
		return nil, fmt.Errorf("%s harus tepat 16 digit hexadecimal", name)
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("%s bukan hexadecimal yang valid: %w", name, err)
	}
	return decoded, nil
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if cfg.role == "sender" {
		err = sender(cfg)
	} else {
		err = receiver(cfg)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}
