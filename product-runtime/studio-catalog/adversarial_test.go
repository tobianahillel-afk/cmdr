package studiocatalog

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// =====================================================================
// 1. Cursor Manipulation and Pagination Boundary Conditions
// =====================================================================

func TestAdversarial_Cursor_CorruptedAndTamperedTokens(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Baseline browse to ensure valid data
	baseResp, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Limit:     2,
	})
	if err != nil {
		t.Fatalf("setup browse failed: %v", err)
	}
	validCursor := baseResp.Pagination.CursorNext
	if validCursor == "" {
		t.Fatalf("expected valid cursor from browse")
	}

	testCases := []struct {
		name        string
		cursor      string
		expectedErr ErrorCode
	}{
		{
			name:        "Garbage non-base64",
			cursor:      "!!not-base64-at-all??$$#%",
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Standard base64 padding rejected by RawURL",
			cursor:      base64.StdEncoding.EncodeToString([]byte(`{"v":1,"t":"tenant-alpha","sf":"updated_at","sd":"desc","id":"skl-01"}`)) + "==",
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Valid base64 but invalid JSON",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{not valid json`)),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Valid JSON with invalid schema version (v=2)",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{"v":2,"t":"tenant-alpha","sf":"updated_at","sd":"desc","id":"skl-01"}`)),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Valid JSON with invalid schema version (v=0)",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{"v":0,"t":"tenant-alpha","sf":"updated_at","sd":"desc","id":"skl-01"}`)),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Valid JSON with negative version (v=-1)",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{"v":-1,"t":"tenant-alpha","sf":"updated_at","sd":"desc","id":"skl-01"}`)),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Missing tiebreaker ID in cursor",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{"v":1,"t":"tenant-alpha","sf":"updated_at","sd":"desc","id":""}`)),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Null tiebreaker ID in cursor",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{"v":1,"t":"tenant-alpha","sf":"updated_at","sd":"desc"}`)),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Wrong data types in JSON (v as string)",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{"v":"one","t":"tenant-alpha","id":"skl-01"}`)),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Massive cursor payload (buffer overflow attempt 64KB)",
			cursor:      base64.RawURLEncoding.EncodeToString([]byte(`{"v":1,"t":"tenant-alpha","id":"skl-01","junk":"` + strings.Repeat("A", 65536) + `"}`)),
			expectedErr: ErrorCodeInvalidCursor, // Anchor won't match or invalid
		},
		{
			name: "Cross-tenant cursor replay (tenant-bravo cursor on tenant-alpha)",
			cursor: EncodeCursor(CursorPayload{
				Version:       1,
				TenantID:      "tenant-bravo",
				SortField:     "updated_at",
				SortDirection: "desc",
				TiebreakerID:  "skl-01",
			}),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name: "Tampered sort field (cursor=name, request=updated_at)",
			cursor: EncodeCursor(CursorPayload{
				Version:       1,
				TenantID:      "tenant-alpha",
				SortField:     "name",
				SortDirection: "desc",
				TiebreakerID:  "skl-01",
			}),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name: "Tampered sort direction (cursor=asc, request=desc)",
			cursor: EncodeCursor(CursorPayload{
				Version:       1,
				TenantID:      "tenant-alpha",
				SortField:     "updated_at",
				SortDirection: "asc",
				TiebreakerID:  "skl-01",
			}),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name: "Non-existent anchor ID (deleted or synthetic)",
			cursor: EncodeCursor(CursorPayload{
				Version:       1,
				TenantID:      "tenant-alpha",
				SortField:     "updated_at",
				SortDirection: "desc",
				TiebreakerID:  "ghost-asset-id-99999",
			}),
			expectedErr: ErrorCodeInvalidCursor,
		},
		{
			name:        "Null byte in cursor string",
			cursor:      validCursor + "\x00extra",
			expectedErr: ErrorCodeInvalidCursor,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Browse(ctx, auth, BrowseRequest{
				TenantRef: "tenant-alpha",
				Limit:     2,
				Cursor:    tc.cursor,
			})
			if err == nil {
				t.Fatalf("expected error %s, got success", tc.expectedErr)
			}
			env, ok := err.(*ErrorEnvelope)
			if !ok {
				t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
			}
			if env.Code != tc.expectedErr {
				t.Fatalf("expected error code %s, got %s (message: %s)", tc.expectedErr, env.Code, env.Message)
			}
			if env.HTTPStatus != 400 {
				t.Fatalf("expected HTTP 400, got %d", env.HTTPStatus)
			}
		})
	}
}

