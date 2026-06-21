package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var adminClient *websocket.Conn

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket ulanishda xato:", err)
		return
	}
	defer conn.Close()

	role := r.URL.Query().Get("role")

	if role == "admin" {
		adminClient = conn
		log.Println("Developper Core (Admin) ulandi!")
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Println("Admin uzildi.")
				adminClient = nil
				break
			}
		}
	} else {
		log.Println("Sinfdosh (Foydalanuvchi) ulandi!")
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Sinfdosh uzildi.")
				break
			}

			if adminClient != nil {
				err = adminClient.WriteMessage(messageType, message)
				if err != nil {
					log.Println("Adminga kadr yuborishda xato:", err)
				}
			}
		}
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/ws", handleWebSocket)

	// Render xosting uchun portni dinamik olish
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	fmt.Println("Server ishga tushdi, port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Serverni yuklashda xatolik:", err)
	}
}
