// Command case-search-workbench runs a read-only local search job over an
// operator-supplied synthetic Case owner snapshot. It is NOT a live SIEM.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func main() {
	if err := runCLI(context.Background(), os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Only a trusted operator-local, bounded regular file is accepted, never
// browser-submitted paths or a network data source.
func readSnapshot(path string) (timeline.MemoryTimelineData, error) {
	if !filepath.IsAbs(path) {
		return timeline.MemoryTimelineData{}, errors.New("absolute local owner snapshot required")
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > 8<<20 {
		return timeline.MemoryTimelineData{}, errors.New("invalid local owner snapshot")
	}
	file, err := os.Open(path)
	if err != nil {
		return timeline.MemoryTimelineData{}, errors.New("cannot open owner snapshot")
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return timeline.MemoryTimelineData{}, errors.New("owner snapshot changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(file, 8<<20+1))
	if err != nil || len(data) > 8<<20 {
		return timeline.MemoryTimelineData{}, errors.New("owner snapshot too large")
	}
	var owner timeline.MemoryTimelineData
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&owner); err != nil {
		return timeline.MemoryTimelineData{}, errors.New("invalid owner snapshot JSON")
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		return timeline.MemoryTimelineData{}, errors.New("owner snapshot must contain one object")
	}
	return owner, nil
}

func runCLI(ctx context.Context, args []string, out io.Writer) error {
	flags := flag.NewFlagSet("case-search-workbench", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	path := flags.String("snapshot", "", "absolute read-only synthetic owner snapshot JSON")
	principal := flags.String("principal", "", "local synthetic principal from the snapshot")
	caseID := flags.String("case", "", "existing authorized Case ID")
	query := flags.String("query", "", "text on released Event ID, source ref or provenance")
	from := flags.String("from", "", "inclusive RFC3339 search start (UTC or offset)")
	to := flags.String("to", "", "exclusive RFC3339 search end (UTC or offset)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *path == "" || *principal == "" || *caseID == "" ||
		*query == "" || *from == "" || *to == "" {
		return errors.New("require snapshot, principal, case, query, from, to")
	}
	if _, err := time.Parse(time.RFC3339, *from); err != nil {
		return errors.New("invalid search start time")
	}
	if _, err := time.Parse(time.RFC3339, *to); err != nil {
		return errors.New("invalid search end time")
	}
	data, err := readSnapshot(*path)
	if err != nil {
		return err
	}
	report, err := runCaseSearch(ctx, data, searchOptions{
		Principal: *principal, CaseID: *caseID, Query: *query, From: *from, To: *to,
	})
	if err != nil {
		return err
	}
	// An unsuccessful authorization or job never writes partial results.
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(report)
}
