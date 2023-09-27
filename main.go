package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DecodedPacket struct {
	Short1     uint16
	Chars1     string
	SingleByte uint8
	Chars2     string
	Short2     uint16
	Chars3     string
	Long       uint32
}

func decodePacket(packet []byte) (DecodedPacket, error) {
	if len(packet) != 44 {
		return DecodedPacket{}, fmt.Errorf("Invalid packet size, expected 44 bytes")
	}

	var decoded DecodedPacket

	decoded.Short1 = binary.BigEndian.Uint16(packet[0:2])
	decoded.Chars1 = string(packet[2:14])
	decoded.SingleByte = packet[14]
	decoded.Chars2 = string(packet[15:23])
	decoded.Short2 = binary.BigEndian.Uint16(packet[23:25])
	decoded.Chars3 = string(packet[25:40])
	decoded.Long = binary.BigEndian.Uint32(packet[40:44])

	return decoded, nil
}

func main() {
	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		packet := []byte{0x04, 0xD2, 0x6B, 0x65, 0x65, 0x70, 0x64, 0x65, 0x63, 0x6F, 0x64, 0x69, 0x6E, 0x67, 0x38, 0x64, 0x6F, 0x6E, 0x74, 0x73, 0x74, 0x6F, 0x70, 0x03, 0x15, 0x63, 0x6F, 0x6E, 0x67, 0x72, 0x61, 0x74, 0x75, 0x6C, 0x61, 0x74, 0x69, 0x6F, 0x6E, 0x73, 0x07, 0x5B, 0xCD, 0x15}

		// decodedPacket, err := decodePacket(packet)
		// if err != nil {
		// 	fmt.Println("Error decoding packet:", err)
		// 	return
		// }

		// fmt.Printf("Decoded struct: %+v\n", decodedPacket)
		// fmt.Fprintf(w, "Hello, World!")

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		fmt.Printf("body struct: %+v\n", body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Process the data from the request body

		decodedPacket, err := decodePacket(packet)
		if err != nil {
			fmt.Println("Error decoding packet:", err)
			return
		}
		// Send a response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, decodedPacket)

	})

	port := 8080
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server listening on port %d...\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Error:", err)
	}
}
