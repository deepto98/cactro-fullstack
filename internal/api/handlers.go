package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deepto98/cargo-fullstack/internal/models"
	"github.com/gorilla/mux"
)

// Handler holds dependencies for API endpoints.
type Handler struct {
	DB *sql.DB
}

// CreatePoll handles POST /api/polls.
// It expects a JSON body with a question and an array of option strings.
func (h *Handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.Question == "" || len(payload.Options) < 2 {
		http.Error(w, "Question and at least two options are required", http.StatusBadRequest)
		return
	}

	// Insert poll.
	var pollID int
	err := h.DB.QueryRow("INSERT INTO polls(question) VALUES($1) RETURNING id", payload.Question).Scan(&pollID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create poll: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert each option.
	for _, opt := range payload.Options {
		_, err := h.DB.Exec("INSERT INTO poll_options(poll_id, option_text) VALUES($1, $2)", pollID, opt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create poll options: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"poll_id": pollID})
}

// GetPoll handles GET /api/polls/{id}.
// It returns the poll details along with options and vote counts.
func (h *Handler) GetPoll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	// Fetch poll question.
	var poll models.Poll
	err = h.DB.QueryRow("SELECT id, question FROM polls WHERE id = $1", pollID).Scan(&poll.ID, &poll.Question)
	if err != nil {
		http.Error(w, "Poll not found", http.StatusNotFound)
		return
	}

	// Fetch options.
	rows, err := h.DB.Query("SELECT id, option_text, vote_count FROM poll_options WHERE poll_id = $1", pollID)
	if err != nil {
		http.Error(w, "Failed to fetch options", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var opt models.Option
		if err := rows.Scan(&opt.ID, &opt.Text, &opt.VoteCount); err != nil {
			http.Error(w, "Error scanning options", http.StatusInternalServerError)
			return
		}
		poll.Options = append(poll.Options, opt)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(poll)
}

// Vote handles POST /api/polls/{id}/vote.
// It expects a JSON body with the option_id to vote for.
func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		OptionID int `json:"option_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update vote count.
	res, err := h.DB.Exec("UPDATE poll_options SET vote_count = vote_count + 1 WHERE id = $1 AND poll_id = $2", payload.OptionID, pollID)
	if err != nil {
		http.Error(w, "Failed to register vote", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Invalid option or poll", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vote registered"})
}

// ListPolls handles GET /api/polls and returns a list of all polls.
func (h *Handler) ListPolls(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, question FROM polls")
	if err != nil {
		http.Error(w, "Failed to query polls", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PollSummary struct {
		ID       int    `json:"id"`
		Question string `json:"question"`
	}

	var polls []PollSummary
	for rows.Next() {
		var p PollSummary
		if err := rows.Scan(&p.ID, &p.Question); err != nil {
			http.Error(w, "Failed to scan poll", http.StatusInternalServerError)
			return
		}
		polls = append(polls, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(polls)
}
