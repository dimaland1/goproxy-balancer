# API Load Balancer en Go

Un load balancer (équilibreur de charge) HTTP écrit en Go qui distribue le trafic entre plusieurs serveurs backend en utilisant l'algorithme Round Robin.

## Fonctionnalités

- **Équilibrage de charge Round Robin** : Distribution équitable des requêtes entre les serveurs configurés
- **Surveillance de santé** : Vérification périodique de la disponibilité des serveurs cibles
- **Tolérance aux pannes** : Détection automatique des serveurs défaillants et redirection du trafic
- **Journalisation avancée** : Rotation des logs et différents niveaux de logs (INFO, WARNING, ERROR)
- **Métriques et statistiques** : Endpoint dédié pour surveiller les performances du load balancer
- **API RESTful** : Endpoints pour vérifier la santé et obtenir des statistiques

## Structure du Projet

```
api-golang/
├── main.go           # Point d'entrée de l'application
├── logger.go         # Configuration du système de journalisation
├── balancer.go       # Implémentation de l'algorithme de load balancing
├── handlers.go       # Gestionnaires HTTP et endpoints API
└── logs/             # Dossier contenant les fichiers de logs
```

## Prérequis

- Go 1.18 ou version ultérieure

## Installation

1. Clonez ce dépôt :
   ```
   git clone https://github.com/votre-username/api-golang.git
   cd api-golang
   ```

2. Installez les dépendances :
   ```
   go get gopkg.in/natefinch/lumberjack.v2
   ```

3. Compilez et exécutez l'application :
   ```
   go run .
   ```
   
   Ou utilisez:
   ```
   go run main.go logger.go balancer.go handlers.go
   ```

## Configuration

Vous pouvez modifier la liste des serveurs cibles dans `main.go`. Par défaut, l'application utilise les serveurs suivants :

```go
targetServers := []string{
    "https://jalalazouzout.vercel.app/",
    "https://example.com/",
    "https://api.example.org/",
    // Ajoutez d'autres serveurs au besoin
}
```

## Utilisation

Une fois lancé, le load balancer sera accessible à l'adresse `http://localhost:8080/`.

### Endpoints disponibles

- **`/`** : Point d'entrée principal qui redirige vers les serveurs cibles
- **`/health`** : Renvoie l'état de santé du load balancer
- **`/stats`** : Fournit des statistiques détaillées sur le load balancer

### Exemple de requête

```bash
# Rediriger vers un serveur cible
curl http://localhost:8080/

# Vérifier l'état de santé du load balancer
curl http://localhost:8080/health

# Obtenir des statistiques
curl http://localhost:8080/stats
```

## Logs

Les fichiers de logs sont automatiquement créés dans le dossier `logs/` avec un nom basé sur la date courante (ex: `2023-04-15.log`). Les logs sont également affichés dans la console.

La rotation des logs est configurée pour :
- Créer un nouveau fichier chaque jour
- Limiter la taille des fichiers à 10 Mo
- Conserver les logs pendant 90 jours
- Compresser les anciens fichiers de logs

## Surveillance de Santé

Le load balancer effectue des vérifications de santé des serveurs cibles toutes les 30 secondes. Si un serveur échoue à 3 vérifications consécutives, il est temporairement retiré de la rotation jusqu'à ce qu'il redevienne disponible.

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou à soumettre une pull request.

## Licence

Ce projet est sous licence MIT.