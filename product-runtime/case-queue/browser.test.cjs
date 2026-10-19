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
 for(const id of ["cases","status","scope","prev","next","page","search","query","clear"])
  nodes.set(id,makeNode(id));
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
 assert.equal(h.nodes.get("cases").children[0].children[0].textContent,"case-a");
 assert.equal(h.requests[0].uri,"/case-queue?limit=25");
 assert.equal(h.requests[0].options.credentials,"same-origin");
 assert.equal(h.requests[0].options.cache,"no-store");
 assert.equal(h.nodes.get("scope").textContent,"tenant-a / prod");
 h.nodes.get("next").events.get("click")();
 await tick();
 assert.equal(h.nodes.get("cases").children[0].children[0].textContent,"case-b");
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
