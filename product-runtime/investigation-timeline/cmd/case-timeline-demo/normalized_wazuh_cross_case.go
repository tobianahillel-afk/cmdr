package main

import (
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Retrospective time co-occurrence of the SAME rule+agent in independently
// authorized Cases, not an alert verdict or proof of global event identity.
type sharedAuthorizedWazuhWindow struct {
	StartUTC string `json:"start_utc"`
	RuleID string `json:"rule_id"`
	AgentID string `json:"agent_id"`
	FirstEvents int `json:"first_events"`
	SecondEvents int `json:"second_events"`
}

func sharedAuthorizedWazuhWindows(first,second timeline.MemoryTimelineData,firstCase,secondCase string)([]sharedAuthorizedWazuhWindow,error){
	a,err:=collectAuthorizedWazuhObservations(first,firstCase)
	if err!=nil{return nil,err}
	b,err:=collectAuthorizedWazuhObservations(second,secondCase)
	if err!=nil{return nil,err}
	type bin struct { window string; cohort authorizedWazuhCohort }
	count:=func(groups map[authorizedWazuhCohort][]authorizedWazuhObservation)map[bin]int{
		out:=make(map[bin]int)
		for cohort,observations:=range groups {
			for _,event:=range observations {
				key:=bin{window:event.at.Truncate(5*time.Minute).Format(time.RFC3339),cohort:cohort}
				out[key]++
			}
		}
		return out
	}
	firstCount,secondCount:=count(a),count(b)
	var result []sharedAuthorizedWazuhWindow
	for key,n:=range firstCount {
		m:=secondCount[key]
		if m==0{continue}
		result=append(result,sharedAuthorizedWazuhWindow{
			StartUTC:key.window,RuleID:key.cohort.rule,AgentID:key.cohort.agent,
			FirstEvents:n,SecondEvents:m,
		})
	}
	sort.Slice(result,func(i,j int)bool{
		a,b:=result[i],result[j]
		if a.StartUTC!=b.StartUTC{return a.StartUTC<b.StartUTC}
		if a.RuleID!=b.RuleID{return a.RuleID<b.RuleID}
		return a.AgentID<b.AgentID
	})
	if len(result)>10{return result[:10],nil}
	return result,nil
}
