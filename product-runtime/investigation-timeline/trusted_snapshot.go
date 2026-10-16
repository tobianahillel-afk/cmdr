package investigationtimeline

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

// MaxTrustedSnapshotBytes bounds a host-supplied, read-only owner snapshot.
// This is not a persistent schema, a remote provider or an ACL service.
const MaxTrustedSnapshotBytes = 8 << 20

// DecodeTrustedMemoryTimelineJSON reads a trusted, local Case/Timeline owner
// snapshot into the existing memory provider data structure. The host must
// protect the file and authenticate the principal independently.
// Construction still validates object ownership and reader relationships.
func DecodeTrustedMemoryTimelineJSON(reader io.Reader) (MemoryTimelineData, error) {
	var data MemoryTimelineData
	if reader == nil {
		return data, errors.New("trusted snapshot reader required")
	}
	raw, err := io.ReadAll(io.LimitReader(reader, MaxTrustedSnapshotBytes+1))
	if err != nil {
		return data, errors.New("could not read trusted snapshot")
	}
	if len(raw) == 0 || len(raw) > MaxTrustedSnapshotBytes {
		return data, errors.New("trusted snapshot missing or exceeds byte limit")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		return MemoryTimelineData{}, errors.New("invalid trusted snapshot JSON")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return MemoryTimelineData{}, errors.New("trusted snapshot contains trailing data")
	}
	if len(data.Principals) == 0 || len(data.Principals) > MaxEntries ||
		len(data.Cases) == 0 || len(data.Cases) > MaxEntries ||
		len(data.Entries) > MaxEntries || len(data.Sources) > MaxEntries {
		return MemoryTimelineData{}, errors.New("trusted snapshot record count invalid")
	}
	return data, nil
}
