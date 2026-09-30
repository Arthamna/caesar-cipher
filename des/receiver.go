package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
)

func receiver(cfg config) error {
	listener, err := net.Listen("tcp", cfg.address)
	if err != nil {
		return fmt.Errorf("menjalankan receiver: %w", err)
	}
	defer listener.Close()

	fmt.Println("Receiver listening on", cfg.address)
	fmt.Println("Block mode:", cfg.blockMode)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error menerima koneksi:", err)
			continue
		}
		handleConnection(conn, cfg)
		conn.Close()
	}
}

func handleConnection(conn net.Conn, cfg config) {
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Println("Error membaca ciphertext:", err)
		}
		return
	}

	cipherHex := scanner.Text()
	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		fmt.Println("Error: ciphertext bukan hexadecimal yang valid:", err)
		return
	}
	plaintext, err := decryptMessage(ciphertext, cfg.key, cfg.iv, cfg.blockMode)
	if err != nil {
		fmt.Println("Error mendekripsi ciphertext:", err)
		return
	}

	fmt.Println("Received Ciphertext:", cipherHex)
	fmt.Println("Decrypted Plaintext:", string(plaintext))
}
