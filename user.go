package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	// UserEntryPointURLName
	UserEntryPointURLName = "addUserToTournament"
	// TournamentKeyFromPayload use to get value for that parameter from payload
	TournamentKeyFromPayload = "tournament_id"
)

var (
	// ErrorNotFoundPayloadParameterTournamentID
	ErrorNotFoundPayloadParameterTournamentID = errors.New("parameter tournament_id is not found")
	// ErrorPayloadEmpty
	ErrorPayloadEmpty = errors.New("parameters are empty")
)

// All Go modules must have a InitModule function with this exact signature.
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// Register the RPC function.
	if err := initializer.RegisterRpc(UserEntryPointURLName, AddUserToTournament); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	return nil
}

// AddUserToTournament
func AddUserToTournament(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	if payload == "" {
		return "", ErrorPayloadEmpty
	}

	meta := make(map[string]interface{})

	if err := json.Unmarshal([]byte(payload), &meta); err != nil {
		return "", err
	}

	if len(meta) == 0 {
		return "", ErrorPayloadEmpty
	}

	if _, ok := meta[TournamentKeyFromPayload]; !ok {
		return "", ErrorNotFoundPayloadParameterTournamentID
	}

	tournamentID := meta[TournamentKeyFromPayload].(string)
	if tournamentID == "" {
		return "", ErrorPayloadEmpty
	}

	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	username := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)

	if !ExistTournamentByID(ctx, logger, nk, tournamentID) {
		err := CreateTournament(ctx, logger, nk, tournamentID)
		if err != nil {
			return "", err
		}
	}

	err := nk.TournamentJoin(ctx, tournamentID, userID, username)
	if err != nil {
		return "", err
	}

	return "Success", nil
}

// ExistTournament check exist particular tournament by id
func ExistTournamentByID(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, id string) bool {

	categoryStart := 1
	categoryEnd := 2
	startTime := int(time.Now().Unix())
	endTime := 0 // All tournaments from the start time.
	limit := 1   // Number to list per page.
	cursor := ""

	if tournaments, err := nk.TournamentList(ctx, categoryStart, categoryEnd, startTime, endTime, limit, cursor); err != nil {
		return false
	} else {
		for _, t := range tournaments.Tournaments {
			if t.GetId() == id {
				return true
			}
		}
	}

	return false
}

// CreateTournament
func CreateTournament(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, id string) error {
	sortOrder := "desc"           // One of: "desc", "asc".
	operator := "best"            // One of: "best", "set", "incr".
	resetSchedule := "0 12 * * *" // Noon UTC each day.
	metadata := map[string]interface{}{}

	title := "Daily Dash"
	description := "Dash past your opponents for high scores and big rewards!"
	category := 1
	startTime := 0       // Start now.
	endTime := 0         // Never end, repeat the tournament each day forever.
	duration := 3600     // In seconds.
	maxSize := 10000     // First 10,000 players who join.
	maxNumScore := 3     // Each player can have 3 attempts to score.
	joinRequired := true // Must join to compete.
	if err := nk.TournamentCreate(ctx, id, sortOrder, operator, resetSchedule, metadata, title, description, category, startTime, endTime, duration, maxSize, maxNumScore, joinRequired); err != nil {
		return err
	}

	return nil
}
