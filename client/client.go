package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	// Connexion au serveur RPC
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Erreur lors de la connexion :", err)
	}
	defer client.Close()

	// Message à envoyer
	message := "Bonjour, serveur !"
	var reply string

	// Appel de la méthode RPC
	err = client.Call("Service.SendMessage", message, &reply)
	if err != nil {
		log.Fatal("Erreur lors de l'appel RPC :", err)
	}

	fmt.Println("Réponse du serveur :", reply)
}
