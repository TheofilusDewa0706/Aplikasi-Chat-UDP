package main

import (
	"fmt"
	"net"
	"strings"
)

var users = make(map[string]*net.UDPAddr)

func main() {
	address := ":8080"
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Server berjalan di", address)

	buffer := make([]byte, 1024)

	for {
		// Menerima pesan dari client
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		message := string(buffer[:n])
		message = strings.TrimSpace(message)

		// Jika pesan adalah format "JOIN:NamaUser", maka ini adalah user baru
		if strings.HasPrefix(message, "JOIN:") {
			userName := strings.TrimPrefix(message, "JOIN:")
			users[userName] = addr
			fmt.Printf("%s has joined the chat from %s\n", userName, addr)

			// Mengirim notifikasi ke semua user lain bahwa user baru telah bergabung
			for name, userAddr := range users {
				if name != userName {
					notification := fmt.Sprintf("%s has joined the chat", userName)
					conn.WriteToUDP([]byte(notification), userAddr)
				}
			}
			continue
		}

		// Jika pesan adalah format "LEAVE:NamaUser", maka user akan meninggalkan chat
		if strings.HasPrefix(message, "LEAVE:") {
			userName := strings.TrimPrefix(message, "LEAVE:")
			fmt.Printf("%s has disconnected\n", userName)

			// Mengirim notifikasi ke semua user lain bahwa user telah meninggalkan chat
			for name, userAddr := range users {
				if name != userName {
					notification := fmt.Sprintf("%s has left the chat", userName)
					conn.WriteToUDP([]byte(notification), userAddr)
				}
			}
			// Hapus user dari daftar
			delete(users, userName)
			continue
		}

		// Broadcast pesan ke semua user
		for _, userAddr := range users {
			if userAddr != addr {
				conn.WriteToUDP(buffer[:n], userAddr)
			}
		}
	}
}
