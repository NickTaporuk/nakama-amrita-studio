package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	// UserPluginName
	UserURLName = "add_user"
	// TournamentKeyFromPayload use to get value for that parameter from payload
	TournamentKeyFromPayload = "tournament_id"
)

var (
	// ErrorNotFoundPayloadParameterTournamentID
	ErrorNotFoundPayloadParameterTournamentID = errors.New("parameter tournament_id is not found")
)

// All Go modules must have a InitModule function with this exact signature.
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// Register the RPC function.
	if err := initializer.RegisterRpc(UserURLName, AddUserToTournament); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	return nil
}

func AddUserToTournament(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	meta := make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &meta); err != nil {
		return "", err
	}

	if _, ok := meta[TournamentKeyFromPayload]; !ok {
		return "", ErrorNotFoundPayloadParameterTournamentID
	}
	tournamentName := meta[TournamentKeyFromPayload].(string)
	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	username := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)

	err := nk.TournamentJoin(ctx, tournamentName, userID, username)
	if err != nil {
		return "", err
	}

	return "Success", nil
}
