package workqueue

import (
	"bytes"
	"html/template"
	"net/url"
)

// The browser is a second representation of the SAME server-authorized
// projection, not a second route, Saved View engine or trusted client.
var queuePage = template.Must(template.New("command-work-queue").Parse(`<!doctype html>
<html lang="fr">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Command — Work Queue</title></head>
<body>
<header><h1>Work Queue</h1><p>Tenant : {{.Tenant}} — Environnement : {{.Environment}}</p></header>
<nav aria-label="Vues système">
{{range .Views}}<a href="{{.Href}}"{{if .Selected}} aria-current="page"{{end}}>{{.Label}}</a> {{end}}
</nav>
<main>
<h2>Vue : {{.Current}}</h2>
<form method="get" action="/command/work-queue">
<input type="hidden" name="view" value="{{.Current}}">
<label for="queue-status">État exact (Incident ou Task)</label>
<input id="queue-status" name="status" type="text" maxlength="64" value="{{.StatusFilter}}">
<label for="queue-sort">Ordre chronologique</label>
<select id="queue-sort" name="sort"><option value="newest"{{if ne .SortOrder "oldest"}} selected{{end}}>Plus récents</option><option value="oldest"{{if eq .SortOrder "oldest"}} selected{{end}}>Plus anciens</option></select>
<button type="submit">Appliquer / effacer la sélection</button>
{{if .StatusFilter}}<p>Filtre actif : {{.StatusFilter}}</p>{{end}}
<p role="status">État : {{.State}}</p>
{{if .Diagnostics}}<ul aria-label="Limites de la projection">{{range .Diagnostics}}<li>{{.}}</li>{{end}}</ul>{{end}}
{{if .Rows}}
<table><caption>Éléments autorisés — première page (50 maximum)</caption>
<thead><tr><th scope="col">Type</th><th scope="col">Identifiant</th><th scope="col">Responsable</th><th scope="col">État</th><th scope="col">Mis à jour (UTC)</th><th scope="col">Fraîcheur source</th><th scope="col">Contexte</th></tr></thead>
<tbody>{{range .LinkedRows}}<tr><td>{{.Kind}}</td><th scope="row">{{.ID}}</th><td>{{.Owner}}</td><td>{{.Status}}</td><td>{{.UpdatedAt}}</td><td>{{if .Freshness}}{{.Freshness.State}} — {{.Freshness.Source}}{{if .Freshness.ObservedAt}} · {{.Freshness.ObservedAt}}{{end}}{{if .Freshness.PolicyVersion}} · policy {{.Freshness.PolicyVersion}}{{end}}{{else}}unknown — source non fournie{{end}}</td><td><button type="submit" name="selected" value="{{.SelectionToken}}">Sélectionner</button></td></tr>{{end}}</tbody>
</table>
{{else}}<p>Aucun élément autorisé dans cette vue.</p>{{end}}
</form>
<nav aria-label="Pagination de la Work Queue">
{{if .HasPrevious}}<a href="{{.FirstPageHref}}">Première page</a>{{end}}
{{if .NextCursor}}<a href="{{.NextPageHref}}">Page suivante</a>{{end}}
</nav>
{{if .Selected}}<section aria-label="Contexte sélectionné"><h3>Élément sélectionné</h3><p>{{.Selected.Kind}} — {{.Selected.ID}}</p><p>Responsable : {{.Selected.Owner}} · État : {{.Selected.Status}} · Mis à jour : {{.Selected.UpdatedAt}}</p><p>Fraîcheur de la source : {{if .Selected.Freshness}}{{.Selected.Freshness.State}} — {{.Selected.Freshness.Source}}{{if .Selected.Freshness.ObservedAt}} · {{.Selected.Freshness.ObservedAt}}{{end}}{{if .Selected.Freshness.PolicyVersion}} · policy {{.Selected.Freshness.PolicyVersion}}{{end}}{{else}}unknown — source non fournie{{end}}</p><p>Cette projection est en lecture seule ; le détail et l’Inspector restent des capacités distinctes.</p></section>{{end}}
</main></body></html>`))

type queuePageView struct {
	Key View
	Label string
	Selected bool
	Href string
}

type queueSelectionRow struct {
	Row
	SelectionToken string
}

func renderQueuePage(result Result) ([]byte, error) {
	labels := [...]string{"Tout", "Incidents", "Tâches", "Non assignés", "Risque SLA", "Mon travail"}
	views := make([]queuePageView, 0, len(SystemViews))
	for i, key := range SystemViews {
		href := "/command/work-queue?view="+url.QueryEscape(string(key))
		if result.StatusFilter != "" {
			href += "&status="+url.QueryEscape(result.StatusFilter)
		}
		if result.SortOrder != "" {
			href += "&sort="+url.QueryEscape(result.SortOrder)
		}
		views = append(views, queuePageView{Key:key, Label:labels[i], Selected:key==result.View, Href:href})
	}
	linked := make([]queueSelectionRow, 0, len(result.Rows))
	for _, row := range result.Rows {
		linked = append(linked, queueSelectionRow{
			Row: row, SelectionToken: string(row.Kind)+":"+row.ID,
		})
	}
	firstPageHref := "/command/work-queue?view="+url.QueryEscape(string(result.View))
	if result.StatusFilter != "" {
		firstPageHref += "&status="+url.QueryEscape(result.StatusFilter)
	}
	if result.SortOrder != "" {
		firstPageHref += "&sort="+url.QueryEscape(result.SortOrder)
	}
	nextPageHref := ""
	if result.NextCursor != "" {
		nextPageHref = firstPageHref+"&cursor="+url.QueryEscape(result.NextCursor)
	}
	var out bytes.Buffer
	err := queuePage.Execute(&out, struct {
		Result
		Current View
		Views []queuePageView
		LinkedRows []queueSelectionRow
		FirstPageHref string
		NextPageHref string
	}{Result:result, Current:result.View, Views:views, LinkedRows:linked,
		FirstPageHref:firstPageHref, NextPageHref:nextPageHref})
	if err != nil { return nil, err }
	return out.Bytes(), nil
}