func TestAdversarial_Pagination_LimitClampingAndOffsets(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	limitsToTest := []struct {
		requestedLimit int
		expectedLimit  int
	}{
		{requestedLimit: 0, expectedLimit: 25},
		{requestedLimit: -1, expectedLimit: 25},
		{requestedLimit: -999999, expectedLimit: 25},
		{requestedLimit: math.MinInt, expectedLimit: 25},
		{requestedLimit: 1, expectedLimit: 1},
		{requestedLimit: 25, expectedLimit: 25},
		{requestedLimit: 50, expectedLimit: 50},
		{requestedLimit: 100, expectedLimit: 100},
		{requestedLimit: 101, expectedLimit: 100},
		{requestedLimit: 10000, expectedLimit: 10000 - 9900}, // 100
		{requestedLimit: math.MaxInt, expectedLimit: 100},
	}

	for _, tc := range limitsToTest {
		t.Run(fmt.Sprintf("Limit_%d", tc.requestedLimit), func(t *testing.T) {
			resp, err := svc.Browse(ctx, auth, BrowseRequest{
				TenantRef: "tenant-alpha",
				Limit:     tc.requestedLimit,
			})
			if err != nil {
				t.Fatalf("browse with limit %d failed: %v", tc.requestedLimit, err)
			}
			if resp.Pagination.Limit != tc.expectedLimit {
				t.Fatalf("expected limit %d, got %d", tc.expectedLimit, resp.Pagination.Limit)
			}
			if len(resp.Items) > tc.expectedLimit {
				t.Fatalf("returned %d items exceeding clamped limit %d", len(resp.Items), tc.expectedLimit)
			}
		})
	}
}

func TestAdversarial_Pagination_EmptyTenantAndTerminalBoundaries(t *testing.T) {
	store := NewInMemoryStore()
	svc := NewCatalogService(store)
	ctx := context.Background()
	auth := adminAuth("tenant-empty")

	// 1. Empty tenant browse
	respEmpty, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-empty",
		Limit:     10,
	})
	if err != nil {
		t.Fatalf("browse on empty tenant failed: %v", err)
	}
	if len(respEmpty.Items) != 0 {
		t.Fatalf("expected 0 items on empty tenant, got %d", len(respEmpty.Items))
	}
	if respEmpty.Pagination.HasMore {
		t.Fatalf("expected HasMore to be false on empty tenant")
	}
	if respEmpty.Pagination.CursorNext != "" {
		t.Fatalf("expected empty CursorNext on empty tenant, got %q", respEmpty.Pagination.CursorNext)
	}
	if respEmpty.Pagination.TotalCount != 0 {
		t.Fatalf("expected TotalCount 0, got %d", respEmpty.Pagination.TotalCount)
	}

	// 2. Add exactly 1 asset to tenant
	singleAsset := makeTestAsset("tenant-empty", "skl-single", AssetKindSkill, "single-skill", LifecycleStatePublished, RiskClassificationLow)
	if err := store.Put(ctx, singleAsset); err != nil {
		t.Fatalf("put failed: %v", err)
	}

	// Browse with limit 1
	resp1, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-empty",
		Limit:     1,
	})
	if err != nil {
		t.Fatalf("browse single failed: %v", err)
	}
	if len(resp1.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp1.Items))
	}
	if resp1.Pagination.HasMore {
		t.Fatalf("expected HasMore false when all items fetched on page 1")
	}
	if resp1.Pagination.CursorNext != "" {
		t.Fatalf("expected CursorNext empty when HasMore false, got %q", resp1.Pagination.CursorNext)
	}

	// 3. Add a second asset
	secondAsset := makeTestAsset("tenant-empty", "skl-second", AssetKindSkill, "second-skill", LifecycleStatePublished, RiskClassificationLow)
	if err := store.Put(ctx, secondAsset); err != nil {
		t.Fatalf("put failed: %v", err)
	}

	// Fetch page 1 (limit 1)
	respP1, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-empty",
		Limit:     1,
	})
	if err != nil {
		t.Fatalf("page 1 failed: %v", err)
	}
	if len(respP1.Items) != 1 || !respP1.Pagination.HasMore || respP1.Pagination.CursorNext == "" {
		t.Fatalf("page 1 invalid: len=%d, hasMore=%v, cursorNext=%q", len(respP1.Items), respP1.Pagination.HasMore, respP1.Pagination.CursorNext)
	}

	// Fetch page 2 (limit 1 with cursor)
	respP2, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-empty",
		Limit:     1,
		Cursor:    respP1.Pagination.CursorNext,
	})
	if err != nil {
		t.Fatalf("page 2 failed: %v", err)
	}
	if len(respP2.Items) != 1 {
		t.Fatalf("expected 1 item on page 2, got %d", len(respP2.Items))
	}
	if respP2.Pagination.HasMore {
		t.Fatalf("expected HasMore false on final page")
	}
	if respP2.Pagination.CursorNext != "" {
		t.Fatalf("expected empty CursorNext on final page, got %q", respP2.Pagination.CursorNext)
	}
	if respP1.Items[0].ID == respP2.Items[0].ID {
		t.Fatalf("duplicate item across pages: %s", respP1.Items[0].ID)
	}
}

