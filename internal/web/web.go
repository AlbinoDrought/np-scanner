package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
	"go.albinodrought.com/neptunes-pride/internal/opsec"
	"go.albinodrought.com/neptunes-pride/internal/types"
)

//go:generate $GOPATH/bin/rice embed-go

type WebOptions struct {
	Address string
}

var DefaultWebOptions = WebOptions{
	Address: ":38080",
}

func Run(ctx context.Context, db matchstore.MatchStore, client npapi.NeptunesPrideClient, options *WebOptions) error {
	if options == nil {
		options = &DefaultWebOptions
	}

	webServer := &webServer{
		ctx,
		db,
		client,
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

	log.Println("Serving on", options.Address)
	return server.ListenAndServe()
}

type webServer struct {
	ctx    context.Context
	db     matchstore.MatchStore
	client npapi.NeptunesPrideClient
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

	for i, creds := range match.PlayerCreds {
		creds.APIKey = ""
		match.PlayerCreds[i] = creds
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
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

	snapshotsToLoad := make(map[int]int64, len(match.PlayerCreds))
	for _, creds := range match.PlayerCreds {
		snapshotsToLoad[creds.PlayerUID] = creds.LatestSnapshot

		customSnapshot := r.URL.Query().Get(strconv.Itoa(creds.PlayerUID))
		if customSnapshot != "" {
			customSnapshotInt, err := strconv.ParseInt(customSnapshot, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Malformed snapshot int for player %v", creds.PlayerUID)
				log.Printf("Malformed snapshot int for match %v player %v, %v: %v", gameNumber, creds.PlayerUID, customSnapshot, err)
				return
			}

			snapshotsToLoad[creds.PlayerUID] = customSnapshotInt
		}
	}

	i := 0
	loadedSnapshots := make([]*types.APIResponse, len(snapshotsToLoad))
	for playerID, snapshotTime := range snapshotsToLoad {
		loadedSnapshots[i], err = ws.db.FindSnapshot(gameNumber, playerID, snapshotTime)
		i++
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Could not load snapshot for player %v at %v", playerID, snapshotTime)
			log.Printf("Could not load snapshot for match %v player %v, %v: %v", gameNumber, playerID, snapshotTime, err)
			return
		}
	}

	mergedSnapshot := opsec.Merge(loadedSnapshots...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mergedSnapshot)
}

func (ws *webServer) Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/matches/{gameNumber}", ws.ShowMatch)
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