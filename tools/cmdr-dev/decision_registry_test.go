package main

import "testing"

func decisionTestGraph() WorkGraph {
	return WorkGraph{Nodes: []WorkNode{{ID: "E5-RES-001A", Status: "READY"}}}
}

func canonicalDecisionPolicy() DecisionPolicy {
	return DecisionPolicy{
		SchemaVersion: 1,
		DefaultPolicy: "deny-unclassified-significant-decision",
		Statuses: []string{"proposed","researching","accepted","rejected","blocked-product","revisit-required","superseded"},
		ProductBoundaries: []string{"engineering-only","product-open-decision","product-behavior","product-security-invariant"},
		RiskTags: []string{
			"local-refactor","representation","serialization","cache","dependency","storage","indexing",
			"performance","concurrency","architecture","security","algorithm","durability","reliability",
		},
		Classes: []DecisionClassPolicy{
			{ID:"A",RecordRequired:false,CurrentExternalResearch:"not-required",AlternativesRequired:false,BenchmarkRule:"not-required",AdversarialReviewRequired:false},
			{ID:"B",RecordRequired:true,CurrentExternalResearch:"if-material-uncertainty",AlternativesRequired:true,BenchmarkRule:"if-performance-sensitive",AdversarialReviewRequired:false},
			{ID:"C",RecordRequired:true,CurrentExternalResearch:"required",AlternativesRequired:true,BenchmarkRule:"if-performance-sensitive",AdversarialReviewRequired:true},
		},
	}
}

func baseDecision() EngineeringDecision {
	return EngineeringDecision{
		ID:"ENG-DEC-0001",DecisionKey:"test.choice",Title:"Test choice",Class:"A",Status:"accepted",
		Owner:"engineering",Question:"Which internal representation?",Rationale:"Equivalent local choice",
		ProductBoundary:"engineering-only",AffectedWorkUnits:[]string{"E5-RES-001A"},
		Constraints:[]string{"preserve external behavior"},RiskTags:[]string{"representation"},
	}
}

func TestDecisionPolicyCanonical(t *testing.T) {
	if err := validateDecisionPolicy(canonicalDecisionPolicy()); err != nil {
		t.Fatal(err)
	}
}

func TestDecisionRegistryAllowsEmptyCanonicalRegistry(t *testing.T) {
	registry := DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions"}
	if err := validateDecisionRegistry(canonicalDecisionPolicy(),registry,decisionTestGraph()); err != nil {
		t.Fatal(err)
	}
}

func TestDecisionRegistryAcceptsClassAEngineeringDecision(t *testing.T) {
	registry := DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{baseDecision()}}
	if err := validateDecisionRegistry(canonicalDecisionPolicy(),registry,decisionTestGraph()); err != nil {
		t.Fatal(err)
	}
}

func TestDecisionRegistryBlocksProductBoundaryAcceptance(t *testing.T) {
	d := baseDecision()
	d.ProductBoundary="product-behavior"
	if err := validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}},decisionTestGraph()); err == nil {
		t.Fatal("expected product-boundary acceptance rejection")
	}
	d.Status="blocked-product"
	if err := validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}},decisionTestGraph()); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptedClassBRequiresEvidenceAndBenchmarkWhenPerformanceSensitive(t *testing.T) {
	d:=baseDecision()
	d.Class="B"
	d.RiskTags=[]string{"performance","cache"}
	d.PerformanceSensitive=true
	registry:=DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}}
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),registry,decisionTestGraph()); err==nil {
		t.Fatal("expected missing evidence rejection")
	}
	d.EvidenceRefs=[]string{"research:packet-1"}
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}},decisionTestGraph()); err==nil {
		t.Fatal("expected missing benchmark rejection")
	}
	d.BenchmarkRefs=[]string{"benchmark:bench-1"}
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}},decisionTestGraph()); err!=nil {
		t.Fatal(err)
	}
}

func TestAcceptedClassCRequiresAdversarialEvidence(t *testing.T) {
	d:=baseDecision()
	d.Class="C"
	d.RiskTags=[]string{"security","architecture"}
	d.CriticalFactors=[]string{"security-critical"}
	d.EvidenceRefs=[]string{"research:packet-1"}
	registry:=DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}}
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),registry,decisionTestGraph()); err==nil {
		t.Fatal("expected missing adversarial evidence rejection")
	}
	d.AdversarialRefs=[]string{"adversarial:review-1"}
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}},decisionTestGraph()); err!=nil {
		t.Fatal(err)
	}
}

func TestCriticalFactorsForceClassC(t *testing.T) {
	d:=baseDecision()
	d.Class="B"
	d.Status="researching"
	d.CriticalFactors=[]string{"algorithm-critical"}
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}},decisionTestGraph()); err==nil {
		t.Fatal("expected critical factor class rejection")
	}
}

func TestDuplicateLiveDecisionKeyFails(t *testing.T) {
	a:=baseDecision()
	b:=baseDecision()
	b.ID="ENG-DEC-0002"
	b.Status="researching"
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{a,b}},decisionTestGraph()); err==nil {
		t.Fatal("expected duplicate live decision_key rejection")
	}
}

func TestSupersededDecisionMustPointToSameKey(t *testing.T) {
	old:=baseDecision()
	old.Status="superseded"
	old.SupersededBy="ENG-DEC-0002"
	newer:=baseDecision()
	newer.ID="ENG-DEC-0002"
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{old,newer}},decisionTestGraph()); err!=nil {
		t.Fatal(err)
	}
	newer.DecisionKey="different.key"
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{old,newer}},decisionTestGraph()); err==nil {
		t.Fatal("expected cross-key supersession rejection")
	}
}

func TestUnknownAffectedWorkUnitFails(t *testing.T) {
	d:=baseDecision()
	d.AffectedWorkUnits=[]string{"E9-NOT-REAL"}
	if err:=validateDecisionRegistry(canonicalDecisionPolicy(),DecisionRegistry{SchemaVersion:1,RegistryKind:"engineering-technical-decisions",Decisions:[]EngineeringDecision{d}},decisionTestGraph()); err==nil {
		t.Fatal("expected unknown work unit rejection")
	}
}
