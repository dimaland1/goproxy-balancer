package main

import (
	"io"
	"log"
	"os"
	"github.com/natefinch/lumberjack"
)

// Niveaux de logs
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// Logger configurable avec niveaux
type ConfigurableLogger struct {
	logger *log.Logger
	level  LogLevel
}

// Méthodes pour chaque niveau de log
func (l *ConfigurableLogger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func (l *ConfigurableLogger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

func (l *ConfigurableLogger) Warn(format string, v ...interface{}) {
	if l.level <= WARNING {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

func (l *ConfigurableLogger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Configuration du système de logging
func SetupLogging() (*log.Logger, *log.Logger, *ConfigurableLogger) {
	// Configuration du writer à rotation
	logWriter := &lumberjack.Logger{
		Filename:   "loadbalancer.log",
		MaxSize:    10,   // taille maximale en MB avant rotation
		MaxBackups: 5,    // nombre de fichiers de backup à conserver
		MaxAge:     28,   // nombre de jours de conservation des logs
		Compress:   true, // compresser les anciens logs
	}

	// Logger pour fichier uniquement
	fileLogger := log.New(logWriter, "[LOAD-BALANCER] ", log.LstdFlags|log.Lshortfile)

	// Logger pour console et fichier
	multiWriter := io.MultiWriter(os.Stdout, logWriter)
	consoleLogger := log.New(multiWriter, "[LOAD-BALANCER] ", log.LstdFlags|log.Lshortfile)

	// Logger configurable
	appLogger := &ConfigurableLogger{
		logger: fileLogger,
		level:  INFO,
	}

	return fileLogger, consoleLogger, appLogger
}
