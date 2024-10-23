package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	serverAddress := "localhost:8080"

	udpAddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	// Meminta nama user saat memulai
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan nama Anda: ")
	userName, _ := reader.ReadString('\n')
	userName = strings.TrimSpace(userName)

	// Mengirim pesan join ke server
	joinMessage := fmt.Sprintf("JOIN:%s", userName)
	_, err = conn.Write([]byte(joinMessage))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Printf("%s has joined the chat\n", userName)

	// Goroutine untuk menangani sinyal keluar (Ctrl+C)
	go func() {
		// Tangkap sinyal `Ctrl+C` atau `os.Interrupt`
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

		// Tunggu sinyal
		<-signalChan
		// Kirim pesan "LEAVE" ke server
		leaveMessage := fmt.Sprintf("LEAVE:%s", userName)
		_, _ = conn.Write([]byte(leaveMessage))
		fmt.Printf("\n%s has left the chat\n", userName)
		os.Exit(0)
	}()

	// Menjalankan goroutine untuk menerima pesan dari server
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error menerima pesan:", err)
				continue
			}

			// Menampilkan pesan yang diterima dari server
			message := string(buffer[:n])

			// Pastikan pesan dari user lain saja yang ditampilkan
			if !strings.HasPrefix(message, userName) {
				fmt.Printf("\r%s\n", message) // Menampilkan pesan user lain
				fmt.Printf("%s: ", userName)  // Prompt untuk pesan setelah user lain mengirim
			}
		}
	}()

	// Loop untuk mengirim pesan dari input user
	for {
		fmt.Printf("%s: ", userName) // Menampilkan prompt user
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Format pesan user
		messageToSend := fmt.Sprintf("%s: %s", userName, message)
		_, err := conn.Write([]byte(messageToSend))
		if err != nil {
			fmt.Println("Error mengirim pesan:", err)
			return
		}
	}
	wg.Wait()
}
