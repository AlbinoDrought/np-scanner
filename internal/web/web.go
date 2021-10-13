package web

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"go.albinodrought.com/neptunes-pride/internal/actions"
	"go.albinodrought.com/neptunes-pride/internal/matches"
	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/notifications"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
	"go.albinodrought.com/neptunes-pride/internal/opsec"
	"go.albinodrought.com/neptunes-pride/internal/types"
)

//go:generate $GOPATH/bin/rice embed-go

type WebOptions struct {
	Address    string
	PollPeriod time.Duration
}

var DefaultWebOptions = WebOptions{
	Address:    ":38080",
	PollPeriod: time.Minute * 5,
}

func Run(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, guard notifications.Guard, sinks []notifications.Sink, options *WebOptions) error {
	if options == nil {
		options = &DefaultWebOptions
	}

	webServer := &webServer{
		ctx,
		db,
		client,
		guard,
		sinks,
	}

	server := &http.Server{
		Addr:    options.Address,
		Handler: webServer.Router(),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	if options.PollPeriod > 0 {
		go webServer.Poll(options.PollPeriod)
		log.Println("automatically polling every", options.PollPeriod)
	}

	log.Println("Serving on", options.Address)
	return server.ListenAndServe()
}

type webServer struct {
	ctx    context.Context
	db     matchstore.MatchStore
	client npapi.NeptunesPrideClient
	guard  notifications.Guard
	sinks  []notifications.Sink
}

func (ws *webServer) authorize(w http.ResponseWriter, r *http.Request, match *matches.Match) (matches.AccessProfile, bool) {
	if !match.HasAccessCode() {
		return matches.PermissiveAccessProfile(), true
	}

	accessProfile, err := match.CheckAccessCode([]byte(r.URL.Query().Get("access_code")))
	if err == nil {
		return accessProfile, true
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Wrong access code"))
	log.Printf("Authorization failed for match %v: %v", match.GameNumber, err)
	return matches.AccessProfile{}, false
}

func (ws *webServer) Poll(period time.Duration) {
	// written this way so the timer fires immediately on fn enter
	// eventually gets reset with the proper period
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-timer.C:
			pollResults, err := actions.PollAllMatches(ws.ctx, ws.db, ws.client, nil)
			if err != nil {
				log.Println("error while doing periodic pull", err)
			}

			for gameNumber, pollResult := range pollResults {
				if !pollResult.Changed {
					continue
				}

				match, err := ws.db.FindMatchOrFail(gameNumber)
				if err != nil {
					log.Println("failed to find match that changed during poll", gameNumber, err)
					continue
				}

				snapshot, err := ws.getMergedSnapshot(match, matches.PermissiveAccessProfile(), map[string]string{})
				if err != nil {
					log.Println("failed to get merged snapshot for notification use", gameNumber, err)
					continue
				}

				notifiables := actions.CheckNotifiables(gameNumber, snapshot)
				err = notifications.SendGuarded(ws.ctx, ws.guard, notifiables, ws.sinks)
				if err != nil {
					log.Println("failed to send notifications", err)
				}
			}

			timer.Reset(period)
		}
	}
}