func TestAdversarial_Pagination_HighVolumeStress(t *testing.T) {
	store := NewInMemoryStore()
	svc := NewCatalogService(store)
	ctx := context.Background()
	auth := adminAuth("tenant-high-volume")

	// Seed 500 assets with diverse kinds, timestamps, and IDs
	totalAssets := 500
	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < totalAssets; i++ {
		kind := AllAssetKinds[i%len(AllAssetKinds)]
		a := makeTestAsset(
			"tenant-high-volume",
			fmt.Sprintf("asset-%05d", i),
			kind,
			fmt.Sprintf("Asset Name %05d", i),
			LifecycleStatePublished,
			RiskClassificationMedium,
		)
		// Deterministic staggered timestamps
		a.UpdatedAt = baseTime.Add(time.Duration(i) * time.Minute)
		a.CreatedAt = a.UpdatedAt
		if err := store.Put(ctx, a); err != nil {
			t.Fatalf("failed to seed asset %d: %v", i, err)
		}
	}

	// Paginate through all 500 items using page limit 37 (prime number to test off-boundary offsets)
	pageSize := 37
	var collectedIDs []string
	seenIDs := make(map[string]bool)
	cursor := ""
	pageCount := 0

	for {
		pageCount++
		if pageCount > 50 {
			t.Fatalf("infinite pagination loop detected! page count: %d", pageCount)
		}

		resp, err := svc.Browse(ctx, auth, BrowseRequest{
			TenantRef:     "tenant-high-volume",
			Limit:         pageSize,
			Cursor:        cursor,
			SortField:     "updated_at",
			SortDirection: "desc",
		})
		if err != nil {
			t.Fatalf("page %d failed: %v", pageCount, err)
		}

		for _, it := range resp.Items {
			if seenIDs[it.ID] {
				t.Fatalf("duplicate item encountered on page %d: %s", pageCount, it.ID)
			}
			seenIDs[it.ID] = true
			collectedIDs = append(collectedIDs, it.ID)
		}

		if !resp.Pagination.HasMore {
			break
		}
		if resp.Pagination.CursorNext == "" {
			t.Fatalf("page %d reported HasMore=true but CursorNext is empty", pageCount)
		}
		cursor = resp.Pagination.CursorNext
	}

	if len(collectedIDs) != totalAssets {
		t.Fatalf("expected %d total collected assets, got %d", totalAssets, len(collectedIDs))
	}

	// Verify exact descending updated_at order (with id ASC tiebreaker)
	// Because asset-00499 has the latest updated_at, it must appear first!
	if collectedIDs[0] != "asset-00499" {
		t.Fatalf("expected first asset to be asset-00499, got %s", collectedIDs[0])
	}
	if collectedIDs[len(collectedIDs)-1] != "asset-00000" {
		t.Fatalf("expected last asset to be asset-00000, got %s", collectedIDs[len(collectedIDs)-1])
	}
}

