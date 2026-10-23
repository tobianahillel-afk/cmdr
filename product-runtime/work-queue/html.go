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
<label for="queue-owner">Responsable exact (facultatif)</label>
<input id="queue-owner" name="owner" type="text" maxlength="128" value="{{.OwnerFilter}}">
<label for="queue-group">Regrouper les éléments autorisés</label>
<select id="queue-group" name="group"><option value=""{{if eq .GroupBy ""}} selected{{end}}>Sans regroupement</option><option value="type"{{if eq .GroupBy "type"}} selected{{end}}>Par type</option><option value="status"{{if eq .GroupBy "status"}} selected{{end}}>Par état</option></select>
<label for="queue-sort">Ordre chronologique</label>
<select id="queue-sort" name="sort"><option value="newest"{{if ne .SortOrder "oldest"}} selected{{end}}>Plus récents</option><option value="oldest"{{if eq .SortOrder "oldest"}} selected{{end}}>Plus anciens</option></select>
<button type="submit">Appliquer / effacer la sélection</button>
</form>
{{if .StatusFilter}}<p>Filtre d’état actif : {{.StatusFilter}}</p>{{end}}
{{if .OwnerFilter}}<p>Filtre de responsable actif : {{.OwnerFilter}}</p>{{end}}
{{if .GroupBy}}<p>Regroupement actif : {{.GroupBy}}. Groupes ordonnés par clé et chronologie conservée dans chaque groupe. Chaque page répète son premier titre de groupe.</p>{{end}}
<p role="status">État : {{.State}}</p>
{{if .Diagnostics}}<ul aria-label="Limites de la projection">{{range .Diagnostics}}<li>{{.}}</li>{{end}}</ul>{{end}}
{{if .ViewFreshness}}<p aria-label="Fraîcheur de la vue autorisée">Fraîcheur de la vue autorisée ({{.ViewFreshness.Total}} éléments) : {{.ViewFreshness.Current}} actuels, {{.ViewFreshness.Aging}} vieillissants, {{.ViewFreshness.Stale}} périmés, {{.ViewFreshness.Partial}} partiels, {{.ViewFreshness.Unavailable}} sources indisponibles, {{.ViewFreshness.Unknown}} inconnus. Décompte du filtre et de la vue autorisés, toutes pages confondues ; aucune fraîcheur déduite sans politique source.</p>{{end}}
{{if .Rows}}
<table><caption>Éléments autorisés — page courante (50 maximum)</caption>
<thead><tr><th scope="col">Type</th><th scope="col">Identifiant</th><th scope="col">Responsable</th><th scope="col">État</th><th scope="col">Mis à jour (UTC)</th><th scope="col">Fraîcheur source</th><th scope="col">Contexte</th></tr></thead>
<tbody>{{range .LinkedRows}}{{if .GroupStart}}<tr><th colspan="7" scope="colgroup">{{.GroupLabel}}</th></tr>{{end}}<tr><td>{{.Kind}}</td><th scope="row">{{.ID}}</th><td>{{.Owner}}</td><td>{{.Status}}</td><td>{{.UpdatedAt}}</td><td>{{if .Freshness}}{{.Freshness.State}} — {{.Freshness.Source}}{{if .Freshness.ObservedAt}} · {{.Freshness.ObservedAt}}{{end}}{{if .Freshness.PolicyVersion}} · policy {{.Freshness.PolicyVersion}}{{end}}{{else}}unknown — source non fournie{{end}}</td><td><form method="get" action="/command/work-queue">
<input type="hidden" name="view" value="{{$.Current}}">
{{if $.StatusFilter}}<input type="hidden" name="status" value="{{$.StatusFilter}}">{{end}}
{{if $.OwnerFilter}}<input type="hidden" name="owner" value="{{$.OwnerFilter}}">{{end}}
{{if $.GroupBy}}<input type="hidden" name="group" value="{{$.GroupBy}}">{{end}}
{{if $.SortOrder}}<input type="hidden" name="sort" value="{{$.SortOrder}}">{{end}}
{{if $.CurrentCursor}}<input type="hidden" name="cursor" value="{{$.CurrentCursor}}">{{end}}
<button type="submit" name="selected" value="{{.SelectionToken}}">Sélectionner</button>
</form></td></tr>{{end}}</tbody>
</table>
{{else}}<p>Aucun élément autorisé dans cette vue.</p>{{end}}
<nav aria-label="Pagination de la Work Queue">
{{if .HasPrevious}}<a href="{{.PreviousPageHref}}">Page précédente</a> <a href="{{.FirstPageHref}}">Première page</a>{{end}}
{{if .NextCursor}}<a href="{{.NextPageHref}}">Page suivante</a>{{end}}
</nav>
{{if .Selected}}<section aria-label="Contexte sélectionné"><h3>Élément sélectionné</h3><p>{{.Selected.Kind}} — {{.Selected.ID}}</p><p>Responsable : {{.Selected.Owner}} · État : {{.Selected.Status}} · Mis à jour : {{.Selected.UpdatedAt}}</p><p>Fraîcheur de la source : {{if .Selected.Freshness}}{{.Selected.Freshness.State}} — {{.Selected.Freshness.Source}}{{if .Selected.Freshness.ObservedAt}} · {{.Selected.Freshness.ObservedAt}}{{end}}{{if .Selected.Freshness.PolicyVersion}} · policy {{.Selected.Freshness.PolicyVersion}}{{end}}{{else}}unknown — source non fournie{{end}}</p><h4>Dimensions de coordination (lecture seule)</h4><dl>
{{range .FactorLines}}<dt>{{.Name}}</dt><dd>{{.Value}}{{if .Source}} — {{.Source}}{{end}}{{if .ObservedAt}} · {{.ObservedAt}}{{end}}</dd>{{end}}
</dl><p>Valeur inconnue : source absente ou non autorisée ; aucune priorité n’est calculée automatiquement. L’Inspector complet reste une capacité distincte.</p><h4>Propriété canonique du produit (lecture seule)</h4>
<p aria-label="Propriété canonique autorisée">{{if .SelectedCanonicalOwnership}}{{if .SelectedCanonicalOwnership.ProductOwner}}Produit propriétaire : {{.SelectedCanonicalOwnership.ProductOwner}}{{else}}Produit propriétaire : inconnu{{end}}{{if .SelectedCanonicalOwnership.RegisterSource}} · registre : {{.SelectedCanonicalOwnership.RegisterSource}}{{else}} · registre : non fourni{{end}}{{if .SelectedCanonicalOwnership.RegisterVersion}} · version : {{.SelectedCanonicalOwnership.RegisterVersion}}{{end}}{{if .SelectedCanonicalOwnership.IdentityStatus}} · identité : {{.SelectedCanonicalOwnership.IdentityStatus}}{{else}} · identité : inconnue{{end}}{{else}}Produit propriétaire : inconnu — registre non fourni{{end}}. Le responsable opérationnel ci-dessus est distinct du propriétaire du produit. Aucun transfert ou changement d'affectation n'est effectué.</p>
{{if eq .Selected.Kind "Incident"}}<p><a href="/command/incident-detail?incident={{urlquery .Selected.ID}}" rel="noreferrer">Ouvrir Incident Detail (lecture seule, autorisation revérifiée)</a></p><h4>Cases Investigate liés (lecture seule)</h4>
{{if .SelectedLinkedCases}}{{if .SelectedLinkedCases.Cases}}
<ul aria-label="Cases Investigate autorisés">{{range .SelectedLinkedCases.Cases}}<li>Case {{.CaseID}} — état {{.Status}} · source de liaison {{.LinkSource}}{{if .LinkVersion}} · version {{.LinkVersion}}{{end}} · fraîcheur du statut {{if .StatusFreshness}}{{.StatusFreshness}}{{else}}unknown{{end}}{{if .StatusSource}} · source du statut {{.StatusSource}}{{else}} · source du statut non fournie{{end}}{{if .StatusObservedAt}} · observé {{.StatusObservedAt}}{{else}} · date inconnue{{end}}{{if .StatusVersion}} · révision {{.StatusVersion}}{{end}} · <a href="/case-timeline?case_id={{urlquery .CaseID}}&amp;origin_incident={{urlquery $.Selected.ID}}" rel="noreferrer">Consulter la chronologie Investigate (autorisation revérifiée à la destination)</a></li>{{end}}</ul>
{{else}}<p>Aucun Case lié visible avec autorisation Investigate actuelle.</p>{{end}}{{else}}<p>Relations inconnues — source Investigate / Object Linking indisponible. Aucun droit Case déduit de l’Incident.</p>{{end}}
{{end}}
<h4>SLA opérationnel (lecture seule)</h4>
<p aria-label="Suivi SLA autorisé">{{if .SelectedSLA}}État : {{.SelectedSLA.State}}{{if .SelectedSLA.PolicySource}} · source : {{.SelectedSLA.PolicySource}}{{end}}{{if .SelectedSLA.PolicyVersion}} · version : {{.SelectedSLA.PolicyVersion}}{{end}}{{if .SelectedSLA.DueAt}} · échéance : {{.SelectedSLA.DueAt}}{{end}}{{if .SelectedSLA.RemainingSeconds}} · restant (secondes) : {{.SelectedSLA.RemainingSeconds}}{{end}}{{if .SelectedSLA.ObservedAt}} · observé : {{.SelectedSLA.ObservedAt}}{{end}}{{if .SelectedSLA.CalendarSource}} · calendrier source : {{.SelectedSLA.CalendarSource}}{{else}} · calendrier : non fourni{{end}}{{if eq .SelectedSLA.State "paused"}} · suspension : {{if .SelectedSLA.PausedAt}}{{.SelectedSLA.PausedAt}}{{else}}horodatage inconnu{{end}} · raison : {{if .SelectedSLA.PauseReason}}{{.SelectedSLA.PauseReason}}{{else}}non fournie{{end}}{{end}}{{else}}unknown — aucun SLA autorisé fourni{{end}}. Aucun seuil ni droit de pause n’est déduit localement.</p></section>{{end}}
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
	GroupStart bool
	GroupLabel string
}

