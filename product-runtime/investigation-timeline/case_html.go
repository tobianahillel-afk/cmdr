package investigationtimeline

import (
	"bytes"
	"html/template"
	"io"
	"mime"
	"net/url"
	"strconv"
	"strings"
)

// The canonical owner-checked Case Timeline read also has a browser
// representation. This template is not an independent permissions engine;
// it receives only already-authorized and bounded ProjectedEntry values.
var timelineHTMLTemplate = template.Must(template.New("case-timeline-html").Funcs(template.FuncMap{
	"sourceURL": localSourceOpenURL,
}).Parse(`<!doctype html>
<html lang="fr">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Investigate — Chronologie du Case</title>
</head>
<body>
<main>
<h1>Chronologie du Case</h1>
<p>Case : {{.CaseID}}</p>
<p role="status">État autorisé : {{.Result.Status}}. Les sources peuvent être partielles ou manquantes. Les périodes à horloge incertaine ne prouvent aucun ordre causal.</p>
<p>Éléments autorisés après filtres : {{.Result.Page.TotalVisible}}. Affichage : {{.Result.Page.Limit}} lignes maximum. La chronologie est réévaluée avec les permissions serveur à chaque requête ; la pagination n'est pas un instantané transactionnel.</p>
<form method="get" action="/case-timeline" aria-label="Filtrer les entrées autorisées">
<input type="hidden" name="case_id" value="{{.CaseID}}">
<label for="timeline-kind">Type</label>
<select id="timeline-kind" name="kind">
<option value="">Tous les types</option>
<option value="observed" {{if eq .Kind "observed"}}selected{{end}}>Observé</option>
<option value="inferred" {{if eq .Kind "inferred"}}selected{{end}}>Inféré</option>
<option value="user-action" {{if eq .Kind "user-action"}}selected{{end}}>Action utilisateur</option>
<option value="decision" {{if eq .Kind "decision"}}selected{{end}}>Décision</option>
<option value="response" {{if eq .Kind "response"}}selected{{end}}>Réponse</option>
</select>
<label for="timeline-search">Rechercher dans les champs autorisés (identifiant, source et provenance)</label>
<input id="timeline-search" name="q" type="search" maxlength="96" value="{{.Query}}">
<label for="timeline-from">Début inclus (RFC3339, ex. 2026-10-21T10:00:00Z)</label>
<input id="timeline-from" name="from" type="text" maxlength="35" value="{{.From}}" placeholder="2026-10-21T10:00:00Z">
<label for="timeline-to">Fin exclue (RFC3339)</label>
<input id="timeline-to" name="to" type="text" maxlength="35" value="{{.To}}" placeholder="2026-10-21T11:00:00Z">
<label for="timeline-limit">Entrées par page (1 à 200)</label>
<input id="timeline-limit" name="limit" type="number" min="1" max="200" required value="{{.Result.Page.Limit}}">
<button type="submit">Appliquer les filtres autorisés</button>
<a href="{{.ResetURL}}">Réinitialiser les filtres</a>
</form>
<p>Les filtres ne changent jamais les permissions : seuls les événements et métadonnées déjà autorisés sont recherchés. Les horaires incertains restent signalés et ne prouvent aucun ordre causal.</p>
{{if .Result.Entries}}
<table>
<caption>Timeline Entry autorisées, lecture seule</caption>
<thead><tr>
<th scope="col">Identifiant</th><th scope="col">Type</th><th scope="col">Heure source</th>
<th scope="col">Fuseau</th><th scope="col">Fiabilité temporelle</th>
<th scope="col">Source</th><th scope="col">Provenance</th><th scope="col">Version</th>
<th scope="col">État de remplacement</th><th scope="col">Version précédente autorisée</th>{{if $.SourceLinks}}<th scope="col">Source autorisée</th>{{end}}
</tr></thead>
<tbody>{{range .Result.Entries}}
<tr>
<th scope="row">{{.ID}}</th><td>{{.Kind}}</td><td>{{.SourceTime}}</td>
<td>{{.Timezone}}</td><td>{{if .TimeUncertain}}Incertitude — ne pas ordonner causalement{{else}}Heure interprétable{{end}}</td>
<td>{{.SourceRef}}</td><td>{{.Provenance}}</td><td>{{.Version}}</td>
<td>{{if .Superseded}}Remplacée (une version ultérieure autorisée est visible){{else}}Non déterminé (aucun successeur autorisé visible dans cette lecture){{end}}</td>
<td>{{if .Supersedes}}{{.Supersedes}}{{else}}—{{end}}</td>
{{if $.SourceLinks}}<td><a href="{{sourceURL $.CaseID .ID $.ReturnURL}}">Ouvrir la source (nouvelle autorisation)</a></td>{{end}}
</tr>{{end}}</tbody>
</table>
{{else}}<p>Aucune entrée autorisée visible avec ces filtres.</p>{{end}}
{{if .SourceLinks}}<p>La source est relue et autorisée séparément à chaque ouverture. Pour revenir à la même page et aux mêmes filtres de Case, utiliser le retour du navigateur.</p>{{end}}
<nav aria-label="Pagination des entrées autorisées">
{{if .HasPrevious}}<a rel="prev" href="{{.PreviousURL}}">Page précédente</a>{{end}}
{{if .HasNext}}<a rel="next" href="{{.NextURL}}">Page suivante</a>{{end}}
</nav>
<p>Ce tableau est une représentation des entrées déjà autorisées. Un statut absent ne prouve pas qu'aucune correction masquée n'existe. La consultation d'un objet source exige un contrôle de permission distinct : les références source ci-dessus ne sont pas des liens donnant accès.</p>
</main>
</body>
</html>`))

