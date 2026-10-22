// A bounded local performance probe for the actual Command Work Queue runtime.
// It measures only in-process projection with one trusted synthetic owner.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	queue "github.com/tobianahillel-afk/cmdr/product-runtime/work-queue"
)

type syntheticOwner struct{}
func (syntheticOwner) List(_ context.Context, tenant, env string) ([]queue.Item,error) {
	return []queue.Item{{Kind:queue.Incident,ID:"incident-perf",Tenant:tenant,Environment:env,Owner:"analyst",Status:"active",UpdatedAt:"2026-10-19T10:00:00Z"}},nil
}
func (syntheticOwner) CanRead(_ context.Context, principal, tenant, env string, kind queue.Kind, id string)(bool,error){
	return principal=="analyst"&&tenant=="perf-tenant"&&env=="perf-env"&&kind==queue.Incident&&id=="incident-perf",nil
}
type observation struct {
	Operation string `json:"operation"`
	Mode string `json:"mode"`
	Operations int `json:"operations"`
	ElapsedNS int64 `json:"elapsed_ns"`
	NSPerOperation float64 `json:"ns_per_operation"`
}
var sink queue.Result
func run(args []string, out io.Writer) error {
	fs:=flag.NewFlagSet("command-work-queue-perf",flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	n:=fs.Int("iterations",10000,"operations")
	mode:=fs.String("mode","latency","mode")
	op:=fs.String("operation","projection","operation")
	if err:=fs.Parse(args);err!=nil{return err}
	if fs.NArg()!=0||*n<1||*n>100000||*mode!="latency"||*op!="projection"{
		return fmt.Errorf("unsupported performance operation, mode or iteration count")
	}
	req:=queue.Request{Principal:"analyst",Tenant:"perf-tenant",Environment:"perf-env",View:queue.All}
	owner:=syntheticOwner{}
	perform:=func()error{
		result,err:=queue.Project(context.Background(),owner,req)
		if err!=nil{return err}
		if len(result.Rows)!=1||result.Rows[0].ID!="incident-perf" {return fmt.Errorf("unexpected Command projection")}
		sink=result
		return nil
	}
	if err:=perform();err!=nil{return err}
	started:=time.Now()
	for i:=0;i<*n;i++ {if err:=perform();err!=nil{return err}}
	elapsed:=time.Since(started)
	if elapsed<=0{return fmt.Errorf("unmeasured Command projection")}
	return json.NewEncoder(out).Encode(observation{
		Operation:*op,Mode:*mode,Operations:*n,
		ElapsedNS:elapsed.Nanoseconds(),NSPerOperation:float64(elapsed.Nanoseconds())/float64(*n),
	})
}
func main(){if err:=run(os.Args[1:],os.Stdout);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(2)}}
