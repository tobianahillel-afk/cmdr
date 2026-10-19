"use strict";
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const test = require("node:test");

function harness(replies) {
 const nodes = new Map();
 const elements = [];
 function makeNode(id = "") {
  const events = new Map();
  const node = {id, textContent:"", value:"", disabled:false, dataset:{},
    children:[], events, focusCount:0,
    appendChild(child){this.children.push(child)},
    replaceChildren(...items){this.children=items},
    addEventListener(name,fn){events.set(name,fn)},
    focus(){this.focusCount++},
  };
  Object.defineProperty(node,"innerHTML",{set(){throw Error("unsafe HTML sink")}});
  elements.push(node);
  return node;
 }
 for(const id of ["cases","status","scope","prev","next","page","search","query","clear",
  "case-status","case-owner","case-freshness","case-sort",
  "selection","selection-status","selection-close","selection-title","selection-id",
  "selection-status-value","selection-owner","selection-updated","selection-freshness",
  "selection-next-action","workspace-handoff","workspace-status","workspace-return",
  "workspace-case"])
  nodes.set(id,makeNode(id));
 nodes.get("selection").hidden=true;
 nodes.get("workspace-handoff").hidden=true;
 const requests=[];
 const fetch=async(uri,options)=>{
  requests.push({uri,options});
  const response=replies.shift();
  if(!response) throw Error("test requests exhausted");
  return {ok:response.ok,json:async()=>response.data};
 };
 const window = {events:new Map(),addEventListener(k,fn){this.events.set(k,fn)}};
 const document = {getElementById(id){return nodes.get(id)},createElement(){return makeNode()}};
 vm.runInNewContext(fs.readFileSync("browser.js","utf8"),
  {fetch,document,window,URLSearchParams,AbortController});
 return {nodes,requests,window};
}
const tick=()=>new Promise(r=>setImmediate(r));
const row=(id,title)=>({case:{id,human_id:id,title,status:"investigating",
 owner:"analyst-a",freshness:"available",updated_at:"2026-10-19T12:00:00Z"}});
test("only current authorized rows render as inert text and pagination never sends tenant",async()=>{
 const dangerous="<img src=x onerror=alert(1)>";
 const h=harness([{ok:true,data:{rows:[row("case-a",dangerous)],tenant_ref:"tenant-a",
  environment_ref:"prod",next_cursor:"next-a"}},
  {ok:true,data:{rows:[row("case-b","Case B")],tenant_ref:"tenant-a",environment_ref:"prod"}}]);
 await tick();
 assert.equal(h.nodes.get("cases").children.length,1);
 assert.equal(h.nodes.get("status").textContent,"1 currently accessible Cases on this page.");
 assert.equal(h.nodes.get("cases").children[0].children.length,5);
 assert.equal(h.nodes.get("cases").children[0].children[0].textContent,dangerous+" — case-a");
 assert.equal(h.requests[0].uri,"/case-queue?limit=25");
 assert.equal(h.requests[0].options.credentials,"same-origin");
 assert.equal(h.requests[0].options.cache,"no-store");
 assert.equal(h.nodes.get("scope").textContent,"tenant-a / prod");
 h.nodes.get("next").events.get("click")();
 await tick();
 assert.equal(h.nodes.get("cases").children[0].children[0].textContent,"Case B — case-b");
 assert.match(h.requests[1].uri,/cursor=next-a/);
 assert.doesNotMatch(h.requests[1].uri,/tenant|permission|access_decision|bearer/);
 assert.equal(h.nodes.get("page").textContent,"Page 2");
});
test("denied page clears prior protected rows and never exposes stale data",async()=>{
 const h=harness([{ok:true,data:{rows:[row("private-case","secret")],next_cursor:"opaque"}},
  {ok:false,data:null}]);
 await tick();
 assert.equal(h.nodes.get("cases").children.length,1);
 h.nodes.get("next").events.get("click")();
 await tick();
 assert.equal(h.nodes.get("cases").children.length,0);
 assert.equal(h.nodes.get("scope").textContent,"");
 assert.match(h.nodes.get("status").textContent,/unavailable/);
 assert.equal(h.nodes.get("next").disabled,true);
});
test("query stays client filter only and pagehide erases Case list",async()=>{
 const h=harness([{ok:true,data:{rows:[row("private","title")]}},
  {ok:true,data:{rows:[]}}]);
 await tick();
 h.nodes.get("query").value="hostile<>&";
 h.nodes.get("search").events.get("submit")({preventDefault(){}});
 await tick();
 assert.equal(h.requests[1].uri,"/case-queue?limit=25&q=hostile%3C%3E%26");
 h.window.events.get("pagehide")();
 assert.equal(h.nodes.get("cases").children.length,0);
});

