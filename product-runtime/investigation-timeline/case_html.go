package investigationtimeline

import (
	"bytes"
	"fmt"
	"html/template"
	"mime"
	"net/url"
	"strconv"
	"strings"
)

const defaultHTMLTimelineLimit = 50

// Only browsers explicitly accepting text/html receive HTML. JSON API clients
// keep the original wire format, including unpaged and paged responses.
func acceptsTimelineHTML(accept string) bool {
	for _, part := range strings.Split(accept, ",") {
		mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(part))
		if err != nil {
			continue
		}
		if mediaType == "text/html" {
			if quality, specified := params["q"]; specified {
				q, err := strconv.ParseFloat(quality, 64)
				if err != nil || q <= 0 || q > 1 {
					return false
				}
			}
			return true
		}
		if mediaType == "application/json" {
			return false
		}
	}
	return false
}

type authorizedTimelineHTML struct {
	CaseID        string
	Status        string
	Entries       []ProjectedEntry
	TotalVisible  int
	PageLimit     int
	PageOffset    int
	HasPrevious   bool
	PreviousHref  string
	HasNext       bool
	NextHref      string
}

// html/template is essential here: every field is owner-supplied evidence
// metadata, never pre-trusted HTML. No browser-side script is needed.
var authorizedTimelineTemplate = template.Must(template.New("case-timeline").Parse(`<!doctype html>
<html lang="fr">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Investigate — Chronologie autorisée du Case</title></head>
<body>
<header><h1>Chronologie Investigate</h1><p>Case autorisé : <strong>{{.CaseID}}</strong></p></header>
<main>
<p role="status">Projection autorisée : {{.Status}}. {{.TotalVisible}} entrée(s) visibles après filtrage, toutes pages confondues. Aucune permission n'est déduite du lien Command.</p>
<p>Une inférence ne constitue pas une observation. Les dates incertaines sont signalées ; l’absence du statut « Remplacée » ne prouve pas l’absence d’une version inaccessible.</p>
{{if .Entries}}
<table>
<caption>Chronologie de Case — entrées autorisées uniquement</caption>
<thead><tr><th scope="col">Entrée</th><th scope="col">Nature</th><th scope="col">Horodatage / fuseau</th><th scope="col">Source / provenance</th><th scope="col">Version</th><th scope="col">Correction</th></tr></thead>
<tbody>{{range .Entries}}
<tr><th scope="row">{{.ID}}</th><td>{{.Kind}}</td>
<td>{{.SourceTime}} · {{.Timezone}}{{if .TimeUncertain}} · Chronologie incertaine{{end}}</td>
<td>{{.SourceRef}} · {{.Provenance}}</td>
<td>{{.Version}}</td>
<td>{{if .Superseded}}Remplacée (version ultérieure autorisée visible){{else}}Statut de remplacement non établi{{end}}{{if .Supersedes}} · Corrige l’entrée autorisée {{.Supersedes}}{{end}}</td></tr>
{{end}}</tbody>
</table>
{{else}}
<p>Aucune entrée autorisée dans cette page ou cette sélection.</p>
{{end}}
<nav aria-label="Pagination de la chronologie">
{{if .HasPrevious}}<a href="{{.PreviousHref}}">Page précédente</a>{{end}}
{{if .HasNext}}<a href="{{.NextHref}}">Page suivante</a>{{end}}
</nav>
</main>
<footer><p>Lecture seule — données et droits réévalués à chaque navigation. Le nombre d’entrées peut changer entre deux pages ; aucune capture temporelle transactionnelle n’est revendiquée.</p>
<a href="/command/work-queue" rel="noreferrer">Ouvrir la Work Queue Command (nouvelle autorisation requise)</a></footer>
</body></html>`))

// renderAuthorizedTimelineHTML has no owner access itself: its sole input is
// the already-filtered Result from the exact same Case HTTP authority path.
// It does not expose hidden successor IDs, and never promotes link navigation
// into any read permission.
func renderAuthorizedTimelineHTML(caseID string, filtered Result, selection timelineReadSelection, original url.Values) ([]byte, error) {
	if filtered.Status == "deny" {
		return nil, fmt.Errorf("denied timeline cannot render")
	}
	if !selection.paged {
		selection.paged = true
		selection.limit = defaultHTMLTimelineLimit
		selection.offset = 0
	}
	page := paginateAuthorizedTimeline(filtered, selection)
	link := func(offset int) string {
		values := make(url.Values, len(original)+2)
		for key, items := range original {
			values[key] = append([]string(nil), items...)
		}
		values.Set("case_id", caseID)
		values.Set("limit", strconv.Itoa(selection.limit))
		values.Set("offset", strconv.Itoa(offset))
		return "/case-timeline?" + values.Encode()
	}
	view := authorizedTimelineHTML{
		CaseID: caseID, Status: page.Status, Entries: page.Entries,
		TotalVisible: page.Page.TotalVisible, PageLimit: selection.limit,
		PageOffset: selection.offset, HasPrevious: selection.offset > 0,
		HasNext: page.Page.HasMore,
	}
	if view.HasPrevious {
		view.PreviousHref = link(max(0, selection.offset-selection.limit))
	}
	if view.HasNext {
		view.NextHref = link(selection.offset+selection.limit)
	}
	var out bytes.Buffer
	if err := authorizedTimelineTemplate.Execute(&out, view); err != nil {
		return nil, fmt.Errorf("render authorized timeline: %w", err)
	}
	return out.Bytes(), nil
}
