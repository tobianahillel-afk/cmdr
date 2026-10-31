package workqueue

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

const queuePageSize = 50
const maxQueueCursorBytes = 256

var ErrCursor = errors.New("work queue page no longer available")

// A cursor is an opaque presentation hint over the CURRENTLY authorized,
// ordered and filtered queue. It is NOT a credential or a signed claim.
// Every request must read and authorize the entire bounded owner result again.
type queuePageCursor struct {
	Version int    `json:"v"`
	Start   int    `json:"start"`
	Digest  string `json:"digest"`
}

// Decode only a bounded first-party page token. No arbitrary client offsets,
// query language, raw owner IDs or cross-tenant context are accepted.
func decodeQueueCursor(raw string) (queuePageCursor, error) {
	var cursor queuePageCursor
	if raw == "" {
		return cursor, nil
	}
	if len(raw) > maxQueueCursorBytes || strings.ContainsAny(raw, "=+/") {
		return cursor, ErrCursor
	}
	bytes, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(bytes) > 180 {
		return cursor, ErrCursor
	}
	dec := json.NewDecoder(strings.NewReader(string(bytes)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cursor); err != nil {
		return queuePageCursor{}, ErrCursor
	}
	var trailing any
	if dec.Decode(&trailing) != io.EOF {
		return queuePageCursor{}, ErrCursor
	}
	digest, err := hex.DecodeString(cursor.Digest)
	if err != nil || len(digest) != sha256.Size || cursor.Version != 1 ||
		cursor.Start < queuePageSize || cursor.Start > 5000 ||
		cursor.Start % queuePageSize != 0 {
		return queuePageCursor{}, ErrCursor
	}
	return cursor, nil
}

func authorizedQueuePage(req Request, all []Row) ([]Row, string, string, bool, string, error) {
	cursor, err := decodeQueueCursor(req.Cursor)
	if err != nil {
		return nil, "", "", false, "", err
	}
	// A hash is a revision check only. Recomputing it over all currently
	// authorized rows prevents stale/revoked pagination from rendering any row.
	// The selected principal/scope/view/filter/order are part of the digest.
	scope, err := json.Marshal(struct {
		Principal string `json:"principal"`
		Tenant string `json:"tenant"`
		Environment string `json:"environment"`
		View View `json:"view"`
		Status string `json:"status"`
		Owner string `json:"owner"`
		Group string `json:"group"`
		Sort string `json:"sort"`
		Rows []Row `json:"rows"`
	}{req.Principal, req.Tenant, req.Environment, req.View, req.StatusFilter, req.OwnerFilter, req.GroupBy, req.SortOrder, all})
	if err != nil {
		return nil, "", "", false, "", ErrIncomplete
	}
	digest := sha256.Sum256(scope)
	currentDigest := hex.EncodeToString(digest[:])
	start := cursor.Start
	if req.Cursor != "" && (cursor.Digest != currentDigest || start >= len(all)) {
		return nil, "", "", false, "", ErrCursor
	}
	// A returning Incident is an untrusted *location hint*, not an ACL.
	// Find it only inside the complete freshly CanRead-authorized cohort.
	// Explicit cursor requests never enter this branch.
	if req.FocusSelection != "" {
		index := -1
		for i, row := range all {
			if string(row.Kind)+":"+row.ID == req.FocusSelection {
				index = i
				break
			}
		}
		if index < 0 {
			return nil, "", "", false, "", ErrCursor
		}
		start = (index / queuePageSize) * queuePageSize
	}
	current := req.Cursor
	if req.FocusSelection != "" && start > 0 {
		current, err = encodeQueueCursor(start, currentDigest)
		if err != nil {
			return nil, "", "", false, "", ErrIncomplete
		}
	}
	stop := start + queuePageSize
	if stop > len(all) {
		stop = len(all)
	}
	page := all[start:stop]
	next := ""
	if stop < len(all) {
		next, err = encodeQueueCursor(stop, currentDigest)
		if err != nil { return nil, "", "", false, "", ErrIncomplete }
	}
	// A previous-page token is ONLY a navigation hint. Page 2 goes back
	// directly to the original first-page URL (no token); later pages carry
	// an earlier offset over the same freshly authorized projection.
	previous := ""
	if start > queuePageSize {
		previous, err = encodeQueueCursor(start-queuePageSize, currentDigest)
		if err != nil { return nil, "", "", false, "", ErrIncomplete }
	}
	return page, next, previous, start > 0, current, nil
}

func encodeQueueCursor(start int, digest string) (string, error) {
	payload, err := json.Marshal(queuePageCursor{Version: 1, Start: start, Digest: digest})
	if err != nil { return "", err }
	return base64.RawURLEncoding.EncodeToString(payload), nil
}
