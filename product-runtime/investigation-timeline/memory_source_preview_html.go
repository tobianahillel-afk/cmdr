package investigationtimeline

import (
	"html/template"
	"io"
	"mime"
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
Ce lien rouvre la chronologie en réévaluant les autorisations et réinitialise les filtres.</p>
<p><a href="/case-timeline?case_id={{urlquery .CaseID}}">Retour à la chronologie du Case</a></p>
</main></body></html>`))

func renderSourcePreviewHTML(out io.Writer, result sourcePreviewResponse) error {
	return sourcePreviewHTMLTemplate.Execute(out, result)
}