// =====================================================================
// 2. Concurrent Operations & Thread Safety
// =====================================================================

func TestAdversarial_InMemoryStore_ConcurrentReadWriteStress(t *testing.T) {
	store := NewInMemoryStore()
	svc := NewCatalogService(store)
	ctx := context.Background()

	// Initial seed
	for i := 0; i < 20; i++ {
		a := makeTestAsset(
			"tenant-stress",
			fmt.Sprintf("initial-%02d", i),
			AssetKindSkill,
			fmt.Sprintf("Initial Skill %02d", i),
			LifecycleStatePublished,
			RiskClassificationLow,
		)
		if err := store.Put(ctx, a); err != nil {
			t.Fatalf("initial put failed: %v", err)
		}
	}

	auth := adminAuth("tenant-stress")
	var wg sync.WaitGroup
	workers := 40
	iterations := 100
	var errorCount int64

	// Concurrently perform Put, Delete, Get, List, Browse, Search, Filter, Inspect, Compare
	for w := 0; w < workers; w++ {
		workerID := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				op := (workerID*iterations + i) % 9
				assetID := fmt.Sprintf("asset-%02d", (workerID+i)%30)

				switch op {
				case 0: // Put
					a := makeTestAsset("tenant-stress", assetID, AssetKindWorkflow, fmt.Sprintf("Wfl %d", i), LifecycleStatePublished, RiskClassificationMedium)
					if err := store.Put(ctx, a); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 1: // Get
					if _, err := store.Get(ctx, "tenant-stress", assetID); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 2: // Delete
					if err := store.Delete(ctx, "tenant-stress", assetID); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 3: // List
					if _, err := store.List(ctx, "tenant-stress"); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 4: // Browse
					if _, err := svc.Browse(ctx, auth, BrowseRequest{
						TenantRef: "tenant-stress",
						Limit:     10,
					}); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 5: // Search
					if _, err := svc.Search(ctx, auth, SearchRequest{
						TenantRef: "tenant-stress",
						QueryText: "Skill",
						Limit:     5,
					}); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 6: // Filter
					if _, err := svc.Filter(ctx, auth, FilterRequest{
						TenantRef: "tenant-stress",
						Criteria: FilterCriteria{
							Kinds: []AssetKind{AssetKindSkill, AssetKindWorkflow},
						},
						Limit: 5,
					}); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 7: // Inspect
					if _, err := svc.Inspect(ctx, auth, InspectRequest{
						TenantRef: "tenant-stress",
						AssetID:   "initial-00",
					}); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 8: // ClonePreparation
					if _, err := svc.ClonePreparation(ctx, auth, ClonePreparationRequest{
						TenantRef:       "tenant-stress",
						TargetTenantRef: "tenant-stress",
						SourceAssetID:   "initial-00",
					}); err != nil {
						atomic.AddInt64(&errorCount, 1)
					}
				}
			}
		}()
	}

	wg.Wait()
	if errs := atomic.LoadInt64(&errorCount); errs > 0 {
		t.Fatalf("encountered %d unexpected errors during concurrent stress", errs)
	}
}

