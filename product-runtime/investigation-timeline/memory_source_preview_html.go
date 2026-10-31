package investigationtimeline

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

// sourcePreviewResponse is a read-only projection issued only after the
// original Case -> Entry -> typed source -> source.read authorization chain.
// It does not contain raw owner records or any browser-provided authority.
type sourcePreviewResponse struct {
	Status    string `json:"status"`
	CaseID    string `json:"case_id"`
	EntryID   string `json:"entry_id"`
	Kind      string `json:"kind"`
	ObjectID  string `json:"object_id"`
	Preview   string `json:"preview"`
	Synthetic bool   `json:"synthetic"`
	// The server validates a same-Case return URL; not part of the JSON API.
	ReturnURL string `json:"-"`
}

// A native browser sends text/html as the first Accept preference. A JSON
// client retains the original response, and q=0 never selects HTML.
func acceptsSourcePreviewHTML(accept string) bool {
	first, _, _ := strings.Cut(accept, ",")
	media, params, err := mime.ParseMediaType(strings.TrimSpace(first))
	if err != nil || media != "text/html" {
		return false
	}
	if q, exists := params["q"]; exists {
		quality, err := strconv.ParseFloat(q, 64)
		return err == nil && quality > 0 && quality <= 1
	}
	return true
}

// The destination is already server-authorized; html/template escapes
// owner-provided text, IDs and attributes. This page never embeds scripts,
// client-provided sourceRef, bearer tokens or a new authorization decision.
var sourcePreviewHTMLTemplate = template.Must(template.New("source-preview").Parse(`<!doctype html>
<html lang="fr">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Investigate — Source autorisée</title></head>
<body><main>
<h1>Source de la chronologie</h1>
<p>Prévisualisation synthétique locale — pas une preuve forensique ni une donnée de production.</p>
<dl>
<dt>Case</dt><dd>{{.CaseID}}</dd>
<dt>Entrée</dt><dd>{{.EntryID}}</dd>
<dt>Propriétaire de l'objet (type)</dt><dd>{{.Kind}}</dd>
<dt>Identifiant de l'objet autorisé</dt><dd>{{.ObjectID}}</dd>
</dl>
<h2>Prévisualisation autorisée</h2>
<pre>{{.Preview}}</pre>
<p>Pour retrouver la même page, les filtres et la position, utilisez le bouton Retour du navigateur.
Ce lien restaure la page et les filtres autorisés de départ, en réévaluant les permissions. Sans contexte de retour, la vue du Case est réinitialisée.</p>
<p><a href="{{.ReturnURL}}">Retour à la chronologie du Case</a></p>
</main></body></html>`))

func renderSourcePreviewHTML(out io.Writer, result sourcePreviewResponse) error {
	return sourcePreviewHTMLTemplate.Execute(out, result)
}

// writeAuthorizedSourcePreview handles only presentation of a caller-verified
// result. Authentication, Case/Entry permission, typed source identity and
// owner.read were already checked in sourcePreviewHandler, not in this helper.
// Keeping serialization separate allows deterministic canceled-request checks.
func writeAuthorizedSourcePreview(w http.ResponseWriter, r *http.Request,
	ctx context.Context, payload sourcePreviewResponse) {
	if ctx.Err() != nil || r.Context().Err() != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	// Content negotiation is deliberately AFTER the entire existing owner
	// chain; switching Accept never grants an additional read right.
	if acceptsSourcePreviewHTML(r.Header.Get("Accept")) {
		var page bytes.Buffer
		if err := renderSourcePreviewHTML(&page, payload); err != nil ||
			ctx.Err() != nil || r.Context().Err() != nil {
			writeCaseTimelineError(w, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
		w.Header().Set("Vary", "Authorization, Accept")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(page.Bytes())
		return
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(encoded, '\n'))
}
