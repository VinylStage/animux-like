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
	// Configure logger to include timestamps and all levels
	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	return file
}

func Start() error {
	logFile := initLogger()
	if logFile != nil {
		defer logFile.Close()
	}

	slog.Info("starting animux daemon", "version", "1.0")

	var err error
	petState, err = pet.LoadState()
	if err != nil {
		if err == pet.ErrNoPet {
			slog.Info("no pet state found, waiting for adoption")
		} else {
			slog.Error("failed to load state", "error", err.Error())
		}
	} else {
		slog.Info("loaded existing pet state", "name", petState.Name, "species", petState.Species)
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
		slog.Error("failed to listen on socket", "error", err.Error())
		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	defer listener.Close()
	defer os.Remove(sockPath)

	mux := http.NewServeMux()
	mux.HandleFunc("/status", withLogging(handleStatus))
	mux.HandleFunc("/feed", withLogging(handleFeed))
	mux.HandleFunc("/play", withLogging(handlePlay))
	mux.HandleFunc("/clean", withLogging(handleClean))
	mux.HandleFunc("/adopt", withLogging(handleAdopt))

	server := &http.Server{Handler: mux}

	go func() {
		if err := server.Serve(listener); err != nil {
			slog.Error("http server stopped", "error", err.Error())
		}
	}()

	slog.Info("daemon running", "socket_path", sockPath, "log_path", getLogPath())

	// Wait for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	slog.Info("shutting down daemon gracefully")
	if petState != nil {
		pet.SaveState(petState)
		slog.Info("saved pet state on shutdown")
	}
	return nil
}

// Middleware to log all HTTP requests
func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("received command request", "path", r.URL.Path, "query", r.URL.RawQuery)
		next.ServeHTTP(w, r)
		slog.Debug("completed command request", "path", r.URL.Path)
	}
}

func applyOfflineDecay() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		return
	}
	
	now := time.Now()
	elapsed := now.Sub(petState.LastUpdate)
	ticks := int(elapsed.Minutes() / 5) // 1 tick = 5 minutes
	
	if ticks > 0 {
		slog.Info("applying offline decay", "elapsed_duration", elapsed.String(), "ticks_missed", ticks)
		decayStats(ticks)
		petState.LastUpdate = now
		pet.SaveState(petState)
	} else {
		slog.Debug("no offline decay needed", "elapsed_duration", elapsed.String())
	}
}

func decayLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		stateMutex.Lock()
		if petState != nil {
			slog.Debug("running scheduled tick decay")
			decayStats(1)
			petState.LastUpdate = time.Now()
			pet.SaveState(petState)
			slog.Info("applied tick decay", 
				"hunger", petState.Hunger, 
				"happiness", petState.Happiness,
				"cleanliness", petState.Cleanliness,
				"is_sick", petState.IsSick)
		}
		stateMutex.Unlock()
	}
}

func decayStats(ticks int) {
	species := pet.SpeciesData[petState.Species]
	
	oldHunger := petState.Hunger
	oldHappy := petState.Happiness
	oldClean := petState.Cleanliness
	
	petState.Hunger -= species.HungerDecay * ticks
	petState.Happiness -= species.HappyDecay * ticks
	petState.Cleanliness -= species.CleanDecay * ticks
	
	if petState.Cleanliness < 50 && petState.Hunger < 30 {
		if !petState.IsSick {
			slog.Warn("pet became sick due to neglect!")
		}
		petState.IsSick = true
	}
	
	petState.Clamp()
	
	slog.Debug("stats decayed", 
		"hunger_diff", petState.Hunger - oldHunger,
		"happiness_diff", petState.Happiness - oldHappy,
		"cleanliness_diff", petState.Cleanliness - oldClean)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	if petState == nil {
		slog.Warn("status check failed: no pet exists")
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(petState)
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		slog.Warn("feed failed: no pet exists")
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	oldHunger := petState.Hunger
	petState.Hunger += 30
	petState.Clamp()
	pet.SaveState(petState)
	slog.Info("pet fed", "old_hunger", oldHunger, "new_hunger", petState.Hunger)
	json.NewEncoder(w).Encode(petState)
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		slog.Warn("play failed: no pet exists")
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	oldHappy := petState.Happiness
	oldHunger := petState.Hunger
	petState.Happiness += 30
	petState.Hunger -= 10
	petState.Clamp()
	pet.SaveState(petState)
	slog.Info("played with pet", 
		"old_happiness", oldHappy, "new_happiness", petState.Happiness,
		"old_hunger", oldHunger, "new_hunger", petState.Hunger)
	json.NewEncoder(w).Encode(petState)
}

func handleClean(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if petState == nil {
		slog.Warn("clean failed: no pet exists")
		http.Error(w, "no pet", http.StatusNotFound)
		return
	}
	oldClean := petState.Cleanliness
	petState.Cleanliness = 100
	wasSick := petState.IsSick
	petState.IsSick = false
	petState.Clamp()
	pet.SaveState(petState)
	slog.Info("cleaned pet", "old_cleanliness", oldClean, "new_cleanliness", petState.Cleanliness, "cured_sickness", wasSick)
	json.NewEncoder(w).Encode(petState)
}

func handleAdopt(w http.ResponseWriter, r *http.Request) {
	species := r.URL.Query().Get("species")
	name := r.URL.Query().Get("name")
	
	slog.Info("attempting to adopt pet", "species", species, "name", name)
	
	if _, ok := pet.SpeciesData[pet.SpeciesType(species)]; !ok {
		slog.Error("adoption failed: invalid species", "requested_species", species)
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
	slog.Info("new pet adopted successfully", "species", species, "name", name, "initial_stats", "100/100/100")
	json.NewEncoder(w).Encode(petState)
}
