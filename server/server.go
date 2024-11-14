package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Définition de la structure du service
type Service struct{}

// Méthode exposée via RPC
func (s *Service) SendMessage(message string, reply *string) error {
	*reply = "Message reçu : " + message
	return nil
}

func main() {
	// Enregistrement du service
	service := new(Service)
	rpc.Register(service)

	// Écoute sur le port 1234
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Erreur lors de l'écoute :", err)
	}
	defer listener.Close()
	fmt.Println("Serveur en attente de connexions sur le port 1234...")

	// Acceptation des connexions entrantes
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Erreur lors de l'acceptation de la connexion :", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