func TestAdversarial_InMemoryStore_ReturnedCopyMutationSafety(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	original := makeTestAsset("tenant-safety", "skl-immutable", AssetKindSkill, "Original Name", LifecycleStatePublished, RiskClassificationLow)
	original.Tags = []string{"tag1", "tag2"}
	original.Parameters = []Parameter{{Name: "param1", Type: "string"}}

	if err := store.Put(ctx, original); err != nil {
		t.Fatalf("put failed: %v", err)
	}

	// 1. Get copy and mutate in-place
	retrieved1, err := store.Get(ctx, "tenant-safety", "skl-immutable")
	if err != nil || retrieved1 == nil {
		t.Fatalf("get failed: %v", err)
	}

	retrieved1.Name = "MUTATED NAME"
	retrieved1.Tags[0] = "MUTATED TAG"
	retrieved1.Tags = append(retrieved1.Tags, "EXTRA TAG")
	retrieved1.Parameters[0].Name = "MUTATED PARAM"
	retrieved1.Parameters = append(retrieved1.Parameters, Parameter{Name: "EXTRA PARAM"})

	// 2. Get another copy and assert it was completely untouched
	retrieved2, err := store.Get(ctx, "tenant-safety", "skl-immutable")
	if err != nil || retrieved2 == nil {
		t.Fatalf("second get failed: %v", err)
	}

	if retrieved2.Name != "Original Name" {
		t.Fatalf("store mutation leak: expected 'Original Name', got %q", retrieved2.Name)
	}
	if retrieved2.Tags[0] != "tag1" || len(retrieved2.Tags) != 2 {
		t.Fatalf("store tag mutation leak: got %v", retrieved2.Tags)
	}
	if retrieved2.Parameters[0].Name != "param1" || len(retrieved2.Parameters) != 1 {
		t.Fatalf("store parameter mutation leak: got %v", retrieved2.Parameters)
	}
}

// =====================================================================
// 3. Search and Filtering under Adversarial Inputs
// =====================================================================

