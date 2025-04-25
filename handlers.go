package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Structure pour les gestionnaires HTTP
type ProxyHandlers struct {
	balancer      *RoundRobinBalancer
	logger        *log.Logger
	consoleLogger *log.Logger
}

// Créer un nouveau gestionnaire HTTP
func NewProxyHandlers(balancer *RoundRobinBalancer, logger, consoleLogger *log.Logger) *ProxyHandlers {
	return &ProxyHandlers{
		balancer:      balancer,
		logger:        logger,
		consoleLogger: consoleLogger,
	}
}

// Gestionnaire principal qui transfère les requêtes aux serveurs backend
func (h *ProxyHandlers) MainHandler(w http.ResponseWriter, r *http.Request) {
	// Obtenir le prochain serveur dans la rotation
	targetURL, proxy := h.balancer.NextServer()

	start := time.Now()
	requestID := fmt.Sprintf("%d", time.Now().UnixNano())

	h.consoleLogger.Printf("[REQ-%s] Proxying request: %s %s to %s", requestID, r.Method, r.URL.Path, targetURL.Host)
	h.logger.Printf("[REQ-%s] New request: %s %s from %s to %s", requestID, r.Method, r.URL.Path, r.RemoteAddr, targetURL.Host)

	// Modification du Host pour le serveur cible
	r.Host = targetURL.Host

	// Forward la requête au serveur cible
	proxy.ServeHTTP(w, r)

	elapsed := time.Since(start)
	h.consoleLogger.Printf("[REQ-%s] Request completed in %v", requestID, elapsed)
	h.logger.Printf("[REQ-%s] Request completed in %v", requestID, elapsed)
}

// Gestionnaire pour l'endpoint de santé
func (h *ProxyHandlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Health check from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)

	servers := h.balancer.GetServers()
	fmt.Fprintf(w, "Load Balancer OK - Managing %d servers", len(servers))
}

// Gestionnaire pour afficher les statistiques et la configuration
func (h *ProxyHandlers) StatsHandler(w http.ResponseWriter, r *http.Request) {
	servers := h.balancer.GetServers()

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Load Balancer Statistics</h1>")
	fmt.Fprintf(w, "<h2>Configured Servers (%d)</h2>", len(servers))
	fmt.Fprintf(w, "<ul>")

	for i, server := range servers {
		fmt.Fprintf(w, "<li>Server %d: %s</li>", i+1, server.Host)
	}

	fmt.Fprintf(w, "</ul>")
}
