package investigationtimeline

// maxTimelinePageSize limits HTTP response size, not the bounded in-process
// projection. Every request must still be authorized against the full Case,
// Timeline Entry and source owner chain before page selection.
const maxTimelinePageSize = 200

type timelineHTTPPage struct {
	Offset       int  `json:"offset"`
	Limit        int  `json:"limit"`
	TotalVisible int  `json:"total_visible"`
	HasMore      bool `json:"has_more"`
}

type timelinePagedHTTPResult struct {
	Result
	Page timelineHTTPPage `json:"page"`
}

// paginateAuthorizedTimeline is only called after OpenCaseTimeline has
// enforced all Case, Timeline Entry and linked source permissions. TotalVisible
// never includes inaccessible objects. Pages are live, not transactionally
// snapshotted across requests; consumers must not assume stable cursors.
func paginateAuthorizedTimeline(in Result, selection timelineReadSelection) timelinePagedHTTPResult {
	total := len(in.Entries)
	start := min(selection.offset, total)
	end := min(start+selection.limit, total)
	page := Result{Status: in.Status, Entries: append([]ProjectedEntry{}, in.Entries[start:end]...)}
	return timelinePagedHTTPResult{
		Result: page,
		Page: timelineHTTPPage{
			Offset: selection.offset, Limit: selection.limit,
			TotalVisible: total, HasMore: end < total,
		},
	}
}