func TestAdversarial_Search_AdversarialInputs(t *testing.T) {
	store := NewInMemoryStore()
	svc := NewCatalogService(store)
	ctx := context.Background()
	auth := adminAuth("tenant-adversarial")

	// Seed unicode, emoji, special chars, and normal assets
	assets := []*Asset{
		makeTestAsset("tenant-adversarial", "uni-01", AssetKindSkill, "威胁情报分析", LifecycleStatePublished, RiskClassificationHigh),
		makeTestAsset("tenant-adversarial", "uni-02", AssetKindSkill, "تحليل التهديدات السيبرانية", LifecycleStatePublished, RiskClassificationHigh),
		makeTestAsset("tenant-adversarial", "uni-03", AssetKindSkill, "ניתוח איומי סייבר", LifecycleStatePublished, RiskClassificationHigh),
		makeTestAsset("tenant-adversarial", "uni-04", AssetKindSkill, "incident-über-triage", LifecycleStatePublished, RiskClassificationLow),
		makeTestAsset("tenant-adversarial", "emo-01", AssetKindSkill, "fire-alert-🔥-shield-🛡️", LifecycleStatePublished, RiskClassificationCritical),
		makeTestAsset("tenant-adversarial", "reg-01", AssetKindSkill, "regex-[a-z0-9]+-test.*", LifecycleStatePublished, RiskClassificationMedium),
		makeTestAsset("tenant-adversarial", "sql-01", AssetKindSkill, "drop-table-';--injection", LifecycleStatePublished, RiskClassificationLow),
	}
	for _, a := range assets {
		if err := store.Put(ctx, a); err != nil {
			t.Fatalf("put failed: %v", err)
		}
	}

	// 1. Empty / whitespace search rejected
	emptyQueries := []string{"", "   ", "\t", "\n", "\r\n", "   \t\n   "}
	for _, eq := range emptyQueries {
		_, err := svc.Search(ctx, auth, SearchRequest{
			TenantRef: "tenant-adversarial",
			QueryText: eq,
		})
		if err == nil {
			t.Fatalf("expected error for empty search query %q, got success", eq)
		}
		env := err.(*ErrorEnvelope)
		if env.Code != ErrorCodeInvalidFilterParameter || env.HTTPStatus != 400 {
			t.Fatalf("expected 400 invalid-filter-parameter, got %d %s", env.HTTPStatus, env.Code)
		}
	}

	// 2. Regex injection attempts (Must be matched as LITERAL substrings, no ReDoS or regex engine interpretation)
	regexAttempts := []struct {
		pattern       string
		expectedMatch string
	}{
		{pattern: ".*", expectedMatch: "reg-01"},
		{pattern: "[a-z0-9]+", expectedMatch: "reg-01"},
		{pattern: "(?i)alert", expectedMatch: ""},    // Literal search for "(?i)alert" will match nothing
		{pattern: "^.*$", expectedMatch: ""},          // Literal search for "^.*$"
		{pattern: "a{100000}", expectedMatch: ""},     // Catastrophic backtracking pattern
		{pattern: "([a-zA-Z]+)*", expectedMatch: ""}, // Nested quantifier ReDoS
	}

	for _, tc := range regexAttempts {
		t.Run("RegexLiteral_"+tc.pattern, func(t *testing.T) {
			resp, err := svc.Search(ctx, auth, SearchRequest{
				TenantRef: "tenant-adversarial",
				QueryText: tc.pattern,
			})
			if err != nil {
				t.Fatalf("search with regex %q failed: %v", tc.pattern, err)
			}
			if tc.expectedMatch != "" {
				found := false
				for _, it := range resp.Items {
					if it.ID == tc.expectedMatch {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected literal regex match for %s, got 0 matches", tc.expectedMatch)
				}
			} else {
				if len(resp.Items) != 0 {
					t.Fatalf("expected 0 matches for non-existent literal regex %q, got %d", tc.pattern, len(resp.Items))
				}
			}
		})
	}

	// 3. Unicode and Multilingual search
	unicodeTests := []struct {
		query         string
		expectedAsset string
	}{
		{query: "威胁情报", expectedAsset: "uni-01"},
		{query: "تحليل", expectedAsset: "uni-02"},
		{query: "ניתוח", expectedAsset: "uni-03"},
		{query: "über", expectedAsset: "uni-04"},
		{query: "ÜBER", expectedAsset: "uni-04"}, // Case insensitivity in unicode
		{query: "🔥", expectedAsset: "emo-01"},
		{query: "🛡️", expectedAsset: "emo-01"},
		{query: "';--", expectedAsset: "sql-01"},
	}

	for _, tc := range unicodeTests {
		t.Run("Unicode_"+tc.query, func(t *testing.T) {
			resp, err := svc.Search(ctx, auth, SearchRequest{
				TenantRef: "tenant-adversarial",
				QueryText: tc.query,
			})
			if err != nil {
				t.Fatalf("unicode search %q failed: %v", tc.query, err)
			}
			if len(resp.Items) == 0 || resp.Items[0].ID != tc.expectedAsset {
				t.Fatalf("expected match %s for query %q, got items: %v", tc.expectedAsset, tc.query, resp.Items)
			}
		})
	}

	// 4. Massive search query (50KB string)
	t.Run("MassiveQueryString", func(t *testing.T) {
		largeQuery := strings.Repeat("A", 50000)
		resp, err := svc.Search(ctx, auth, SearchRequest{
			TenantRef: "tenant-adversarial",
			QueryText: largeQuery,
		})
		if err != nil {
			t.Fatalf("large query failed: %v", err)
		}
		if len(resp.Items) != 0 {
			t.Fatalf("expected 0 matches for large synthetic query, got %d", len(resp.Items))
		}
	})
}

func TestAdversarial_Filter_AdversarialCriteria(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// 1. Non-existent / malicious AssetKind
	respKind, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			Kinds: []AssetKind{AssetKind("malicious_kind"), AssetKind("'; DROP TABLE--")},
		},
	})
	if err != nil {
		t.Fatalf("filter with invalid kind failed: %v", err)
	}
	if len(respKind.Items) != 0 {
		t.Fatalf("expected 0 items for unknown kind, got %d", len(respKind.Items))
	}

	// 2. Non-existent Lifecycle state
	respLife, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			LifecycleStates: []LifecycleState{LifecycleState("fake-state")},
		},
	})
	if err != nil {
		t.Fatalf("filter with invalid lifecycle state failed: %v", err)
	}
	if len(respLife.Items) != 0 {
		t.Fatalf("expected 0 items for fake lifecycle, got %d", len(respLife.Items))
	}

	// 3. Filter by unknown owner under REQ-OBJ-009
	respUnknownOwner, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			OwnerID: "unknown",
		},
	})
	if err != nil {
		t.Fatalf("filter by owner unknown failed: %v", err)
	}
	if len(respUnknownOwner.Items) != 1 || respUnknownOwner.Items[0].ID != "skl-no-owner" {
		t.Fatalf("expected skl-no-owner for owner=unknown, got %v", respUnknownOwner.Items)
	}

	// 4. Tags filter: conjunction with non-existent tag
	respTags, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			Tags: []string{"security", "non-existent-tag-999"},
		},
	})
	if err != nil {
		t.Fatalf("filter with tags failed: %v", err)
	}
	if len(respTags.Items) != 0 {
		t.Fatalf("expected 0 items for non-existent tag conjunction, got %d", len(respTags.Items))
	}

	// 5. Empty criteria: returns all visible items
	respAll, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria:  FilterCriteria{},
	})
	if err != nil {
		t.Fatalf("empty criteria filter failed: %v", err)
	}
	if len(respAll.Items) != 10 {
		t.Fatalf("expected 10 items for empty criteria, got %d", len(respAll.Items))
	}
}