test("Inspect performs an independent server-authorized GET with exact opaque Case ID and restores focus",async()=>{
 const caseID=" case-a ";
 const h=harness([{ok:true,data:{rows:[row(caseID,"Current Case")],tenant_ref:"tenant-a"}},
  {ok:true,data:{case:{id:caseID,title:"<b>secure title</b>",owner:"analyst-a",
   status:"review",freshness:"available",updated_at:"2026-10-19T12:05:00Z",
   next_action:"Investigate"},partial:false,stale:false}}]);
 await tick();
 const button=h.nodes.get("cases").children[0].children[0].children[0];
 assert.equal(button.textContent,"Inspect");
 button.events.get("click")();
 await tick();
 assert.match(h.requests[1].uri,/^\/case-queue\?limit=25&inspect_case=/);
 assert.match(h.requests[1].uri,/inspect_case=\+case-a\+/);
 assert.equal(h.requests[1].options.credentials,"same-origin");
 assert.equal(h.requests[1].options.cache,"no-store");
 assert.equal(h.nodes.get("selection").hidden,false);
 assert.equal(h.nodes.get("selection-id").textContent,caseID);
 assert.equal(h.nodes.get("selection-title").textContent,"<b>secure title</b>");
 assert.equal(h.nodes.get("selection-close").focusCount,1);
 h.window.events.get("keydown")({key:"Escape",preventDefault(){}});
 assert.equal(h.nodes.get("selection").hidden,true);
 assert.equal(button.focusCount,1);
 assert.equal(h.nodes.get("selection-id").textContent,"");
});
test("denied independent Inspector request never reuses previously visible selection",async()=>{
 const h=harness([
  {ok:true,data:{rows:[row("case-a","Visible"),row("case-b","Hidden")]}},
  {ok:true,data:{case:{id:"case-a",title:"Sensitive Case A",owner:"analyst-a"}}},
  {ok:false,data:null}
 ]);
 await tick();
 const rows=h.nodes.get("cases").children;
 rows[0].children[0].children[0].events.get("click")();
 await tick();
 assert.equal(h.nodes.get("selection-title").textContent,"Sensitive Case A");
 rows[1].children[0].children[0].events.get("click")();
 await tick();
 assert.equal(h.nodes.get("selection").hidden,true);
 assert.equal(h.nodes.get("selection-title").textContent,"");
 assert.match(h.nodes.get("selection-status").textContent,/permission changed/);
});
test("Inspector ignores a wrongly-authorized different Case and clears details on navigation",async()=>{
 const h=harness([{ok:true,data:{rows:[row("case-a","Allowed")],next_cursor:"cursor-1"}},
  {ok:true,data:{case:{id:"case-other",title:"Other Case"}}},
  {ok:true,data:{rows:[]}}]);
 await tick();
 h.nodes.get("cases").children[0].children[0].children[0].events.get("click")();
 await tick();
 assert.equal(h.nodes.get("selection").hidden,true);
 assert.equal(h.nodes.get("selection-title").textContent,"");
 h.nodes.get("next").events.get("click")();
 await tick();
 assert.equal(h.nodes.get("selection").hidden,true);
 assert.equal(h.nodes.get("cases").children.length,0);
});

