package caesar

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"

)

func caesar(text string, shift int) string {
	var result string
	for _, c := range text {
		if c >= 'a' && c <= 'z' {
			result += string((c-'a'+rune(shift))%26 + 'a')
		} else if c >= 'A' && c <= 'Z' {
			result += string((c-'A'+rune(shift))%26 + 'A')
		} else {
			result += string(c)
		}
	}
	return result
}

func sender(address string, shift int) {
	fmt.Print("Masukkan pesan: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	plaintext := scanner.Text()
	ciphertext := caesar(plaintext, shift)
	fmt.Printf("Address = %q\n", address)

	fmt.Printf("Ciphertext: %s\n", ciphertext)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\n", ciphertext) // send to connection server
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}

func receiver(address string, shift int) {
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

		if scanner.Scan() { // scan return bool value
			ciphertext := scanner.Text()
			plaintext := caesar(ciphertext, -shift) // decypher just reverse the shift
			fmt.Println("Received Ciphertext:", ciphertext)
			fmt.Println("Decrypted Plaintext:", plaintext)
		}
		conn.Close()
	}
}

func main(){
	mode := flag.String("mode", "", "sender or receiver")
	address := flag.String("addr", ":8080", "address to connect to")
	shift := flag.Int("shift", 3, "shift value for Caesar cipher")
	flag.Parse()

	if *mode == "sender" {
		sender(*address, *shift)
	} else if *mode == "receiver" {
		receiver(*address, *shift)
	} else {
		fmt.Println("Invalid mode. Please specify -mode=sender or -mode=receiver.")
	}
}