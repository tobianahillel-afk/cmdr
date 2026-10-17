package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"
)

// This explicit local receipt is sensitive: each returned event was already
// authorized by the existing independent Case/Timeline/typed-source handlers.
// Matching a time window never establishes causality or shared source identity.
type localCaseWindowSequence struct {
	Format      string           `json:"format"`
	Synthetic   bool             `json:"synthetic"`
	FirstCase   string           `json:"first_case"`
	SecondCase  string           `json:"second_case"`
	StartUTC    string           `json:"start_utc"`
	EndUTC      string           `json:"end_utc_exclusive"`
	FirstEvents int              `json:"first_events"`
	SecondEvents int             `json:"second_events"`
	Events      []localCaseEvent `json:"events"`
}

func parseCertainFiveMinuteWindow(value string) (time.Time, error) {
	start, err := time.Parse(time.RFC3339, value)
	if err != nil || start.Format(time.RFC3339) != value ||
		start.Location() != time.UTC || start.Second() != 0 ||
		start.Minute()%5 != 0 || start.Nanosecond() != 0 {
		return time.Time{}, errors.New("expected canonical UTC five-minute window start")
	}
	return start, nil
}

// No output is emitted until both owners have independently completed their
// full authorized paginated Case reads, including typed source checks.
func writeDemoCaseWindowSequence(ctx context.Context, principal, snapshot, firstID, secondID, window string, out io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || firstID == "" || secondID == "" || firstID == secondID ||
		len(firstID) > 1024 || len(secondID) > 1024 {
		return errors.New("two distinct authorized Cases and analyst required")
	}
	start, err := parseCertainFiveMinuteWindow(window)
	if err != nil {
		return err
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return err
	}
	full, err := collectAuthorizedCrossCaseSequence(ctx, handler, token, firstID, secondID)
	if err != nil {
		return err
	}
	end := start.Add(5 * time.Minute)
	report := localCaseWindowSequence{
		Format: "cmdr-local-cross-case-time-window-v1", Synthetic: true,
		FirstCase: firstID, SecondCase: secondID,
		StartUTC: start.Format(time.RFC3339), EndUTC: end.Format(time.RFC3339),
		Events: make([]localCaseEvent, 0),
	}
	for _, item := range full.Events {
		if err := ctx.Err(); err != nil {
			return err
		}
		if item.Entry.TimeUncertain {
			continue
		}
		instant, err := time.Parse(time.RFC3339Nano, item.Entry.SourceTime)
		if err != nil {
			return errors.New("invalid authorized certain timestamp")
		}
		if instant.Before(start) || !instant.Before(end) {
			continue
		}
		report.Events = append(report.Events, item)
		switch item.CaseID {
		case firstID:
			report.FirstEvents++
		case secondID:
			report.SecondEvents++
		default:
			return errors.New("invalid authorized Case in timeline")
		}
	}
	if report.FirstEvents == 0 || report.SecondEvents == 0 {
		return errors.New("no shared certain window in the two authorized Cases")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(report)
}
