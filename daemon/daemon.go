package daemon

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/VinylStage/animux-like/pet"
)

var (
	stateMutex sync.Mutex
	petState   *pet.State
)

func getSocketPath() string {
	xdgRuntime := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntime == "" {
		xdgRuntime = "/tmp"
	}
	return filepath.Join(xdgRuntime, "animux.sock")
}

func getLogPath() string {
	xdgState := os.Getenv("XDG_STATE_HOME")
	if xdgState == "" {
		home, _ := os.UserHomeDir()
		xdgState = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(xdgState, "animux", "animux.log")
}

func initLogger() *os.File {
	logPath := getLogPath()
	os.MkdirAll(filepath.Dir(logPath), 0755)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
		return nil
	}
	logger := slog.New(slog.NewJSONHandler(file, nil))
	slog.SetDefault(logger)
	return file
}

func Start() error {
	logFile := initLogger()
	if logFile != nil {
		defer logFile.Close()
	}

	slog.Info("starting animux daemon")

	var err error
	petState, err = pet.LoadState()
	if err != nil {
		if err == pet.ErrNoPet {
			slog.Info("no pet state found, waiting for adoption")
		} else {
			slog.Error("failed to load state", "error", err)
		}
	} else {
		applyOfflineDecay()
	}

	// Start decay loop
	go decayLoop()

	// Setup HTTP server on Unix Socket
	sockPath := getSocketPath()
	os.Remove(sockPath)
	os.MkdirAll(filepath.Dir(sockPath), 0755)

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	defer listener.Close()
	defer os.Remove(sockPath)

	mux := http.NewServeMux()
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/feed", handleFeed)
	mux.HandleFunc("/play", handlePlay)
	mux.HandleFunc("/clean", handleClean)
	mux.HandleFunc("/adopt", handleAdopt)

	server := &http.Server{Handler: mux}

	go func() {
		if err := server.Serve(listener); err != nil {
			slog.Error("server stopped", "error", err)
		}
	}()

	slog.Info("daemon running, listening on socket", "path", sockPath)

	// Wait for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	slog.Info("shutting down daemon")
	if petState != nil {
		pet.SaveState(petState)
	}
	return nil
}

func applyOfflineDecay() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		return
	}
	
	// Calculate ticks missed
	now := time.Now()
	elapsed := now.Sub(petState.LastUpdate)
	ticks := int(elapsed.Minutes() / 5) // 1 tick = 5 minutes
	
	if ticks > 0 {
		decayStats(ticks)
		petState.LastUpdate = now
		pet.SaveState(petState)
		slog.Info("applied offline decay", "ticks", ticks)
	}
}

func decayLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		stateMutex.Lock()
		if petState != nil {
			decayStats(1)
			petState.LastUpdate = time.Now()
			pet.SaveState(petState)
			slog.Info("applied tick decay", "hunger", petState.Hunger, "happiness", petState.Happiness)
		}
		stateMutex.Unlock()
	}
}

func decayStats(ticks int) {
	species := pet.SpeciesData[petState.Species]
	petState.Hunger -= species.HungerDecay * ticks
	petState.Happiness -= species.HappyDecay * ticks
	petState.Cleanliness -= species.CleanDecay * ticks
	
	if petState.Cleanliness < 50 && petState.Hunger < 30 {
		petState.IsSick = true
	}
	petState.Clamp()
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	if petState == nil {
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(petState)
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	petState.Hunger += 30
	petState.Clamp()
	pet.SaveState(petState)
	slog.Info("pet fed")
	json.NewEncoder(w).Encode(petState)
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	petState.Happiness += 30
	petState.Hunger -= 10
	petState.Clamp()
	pet.SaveState(petState)
	slog.Info("played with pet")
	json.NewEncoder(w).Encode(petState)
}

func handleClean(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	petState.Cleanliness = 100
	petState.IsSick = false
	petState.Clamp()
	pet.SaveState(petState)
	slog.Info("cleaned pet")
	json.NewEncoder(w).Encode(petState)
}

func handleAdopt(w http.ResponseWriter, r *http.Request) {
	species := r.URL.Query().Get("species")
	name := r.URL.Query().Get("name")
	
	if _, ok := pet.SpeciesData[pet.SpeciesType(species)]; !ok {
		http.Error(w, "invalid species", http.StatusBadRequest)
		return
	}

	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	petState = &pet.State{
		Name:        name,
		Species:     pet.SpeciesType(species),
		Hunger:      100,
		Happiness:   100,
		Cleanliness: 100,
		IsSick:      false,
		LastUpdate:  time.Now(),
	}
	pet.SaveState(petState)
	slog.Info("new pet adopted", "species", species, "name", name)
	json.NewEncoder(w).Encode(petState)
}
