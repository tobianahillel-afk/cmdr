package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type commandWorkQueueProbe struct {
	moduleRoot string
	binaryPath string
}
func prepareCommandWorkQueueProbe(root string) (commandWorkQueueProbe,func(),error){
	moduleRoot,err:=resolveRepoPath(root,"product-runtime/work-queue",false)
	if err!=nil{return commandWorkQueueProbe{},func(){},fmt.Errorf("Command queue unavailable: %w",err)}
	dir,err:=os.MkdirTemp(root,".cmdr-command-work-queue-perf-")
	if err!=nil{return commandWorkQueueProbe{},func(){},err}
	cleanup:=func(){_ = os.RemoveAll(dir)} // #nosec G703 - created below checked repository root
	name:="command-work-queue-perf"
	if runtime.GOOS=="windows"{name+=".exe"}
	binary:=filepath.Join(dir,name)
	if _,err:=runPerformanceProcess(moduleRoot,"go","build","-trimpath","-o",binary,"./cmd/perf-probe");err!=nil{
		cleanup()
		return commandWorkQueueProbe{},func(){},fmt.Errorf("build Command queue probe: %w",err)
	}
	return commandWorkQueueProbe{moduleRoot:moduleRoot,binaryPath:binary},cleanup,nil
}
func (p commandWorkQueueProbe) run(iterations int)(pilotProjectionProbeObservation,error){
	if iterations<1||iterations>100000{return pilotProjectionProbeObservation{},fmt.Errorf("invalid iteration count")}
	raw,err:=runPerformanceProcess(p.moduleRoot,p.binaryPath,"-iterations",strconv.Itoa(iterations),"-mode","latency","-operation","projection")
	if err!=nil{return pilotProjectionProbeObservation{},err}
	var sample pilotProjectionProbeObservation
	if err:=json.Unmarshal([]byte(strings.TrimSpace(raw)),&sample);err!=nil{return sample,err}
	if sample.Operation!="projection"||sample.Mode!="latency"||sample.Operations!=iterations||
		sample.ElapsedNS<=0||sample.NSPerOperation<=0||math.IsNaN(sample.NSPerOperation)||math.IsInf(sample.NSPerOperation,0){
		return sample,fmt.Errorf("invalid Command queue probe output")
	}
	return sample,nil
}