test("applied native status/owner/freshness/order persist across pagination and Inspector reauthorization",async()=>{
 const id=" exact-case-id ";
 const h=harness([
  {ok:true,data:{rows:[row(id,"Restricted")],next_cursor:"page-2"}},
  {ok:true,data:{rows:[row(id,"Currently filtered")],next_cursor:"page-3"}},
  {ok:true,data:{rows:[row("later","Elsewhere")],next_cursor:""}},
  {ok:true,data:{case:{id:"later",title:"Allowed after new authorization",status:"review"}}}
 ]);
 await tick();
 h.nodes.get("query").value="Case Title";
 h.nodes.get("case-status").value="review";
 h.nodes.get("case-owner").value="analyst-a";
 h.nodes.get("case-freshness").value="available";
 h.nodes.get("case-sort").value="updated_at:desc";
 h.nodes.get("search").events.get("submit")({preventDefault(){}});
 await tick();
 assert.equal(h.requests[1].uri,
   "/case-queue?limit=25&q=Case+Title&status=review&owner=analyst-a&freshness=available&sort=updated_at%3Adesc");
 h.nodes.get("next").events.get("click")();
 await tick();
 assert.match(h.requests[2].uri,/cursor=page-3/);
 assert.match(h.requests[2].uri,/status=review/);
 assert.match(h.requests[2].uri,/sort=updated_at%3Adesc/);
 // A user can edit controls but NOT apply them: Inspector must continue
 // checking exact prior filtered page, not an unsubmitted different scope.
 h.nodes.get("case-status").value="secret-unsubmitted";
 h.nodes.get("case-sort").value="title:asc";
 const inspect=h.nodes.get("cases").children[0].children[0].children[0];
 inspect.events.get("click")();
 await tick();
 assert.match(h.requests[3].uri,/inspect_case=later/);
 assert.match(h.requests[3].uri,/status=review/);
 assert.match(h.requests[3].uri,/sort=updated_at%3Adesc/);
 assert.doesNotMatch(h.requests[3].uri,/secret-unsubmitted|sort=title/);
 assert.doesNotMatch(h.requests[3].uri,/tenant|permission|access_decision|bearer/);
});
test("filter reset clears every applied native field and revokes old page and selected detail",async()=>{
 const h=harness([
  {ok:true,data:{rows:[row("a","Sensitive")],next_cursor:"old-page"}},
  {ok:true,data:{rows:[row("b","Filtered")],next_cursor:"new-page"}},
  {ok:true,data:{case:{id:"b",title:"Private selected details"}}},
  {ok:true,data:{rows:[],next_cursor:""}}
 ]);
 await tick();
 h.nodes.get("case-owner").value="analyst-a";
 h.nodes.get("case-status").value="review";
 h.nodes.get("case-sort").value="title:asc";
 h.nodes.get("search").events.get("submit")({preventDefault(){}});
 await tick();
 h.nodes.get("cases").children[0].children[0].children[0].events.get("click")();
 await tick();
 assert.equal(h.nodes.get("selection").hidden,false);
 h.nodes.get("clear").events.get("click")();
 await tick();
 assert.equal(h.nodes.get("selection").hidden,true);
 assert.equal(h.nodes.get("selection-title").textContent,"");
 assert.equal(h.nodes.get("case-status").value,"");
 assert.equal(h.nodes.get("case-owner").value,"");
 assert.equal(h.nodes.get("case-freshness").value,"");
 assert.equal(h.nodes.get("case-sort").value,"");
 assert.equal(h.requests[3].uri,"/case-queue?limit=25");
 assert.equal(h.nodes.get("next").disabled,true);
});

