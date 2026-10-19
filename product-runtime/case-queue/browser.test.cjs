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
  "selection","selection-status","selection-close","selection-title","selection-id",
  "selection-status-value","selection-owner","selection-updated","selection-freshness",
  "selection-next-action"])
  nodes.set(id,makeNode(id));
 nodes.get("selection").hidden=true;
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
