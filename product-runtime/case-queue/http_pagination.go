package casequeue

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash"
	"io"
	"strings"
)

// Case Queue cursors are presentation-only offsets over a freshly
// authorized, fully projected list. They do not grant Case or Saved View
// permissions and are not durable snapshot identifiers.
type queuePageCursor struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Digest string `json:"digest"`
}

type queuePageResponse struct {
	*Projection
	NextCursor string `json:"next_cursor,omitempty"`
}

func queueProjectionDigest(p *Projection) (string, error) {
	var h hash.Hash = sha256.New()
	if err := json.NewEncoder(h).Encode(p); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Re-evaluate the WHOLE queue before paging. The digest binds each cursor to
// the exact currently readable projection, filter/sort output, scope and
// Shared Saved View effects; any change invalidates the previous cursor.
// A forged cursor can at most change an offset within the current authorized
// rows; neither its offset nor digest is treated as authority.
func paginateQueueProjection(p *Projection, limit int, opaque string) (queuePageResponse, error) {
	if p == nil || limit < 1 || limit > 100 {
		return queuePageResponse{}, errors.New("invalid Case Queue page")
	}
	digest, err := queueProjectionDigest(p)
	if err != nil {
		return queuePageResponse{}, err
	}
	start := 0
	if opaque != "" {
		if len(opaque) > 256 {
			return queuePageResponse{}, errors.New("invalid Case Queue cursor")
		}
		raw, err := base64.RawURLEncoding.DecodeString(opaque)
		if err != nil || len(raw) > 192 {
			return queuePageResponse{}, errors.New("invalid Case Queue cursor")
		}
		decoder := json.NewDecoder(strings.NewReader(string(raw)))
		decoder.DisallowUnknownFields()
		var saved queuePageCursor
		if decoder.Decode(&saved) != nil || decoder.Decode(new(any)) != io.EOF ||
			saved.Limit != limit || saved.Offset < 1 || saved.Offset >= len(p.Rows) ||
			saved.Digest != digest {
			return queuePageResponse{}, errors.New("stale or invalid Case Queue cursor")
		}
		start = saved.Offset
	}
	end := start + limit
	if end > len(p.Rows) {
		end = len(p.Rows)
	}
	// Clone ONLY the response header, never mutate the canonical projection.
	page := *p
	page.Rows = append([]Row{}, p.Rows[start:end]...)
	out := queuePageResponse{Projection: &page}
	if end < len(p.Rows) {
		next, err := json.Marshal(queuePageCursor{Offset: end, Limit: limit, Digest: digest})
		if err != nil {
			return queuePageResponse{}, err
		}
		out.NextCursor = base64.RawURLEncoding.EncodeToString(next)
	}
	return out, nil
}