// =====================================================================
// 4. Execution Verb Blocking under All Request Variations
// =====================================================================

func TestAdversarial_ExecutionVerbBlocking_AllVariations(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	verbs := []struct {
		name string
		call func(context.Context, AuthContext, any) error
	}{
		{name: "Invoke", call: svc.Invoke},
		{name: "Execute", call: svc.Execute},
		{name: "Run", call: svc.Run},
		{name: "Call", call: svc.Call},
		{name: "Trigger", call: svc.Trigger},
	}

	payloads := []struct {
		payloadName string
		payload     any
	}{
		{payloadName: "Nil payload", payload: nil},
		{payloadName: "Empty struct", payload: struct{}{}},
		{payloadName: "Command injection map", payload: map[string]any{"command": "rm -rf /", "target": "10.0.0.1", "action": "execute"}},
		{payloadName: "String script", payload: "#!/bin/bash\ncat /etc/shadow"},
		{payloadName: "Integer flag", payload: 42},
		{payloadName: "Slice of parameters", payload: []string{"param1", "param2"}},
		{payloadName: "Raw JSON payload", payload: json.RawMessage(`{"dangerous": true, "operation": "execute"}`)},
	}

	authContexts := []struct {
		authName string
		auth     AuthContext
	}{
		{authName: "Super Admin Auth", auth: adminAuth("tenant-alpha")},
		{authName: "Empty/Anonymous Auth", auth: AuthContext{}},
		{authName: "Cross-Tenant ATS Auth", auth: NewAuthContext("usr-mssp", "tenant-alpha", []string{"tenant-alpha", "tenant-bravo"}, PermStudioLibraryRead)},
		{authName: "Root Service Account", auth: NewAuthContext("svc-root", "system", []string{"system"}, "*")},
	}

	for _, v := range verbs {
		for _, p := range payloads {
			for _, a := range authContexts {
				testName := fmt.Sprintf("%s/%s/%s", v.name, p.payloadName, a.authName)
				t.Run(testName, func(t *testing.T) {
					err := v.call(ctx, a.auth, p.payload)
					if err == nil {
						t.Fatalf("CRITICAL SECURITY VIOLATION: Execution verb %s succeeded with %s and %s!", v.name, p.payloadName, a.authName)
					}

					env, ok := err.(*ErrorEnvelope)
					if !ok {
						t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
					}

					if env.Code != ErrorCodeExecutionBlocked {
						t.Fatalf("expected error code %s, got %s", ErrorCodeExecutionBlocked, env.Code)
					}
					if env.HTTPStatus != 400 {
						t.Fatalf("expected HTTP 400, got %d", env.HTTPStatus)
					}
					if env.Retryable {
						t.Fatalf("expected non-retryable error, got retryable=true")
					}
					if !strings.HasPrefix(env.CorrelationID, "corr-") {
						t.Fatalf("expected correlation ID starting with 'corr-', got %q", env.CorrelationID)
					}
					if !strings.Contains(env.Message, "REQ-PROD-006") {
						t.Fatalf("expected message to cite REQ-PROD-006, got %q", env.Message)
					}
					if !strings.Contains(strings.ToLower(env.Message), strings.ToLower(v.name)) {
						t.Fatalf("expected message to cite verb %q, got %q", v.name, env.Message)
					}
				})
			}
		}
	}
}

