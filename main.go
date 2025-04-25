package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// Configuration du système de logging
	fileLogger, consoleLogger, appLogger := SetupLogging()

	consoleLogger.Println("---> Starting the server with round-robin load balancing...")

	// Configuration de la rotation quotidienne des logs
	go func() {
		// Force une rotation quotidienne à minuit
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			consoleLogger.Println("Performing scheduled log rotation")
			// La rotation est gérée automatiquement par lumberjack
		}
	}()

	// Liste des serveurs cibles
	targetServers := []string{
		"https://jalalazouzout.vercel.app/",
		"https://example.com/",
		"https://api.example.org/",
		// Ajoutez d'autres serveurs au besoin
	}

	// Créer le balancer round-robin
	balancer, err := NewRoundRobinBalancer(targetServers, fileLogger)
	if err != nil {
		consoleLogger.Println("Erreur lors de la création du balancer:", err)
		os.Exit(1)
	}

	// Utilisation du logger configurable
	appLogger.Info("Server started with %d targets", len(targetServers))

	// Créer les gestionnaires HTTP
	handlers := NewProxyHandlers(balancer, fileLogger, consoleLogger)

	// Configuration des routes
	http.HandleFunc("/", handlers.MainHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/stats", handlers.StatsHandler)

	// Démarrer le serveur proxy sur le port 8080
	consoleLogger.Println("---> Démarrage du load balancer sur le port 8080...")
	consoleLogger.Println("---> Serveurs cibles configurés:")

	for i, server := range targetServers {
		consoleLogger.Printf("     [%d] %s", i+1, server)
	}

	portNumber := 8080
	serverAddr := fmt.Sprintf(":%d", portNumber)

	consoleLogger.Printf("---> Server listening at http://localhost%s", serverAddr)
	err = http.ListenAndServe(serverAddr, nil)
	if err != nil {
		consoleLogger.Println("Erreur lors du démarrage du serveur:", err)
		os.Exit(1)
	}
}