test("browser Workspace handoff reauthorizes exact Case and POST return rechecks the original Queue page",async()=>{
 const h=harness([
  {ok:true,data:{rows:[row("case-a","Sensitive Case A")]}},
  {ok:true,data:{screen_id:"INV-CAS-001",case_id:"case-a",return_token:"opaque-state-only"}},
  {ok:true,data:{origin:"case-queue",selected_case_id:"case-a",
   page_limit:25,page_cursor:"",search_query:"",filters:[],sort:[],scroll_offset:21}},
  {ok:true,data:{rows:[row("case-a","Freshly authorized Case A")]}}
 ]);
 h.window.scrollY=21;
 h.window.scrollTo=(x,y)=>{h.window.lastScroll=[x,y]};
 await tick();
 const rowNode=h.nodes.get("cases").children[0];
 const open=rowNode.children[0].children[1];
 assert.equal(open.textContent,"Open workspace handoff");
 open.events.get("click")();
 await tick();
 assert.match(h.requests[1].uri,/^\/case-queue\?limit=25&open_case=case-a&scroll=21$/);
 assert.equal(h.requests[1].options.credentials,"same-origin");
 assert.equal(h.nodes.get("workspace-handoff").hidden,false);
 assert.equal(h.nodes.get("workspace-case").textContent,"case-a");
 assert.equal(h.nodes.get("cases").children.length,0);
 assert.equal(h.nodes.get("workspace-return").focusCount,1);
 h.nodes.get("workspace-return").events.get("click")();
 await tick();await tick();
 assert.equal(h.requests[2].uri,"/case-queue/return");
 assert.equal(h.requests[2].options.method,"POST");
 assert.equal(h.requests[2].options.headers["Content-Type"],"application/json");
 assert.deepEqual(JSON.parse(h.requests[2].options.body),{return_token:"opaque-state-only"});
 assert.doesNotMatch(h.requests.map(x=>x.uri).join("|"),/opaque-state-only|return_token=|bearer|tenant=/);
 assert.equal(h.requests[3].uri,"/case-queue?limit=25");
 assert.equal(h.nodes.get("workspace-handoff").hidden,true);
 assert.equal(h.nodes.get("workspace-case").textContent,"");
 assert.equal(h.nodes.get("cases").children[0].children[0].textContent,"Freshly authorized Case A — case-a");
 assert.equal(h.window.lastScroll[0],0);
 assert.equal(h.window.lastScroll[1],21);
});
test("Case Workspace handoff denial or wrong Case never retains old sensitive rows",async()=>{
 const h=harness([
  {ok:true,data:{rows:[row("case-a","Sensitive A")]}},
  {ok:true,data:{screen_id:"INV-CAS-001",case_id:"different-case",return_token:"invalid"}},
 ]);
 await tick();
 h.nodes.get("cases").children[0].children[0].children[1].events.get("click")();
 await tick();
 assert.equal(h.nodes.get("cases").children.length,0);
 assert.equal(h.nodes.get("workspace-handoff").hidden,true);
 assert.equal(h.nodes.get("workspace-case").textContent,"");
 assert.match(h.nodes.get("status").textContent,/handoff unavailable/);
});
test("return denial drops transient state and never redisplays stale Case",async()=>{
 const h=harness([
  {ok:true,data:{rows:[row("case-a","Sensitive A")]}},
  {ok:true,data:{screen_id:"INV-CAS-001",case_id:"case-a",return_token:"unsigned-state"}},
  {ok:false,data:null}
 ]);
 await tick();
 h.nodes.get("cases").children[0].children[0].children[1].events.get("click")();
 await tick();
 h.nodes.get("workspace-return").events.get("click")();
 await tick();await tick();
 assert.equal(h.nodes.get("cases").children.length,0);
 assert.equal(h.nodes.get("workspace-handoff").hidden,true);
 assert.match(h.nodes.get("status").textContent,/return denied/);
 assert.equal(h.requests.length,3);
});
test("Workspace handoff preserves applied filtered page and rejects a changed returned projection",async()=>{
 const h=harness([
  {ok:true,data:{rows:[row("case-a","Initial")]}},
  {ok:true,data:{rows:[row("case-b","Filtered")],next_cursor:"next-filtered"}},
  {ok:true,data:{screen_id:"INV-CAS-001",case_id:"case-b",return_token:"opaque-2"}},
  {ok:true,data:{origin:"case-queue",selected_case_id:"case-b",
   page_limit:25,page_cursor:"",search_query:"search-value",
   filters:[{field:"status",values:["wrong-status"]}],sort:[],scroll_offset:0}}
 ]);
 await tick();
 h.nodes.get("query").value="search-value";
 h.nodes.get("case-status").value="review";
 h.nodes.get("search").events.get("submit")({preventDefault(){}});
 await tick();
 h.nodes.get("cases").children[0].children[0].children[1].events.get("click")();
 await tick();
 assert.match(h.requests[2].uri,/q=search-value&status=review&open_case=case-b/);
 h.nodes.get("workspace-return").events.get("click")();
 await tick();await tick();
 assert.equal(h.nodes.get("cases").children.length,0);
 assert.equal(h.nodes.get("workspace-handoff").hidden,true);
 assert.match(h.nodes.get("status").textContent,/return denied/);
 assert.equal(h.requests.length,4);
});

test("Go-normalized case filters still restore a live browser page without widening access",async()=>{
 const h=harness([
  {ok:true,data:{rows:[row("case-a","Initial")]}},
  {ok:true,data:{rows:[row("case-a","Filtered")]}},
  {ok:true,data:{screen_id:"INV-CAS-001",case_id:"case-a",return_token:"nonce-nonauthorizing"}},
  {ok:true,data:{origin:"case-queue",selected_case_id:"case-a",
   page_limit:25,page_cursor:"",search_query:"",
   filters:[{field:"status",values:["review"]},{field:"owner",values:["analyst-a"]}],
   sort:[],scroll_offset:0}},
  {ok:true,data:{rows:[row("case-a","Current owner authorized")]}}
 ]);
 await tick();
 h.nodes.get("case-status").value="Review";
 h.nodes.get("case-owner").value="Analyst-A";
 h.nodes.get("search").events.get("submit")({preventDefault(){}});
 await tick();
 h.nodes.get("cases").children[0].children[0].children[1].events.get("click")();
 await tick();
 h.nodes.get("workspace-return").events.get("click")();
 await tick();await tick();
 assert.equal(h.nodes.get("workspace-handoff").hidden,true);
 assert.equal(h.nodes.get("cases").children[0].children[0].textContent,"Current owner authorized — case-a");
 assert.match(h.requests[4].uri,/status=Review/);
 assert.match(h.requests[4].uri,/owner=Analyst-A/);
});