// =====================================================================
// 5. Tenant Scoping, Wildcards, and Anti-Oracle Isolation
// =====================================================================

func TestAdversarial_TenantScoping_PathTraversalAndWildcards(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	adversarialTenants := []struct {
		tenant      string
		expectedErr ErrorCode
	}{
		{tenant: "", expectedErr: ErrorCodeMissingTenant},
		{tenant: "   ", expectedErr: ErrorCodeMissingTenant},
		{tenant: "*", expectedErr: ErrorCodeWildcardTenant},
		{tenant: "%", expectedErr: ErrorCodeWildcardTenant},
		{tenant: "all", expectedErr: ErrorCodeWildcardTenant},
		{tenant: "ALL", expectedErr: ErrorCodeWildcardTenant},
		{tenant: "All", expectedErr: ErrorCodeWildcardTenant},
		{tenant: " * ", expectedErr: ErrorCodeWildcardTenant},
		// Tenant not in caller's ATS
		{tenant: "tenant-unauthorized", expectedErr: ErrorCodeUnauthorizedTenantInSet},
		{tenant: "tenant-bravo", expectedErr: ErrorCodeUnauthorizedTenantInSet},
	}

	for _, tc := range adversarialTenants {
		t.Run("BrowseTenant_"+tc.tenant, func(t *testing.T) {
			_, err := svc.Browse(ctx, auth, BrowseRequest{
				TenantRef: tc.tenant,
			})
			if err == nil {
				t.Fatalf("expected error for tenant %q, got success", tc.tenant)
			}
			env := err.(*ErrorEnvelope)
			if env.Code != tc.expectedErr {
				t.Fatalf("expected error code %s for tenant %q, got %s", tc.expectedErr, tc.tenant, env.Code)
			}
		})
	}
}

func TestAdversarial_ExistenceOracle_ResponseIndistinguishability(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	// Caller is only authorized for tenant-alpha
	auth := adminAuth("tenant-alpha")

	// Case 1: Query an asset ID that exists in tenant-bravo (foreign tenant)
	_, errForeign := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-bravo-secret", // Exists in tenant-bravo!
	})
	if errForeign == nil {
		t.Fatalf("expected error querying foreign asset")
	}
	envForeign := errForeign.(*ErrorEnvelope)

	// Case 2: Query an asset ID that does NOT exist anywhere in the universe
	_, errNonExistent := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "ghost-asset-404-nonexistent",
	})
	if errNonExistent == nil {
		t.Fatalf("expected error querying non-existent asset")
	}
	envNonExistent := errNonExistent.(*ErrorEnvelope)

	// Both must return identical ErrorCode, HTTP status, and retryability
	if envForeign.Code != ErrorCodeAssetNotFound || envNonExistent.Code != ErrorCodeAssetNotFound {
		t.Fatalf("oracle leak in error code: foreign=%s, nonexistent=%s", envForeign.Code, envNonExistent.Code)
	}
	if envForeign.HTTPStatus != 404 || envNonExistent.HTTPStatus != 404 {
		t.Fatalf("oracle leak in HTTP status: foreign=%d, nonexistent=%d", envForeign.HTTPStatus, envNonExistent.HTTPStatus)
	}
	if envForeign.Retryable != envNonExistent.Retryable {
		t.Fatalf("oracle leak in retryable: foreign=%v, nonexistent=%v", envForeign.Retryable, envNonExistent.Retryable)
	}

	// Message structure must be identically structured ("Asset %q was not found in tenant %q")
	expectedForeignMsg := `Asset "skl-bravo-secret" was not found in tenant "tenant-alpha"`
	if envForeign.Message != expectedForeignMsg {
		t.Fatalf("expected message %q, got %q", expectedForeignMsg, envForeign.Message)
	}
}
