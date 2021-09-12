package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.albinodrought.com/neptunes-pride/internal/matchstore"
	"go.albinodrought.com/neptunes-pride/internal/npapi"
)

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

	match, err := ws.db.FindMatchOrFail(vars["gameNumber"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Match not found"))
		log.Printf("Match %v not found: %v", vars, err)
		return
	}

	for i, creds := range match.PlayerCreds {
		creds.APIKey = ""
		match.PlayerCreds[i] = creds
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func (ws *webServer) Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/matches/{gameNumber}", ws.ShowMatch)

	return r
}