func (ws *webServer) IndexMatch(w http.ResponseWriter, r *http.Request) {
	allMatches := []matches.Match{}

	err := ws.db.EachMatch(true, func(gameNumber string, match *matches.Match) {
		match.OldAccessCode = nil
		match.AccessProfiles = nil
		match.PlayerCreds = nil
		allMatches = append(allMatches, *match)
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		log.Printf("Error indexing matches: %v", err)
		return
	}

	sort.SliceStable(allMatches, func(i, j int) bool {
		if allMatches[i].Finished && !allMatches[j].Finished {
			// i is finished, j is not finished, move j before i
			return false
		}

		if !allMatches[i].Finished && allMatches[j].Finished {
			// i is not finished, j is finished, move i before j
			return true
		}

		// sort by name
		return strings.Compare(allMatches[i].Name, allMatches[j].Name) == -1
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allMatches)
}

func (ws *webServer) ShowMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameNumber := vars["gameNumber"]

	match, err := ws.db.FindMatchOrFail(gameNumber)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Match not found"))
		log.Printf("Match %v not found: %v", gameNumber, err)
		return
	}

	accessProfile, ok := ws.authorize(w, r, match)
	if !ok {
		return
	}

	match.OldAccessCode = nil
	match.AccessProfiles = nil

	visibleCredsMap := map[int]matches.PlayerCreds{}
	for i, creds := range match.PlayerCreds {
		if !accessProfile.CanViewPlayerID(i) {
			continue
		}
		creds.APIKey = ""
		visibleCredsMap[i] = creds
	}
	match.PlayerCreds = visibleCredsMap

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

const maxPlayerSnapshotLimit = 1000 // arbitrary limit

func (ws *webServer) AddApiKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameNumber := vars["gameNumber"]

	match, err := ws.db.FindMatchOrFail(gameNumber)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Match not found"))
		log.Printf("Match %v not found: %v", gameNumber, err)
		return
	}

	if _, ok := ws.authorize(w, r, match); !ok {
		return
	}

	key := r.URL.Query().Get("api-key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed ?api-key="))
		return
	}

	err = actions.SetCredentials(r.Context(), ws.db, ws.client, gameNumber, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to add key"))
		log.Printf("Failed to set credentials for %v: %v", gameNumber, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ws *webServer) IndexPlayerSnapshots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameNumber := vars["gameNumber"]

	match, err := ws.db.FindMatchOrFail(gameNumber)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Match not found"))
		log.Printf("Match %v not found: %v", gameNumber, err)
		return
	}

	accessProfile, ok := ws.authorize(w, r, match)
	if !ok {
		return
	}

	strPlayerID := vars["player"]

	playerID, err := strconv.Atoi(strPlayerID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed {player}"))
		return
	}

	if !accessProfile.CanViewPlayerID(playerID) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Not allowed to view {player}"))
		return
	}

	strLimit := r.URL.Query().Get("limit")
	if strLimit == "" {
		strLimit = "50"
	}

	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed ?limit"))
		return
	}

	if limit > maxPlayerSnapshotLimit {
		limit = maxPlayerSnapshotLimit
	}

	snapshotTimes, err := ws.db.ListSnapshotTimes(gameNumber, playerID, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		log.Printf("Error indexing matches: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshotTimes)
}

var ErrNoSnapshotsLoaded = errors.New("no snapshots loaded")

func (ws *webServer) getMergedSnapshot(match *matches.Match, accessProfile matches.AccessProfile, overrides map[string]string) (*types.APIResponse, error) {
	ignoredSnapshots := 0
	snapshotsToLoad := make(map[int]int64, len(match.PlayerCreds))
	for _, creds := range match.PlayerCreds {
		if !accessProfile.CanViewPlayerID(creds.PlayerUID) {
			continue
		}

		if creds.PollingDisabled {
			// these were messing up the "Last Polled" time, only load these if specified by user below
			snapshotsToLoad[creds.PlayerUID] = 0
		} else {
			snapshotsToLoad[creds.PlayerUID] = creds.LatestSnapshot
		}

		customSnapshot, ok := overrides[strconv.Itoa(creds.PlayerUID)]
		if ok && customSnapshot != "" && customSnapshot != "latest" {
			customSnapshotInt, err := strconv.ParseInt(customSnapshot, 10, 64)
			if err != nil {
				log.Printf("Malformed snapshot int for match %v player %v, %v: %v", match.GameNumber, creds.PlayerUID, customSnapshot, err)
				return nil, err
			}

			snapshotsToLoad[creds.PlayerUID] = customSnapshotInt
		}

		if snapshotsToLoad[creds.PlayerUID] == 0 {
			ignoredSnapshots++
		}
	}

	var err error
	i := 0
	loadedSnapshots := make([]*types.APIResponse, len(snapshotsToLoad)-ignoredSnapshots)
	for playerID, snapshotTime := range snapshotsToLoad {
		if snapshotTime == 0 {
			continue // ignored
		}

		loadedSnapshots[i], err = ws.db.FindSnapshot(match.GameNumber, playerID, snapshotTime)
		i++
		if err != nil {
			log.Printf("Could not load snapshot for match %v player %v, %v: %v", match.GameNumber, playerID, snapshotTime, err)
			return nil, err
		}
	}

	if len(loadedSnapshots) == 0 {
		return nil, ErrNoSnapshotsLoaded
	}

	return opsec.Merge(loadedSnapshots...), nil
}

func (ws *webServer) ShowMergedSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameNumber := vars["gameNumber"]

	match, err := ws.db.FindMatchOrFail(gameNumber)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Match not found"))
		log.Printf("Match %v not found: %v", gameNumber, err)
		return
	}

	accessProfile, ok := ws.authorize(w, r, match)
	if !ok {
		return
	}

	overrides := make(map[string]string)

	for _, creds := range match.PlayerCreds {
		stringID := strconv.Itoa(creds.PlayerUID)
		overrides[stringID] = r.URL.Query().Get(stringID)
	}

	mergedSnapshot, err := ws.getMergedSnapshot(match, accessProfile, overrides)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error merging snapshot"))
		log.Printf("Failed to get merged snapshot for match %v with overrides %+v: %v", gameNumber, overrides, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mergedSnapshot)
}

func (ws *webServer) Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/matches", ws.IndexMatch)
	r.HandleFunc("/api/matches/{gameNumber}", ws.ShowMatch)
	r.HandleFunc("/api/matches/{gameNumber}/api-key", ws.AddApiKey)
	r.HandleFunc("/api/matches/{gameNumber}/player-snapshots/{player}", ws.IndexPlayerSnapshots)
	r.HandleFunc("/api/matches/{gameNumber}/merged-snapshot", ws.ShowMergedSnapshot)

	box := rice.MustFindBox("packaged")
	r.PathPrefix("/").Handler(http.FileServer(&SPAFileSystem{box.HTTPBox()}))

	return r
}

// SPAFileSystem allows you to use history-based routing.
// All calls will load the index.html SPA unless a file asset is found.
type SPAFileSystem struct {
	http.FileSystem
}

func (fs *SPAFileSystem) forceFallback() (http.File, error) {
	return fs.FileSystem.Open("/index.html")
}

// Open a file or respond with the fallback contents
func (fs *SPAFileSystem) Open(name string) (http.File, error) {
	file, err := fs.FileSystem.Open(name)

	// prevent redirect loops, ignore root /
	if strings.TrimLeft(name, "/") == "" {
		return file, err
	}

	// api calls should never load SPA
	if strings.HasPrefix(name, "/api") {
		return file, err
	}

	// load SPA if file doesn't exist
	if err != nil && os.IsNotExist(err) {
		return fs.forceFallback()
	}

	// file actually exists, return it
	return file, err
}
