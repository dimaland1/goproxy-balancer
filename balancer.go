package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Structure pour gérer les serveurs cibles avec round-robin
type RoundRobinBalancer struct {
	serverURLs []*url.URL               // Liste des URLs des serveurs cibles
	current    int                      // Index du serveur actuel dans la rotation
	mutex      sync.Mutex               // Mutex pour la synchronisation des accès concurrents
	proxies    []*httputil.ReverseProxy // Liste des reverse proxies pour chaque serveur cible
	logger     *log.Logger              // Logger pour enregistrer les activités
}

// Créer un nouveau balancer
func NewRoundRobinBalancer(targets []string, logger *log.Logger) (*RoundRobinBalancer, error) {
	var serverURLs []*url.URL
	var proxies []*httputil.ReverseProxy

	for _, target := range targets {
		url, err := url.Parse(target)
		if err != nil {
			return nil, fmt.Errorf("erreur parsing URL %s: %w", target, err)
		}

		serverURLs = append(serverURLs, url)
		proxies = append(proxies, httputil.NewSingleHostReverseProxy(url))
	}

	return &RoundRobinBalancer{
		serverURLs: serverURLs,
		current:    0,
		mutex:      sync.Mutex{},
		proxies:    proxies,
		logger:     logger,
	}, nil
}

// Obtenir le prochain serveur dans la rotation
func (b *RoundRobinBalancer) NextServer() (*url.URL, *httputil.ReverseProxy) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	server := b.serverURLs[b.current]
	proxy := b.proxies[b.current]

	// Passer au serveur suivant
	b.current = (b.current + 1) % len(b.serverURLs)

	b.logger.Printf("Selected server: %s (index: %d)", server.Host, b.current)

	return server, proxy
}

// Ajouter un serveur au balancer
func (b *RoundRobinBalancer) AddServer(target string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	url, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("erreur parsing URL %s: %w", target, err)
	}

	b.serverURLs = append(b.serverURLs, url)
	b.proxies = append(b.proxies, httputil.NewSingleHostReverseProxy(url))

	b.logger.Printf("Added new server: %s", url.Host)
	return nil
}

// Obtenir la liste des serveurs
func (b *RoundRobinBalancer) GetServers() []*url.URL {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Créer une copie pour éviter les modifications externes
	servers := make([]*url.URL, len(b.serverURLs))
	copy(servers, b.serverURLs)

	return servers
}