type factorLine struct {
	Name string
	Value string
	Source string
	ObservedAt string
}

func selectedFactorLines(f *AuthorizedFactors) []factorLine {
	var d AuthorizedFactors
	if f != nil { d = *f }
	names := [...]string{"Priorité", "Sévérité", "Urgence", "Impact", "Confiance", "Risque SLA"}
	values := [...] *Dimension{d.Priority, d.Severity, d.Urgency, d.Impact, d.Confidence, d.SLARisk}
	lines := make([]factorLine,0,len(names))
	for i,dim := range values {
		line := factorLine{Name:names[i],Value:"inconnue"}
		if dim != nil {
			line.Value, line.Source, line.ObservedAt = dim.Value,dim.Source,dim.ObservedAt
		}
		lines=append(lines,line)
	}
	return lines
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
		if result.OwnerFilter != "" {
			href += "&owner="+url.QueryEscape(result.OwnerFilter)
		}
		if result.GroupBy != "" {
			href += "&group="+url.QueryEscape(result.GroupBy)
		}
		views = append(views, queuePageView{Key:key, Label:labels[i], Selected:key==result.View, Href:href})
	}
	linked := make([]queueSelectionRow, 0, len(result.Rows))
	previousGroup := ""
	for i, row := range result.Rows {
		group := ""
		switch result.GroupBy {
		case "type": group = string(row.Kind)
		case "status": group = row.Status
		}
		linked = append(linked, queueSelectionRow{
			Row: row, SelectionToken: string(row.Kind)+":"+row.ID,
			GroupStart: result.GroupBy != "" && (i == 0 || group != previousGroup),
			GroupLabel: group,
		})
		previousGroup = group
	}
	firstPageHref := "/command/work-queue?view="+url.QueryEscape(string(result.View))
	if result.StatusFilter != "" {
		firstPageHref += "&status="+url.QueryEscape(result.StatusFilter)
	}
	if result.SortOrder != "" {
		firstPageHref += "&sort="+url.QueryEscape(result.SortOrder)
	}
	if result.OwnerFilter != "" {
		firstPageHref += "&owner="+url.QueryEscape(result.OwnerFilter)
	}
	if result.GroupBy != "" {
		firstPageHref += "&group="+url.QueryEscape(result.GroupBy)
	}
	nextPageHref := ""
	if result.NextCursor != "" {
		nextPageHref = firstPageHref+"&cursor="+url.QueryEscape(result.NextCursor)
	}
	previousPageHref := firstPageHref
	if result.PreviousCursor != "" {
		previousPageHref = firstPageHref+"&cursor="+url.QueryEscape(result.PreviousCursor)
	}
	var out bytes.Buffer
	err := queuePage.Execute(&out, struct {
		Result
		Current View
		Views []queuePageView
		LinkedRows []queueSelectionRow
		FirstPageHref string
		PreviousPageHref string
		NextPageHref string
		FactorLines []factorLine
	}{Result:result, Current:result.View, Views:views, LinkedRows:linked,
		FirstPageHref:firstPageHref, PreviousPageHref:previousPageHref, NextPageHref:nextPageHref,
		FactorLines:selectedFactorLines(result.SelectedFactors)})
	if err != nil { return nil, err }
	return out.Bytes(), nil
}
