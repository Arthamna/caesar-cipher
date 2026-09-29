package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
)

func generateKey(length int) []byte {
	key := make([]byte, length)
	rand.Read(key)
	return key
}

func otp(text, key []byte) []byte {
	result := make([]byte, len(text))

	for i := range text {
		result[i] = text[i] ^ key[i]
	}

	return result
}

func sender(address string) {
	fmt.Print("Masukkan pesan: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	plaintext := []byte(scanner.Text())

	key := generateKey(len(plaintext))
	ciphertext := otp(plaintext, key)

	keyHex := hex.EncodeToString(key)
	cipherHex := hex.EncodeToString(ciphertext)

	fmt.Printf("Address = %q\n", address)
	fmt.Println("Key:", keyHex)
	fmt.Println("Ciphertext:", cipherHex)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", keyHex)
	fmt.Fprintf(conn, "%s\n", cipherHex)
}

func receiver(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Receiver listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		scanner := bufio.NewScanner(conn)

		scanner.Scan()
		keyHex := scanner.Text()

		scanner.Scan()
		cipherHex := scanner.Text()

		key, _ := hex.DecodeString(keyHex)
		ciphertext, _ := hex.DecodeString(cipherHex)

		plaintext := otp(ciphertext, key)

		fmt.Println("Received Key:", keyHex)
		fmt.Println("Received Ciphertext:", cipherHex)
		fmt.Println("Decrypted Plaintext:", string(plaintext))

		conn.Close()
	}
}

func main() {
	mode := flag.String("mode", "", "sender or receiver")
	address := flag.String("addr", ":8080", "address to connect to")
	flag.Parse()

	if *mode == "sender" {
		sender(*address)
	} else if *mode == "receiver" {
		receiver(*address)
	} else {
		fmt.Println("Invalid mode. Please specify -mode=sender or -mode=receiver.")
	}
}