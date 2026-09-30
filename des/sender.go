package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
)

func sender(cfg config) error {
	fmt.Print("Masukkan pesan: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("membaca plaintext: %w", err)
		}
		return nil
	}

	ciphertext, err := encryptMessage([]byte(scanner.Text()), cfg.key, cfg.iv, cfg.blockMode)
	if err != nil {
		return fmt.Errorf("mengenkripsi plaintext: %w", err)
	}
	cipherHex := hex.EncodeToString(ciphertext)

	fmt.Printf("Address: %s\n", cfg.address)
	fmt.Printf("Block mode: %s\n", cfg.blockMode)
	fmt.Println("Ciphertext:", cipherHex)

	conn, err := net.Dial("tcp", cfg.address)
	if err != nil {
		return fmt.Errorf("menghubungi receiver: %w", err)
	}
	defer conn.Close()

	if _, err := fmt.Fprintln(conn, cipherHex); err != nil {
		return fmt.Errorf("mengirim ciphertext: %w", err)
	}
	return nil
}