type timelineHTMLView struct {
	CaseID      string
	Kind        string
	Query       string
	From        string
	To          string
	ResetURL    string
	ReturnURL   string
	SourceLinks bool
	Result      timelinePagedHTTPResult
	HasPrevious bool
	HasNext     bool
	PreviousURL string
	NextURL     string
}

// HTML is explicit content negotiation, not an authorization signal. A
// standard navigation Accept header starts with text/html; a JSON API
// request continues receiving its unchanged JSON representation.
func wantsTimelineHTML(accept string) bool {
	if len(accept) > 1024 {
		return false
	}
	first := strings.TrimSpace(strings.SplitN(accept, ",", 2)[0])
	media, options, err := mime.ParseMediaType(first)
	if err != nil || media != "text/html" {
		return false
	}
	if q, exists := options["q"]; exists {
		weight, err := strconv.ParseFloat(q, 64)
		if err != nil || weight <= 0 || weight > 1 {
			return false
		}
	}
	return true
}

// Copy only validated query keys. No principal, scope or permission travels
// in an HTML link; the same server authorization runs on every new page.
// Only the fixed loopback source-open route is ever emitted. Both selectors
// are opaque identifiers already authorized in the displayed Case projection.
func localSourceOpenURL(caseID, entryID, returnURL string) string {
	q := url.Values{"case_id": {caseID}, "entry_id": {entryID}}
	if returnURL != "" {
		q.Set("return_to", returnURL)
	}
	return "/demo-source-open?" + q.Encode()
}

func timelineHTMLPageURL(values url.Values, limit, offset int) string {
	q := make(url.Values, len(values)+2)
	for k, entries := range values {
		if k != "limit" && k != "offset" && len(entries) == 1 {
			q.Set(k, entries[0])
		}
	}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	return "/case-timeline?" + q.Encode()
}

func renderTimelineHTML(w io.Writer, values url.Values, caseID string, selection timelineReadSelection, authorized Result, sourceLinks ...bool) error {
	if !selection.paged {
		selection.paged = true
		selection.limit = 50
	}
	page := paginateAuthorizedTimeline(authorized, selection)
	view := timelineHTMLView{
		CaseID: caseID, Result: page,
		SourceLinks: len(sourceLinks) != 0 && sourceLinks[0],
		Kind:        values.Get("kind"), Query: values.Get("q"),
		From: values.Get("from"), To: values.Get("to"),
		ResetURL:    timelineHTMLPageURL(url.Values{"case_id": {caseID}}, page.Page.Limit, 0),
		ReturnURL:   timelineHTMLPageURL(values, page.Page.Limit, selection.offset),
		HasPrevious: selection.offset > 0,
		HasNext:     page.Page.HasMore,
	}
	if view.HasPrevious {
		view.PreviousURL = timelineHTMLPageURL(values, selection.limit, max(0, selection.offset-selection.limit))
	}
	if view.HasNext {
		view.NextURL = timelineHTMLPageURL(values, selection.limit, selection.offset+selection.limit)
	}
	// No partial response if rendering fails: caller writes the buffered body
	// only after the complete template succeeds.
	var encoded bytes.Buffer
	if err := timelineHTMLTemplate.Execute(&encoded, view); err != nil {
		return err
	}
	_, err := io.Copy(w, &encoded)
	return err
}
