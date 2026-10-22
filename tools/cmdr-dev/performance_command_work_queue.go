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

func prepareCommandWorkQueueProbe(root string) (commandWorkQueueProbe, func(), error) {
	moduleRoot, err := resolveRepoPath(root, "product-runtime/work-queue", false)
	if err != nil {
		return commandWorkQueueProbe{}, func() {}, fmt.Errorf("Command queue unavailable: %w", err)
	}
	dir, err := os.MkdirTemp(root, ".cmdr-command-work-queue-perf-")
	if err != nil {
		return commandWorkQueueProbe{}, func() {}, err
	}
	cleanup := func() { _ = os.RemoveAll(dir) } // #nosec G703 - created below checked repository root
	name := "command-work-queue-perf"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(dir, name)
	if _, err := runPerformanceProcess(moduleRoot, "go", "build", "-trimpath", "-o", binary, "./cmd/perf-probe"); err != nil {
		cleanup()
		return commandWorkQueueProbe{}, func() {}, fmt.Errorf("build Command queue probe: %w", err)
	}
	return commandWorkQueueProbe{moduleRoot: moduleRoot, binaryPath: binary}, cleanup, nil
}

func (p commandWorkQueueProbe) run(iterations int) (pilotProjectionProbeObservation,error) {
	return p.runMode(iterations,"latency",0)
}
func (p commandWorkQueueProbe) runResource(iterations,memoryMiB int) (pilotProjectionProbeObservation,error) {
	return p.runMode(iterations,"resource",memoryMiB)
}
func (p commandWorkQueueProbe) runMode(iterations int,mode string,memoryMiB int) (pilotProjectionProbeObservation,error) {
	if iterations<1||iterations>100000{return pilotProjectionProbeObservation{},fmt.Errorf("invalid iteration count")}
	if mode!="latency"&&mode!="resource"{return pilotProjectionProbeObservation{},fmt.Errorf("invalid mode")}
	args:=[]string{"-iterations",strconv.Itoa(iterations),"-mode",mode,"-operation","projection"}
	if mode=="resource" {
		if memoryMiB<64||memoryMiB>8192{return pilotProjectionProbeObservation{},fmt.Errorf("invalid memory limit")}
		args=append(args,"-memory-limit-mib",strconv.Itoa(memoryMiB))
	}
	raw,err:=runPerformanceProcess(p.moduleRoot,p.binaryPath,args...)
	if err!=nil{return pilotProjectionProbeObservation{},err}
	var sample pilotProjectionProbeObservation
	if err:=json.Unmarshal([]byte(strings.TrimSpace(raw)),&sample);err!=nil{return sample,err}
	if sample.Operation!="projection"||sample.Mode!=mode||sample.Operations<iterations||sample.Operations>1000000||sample.ElapsedNS<=0||sample.NSPerOperation<=0||math.IsNaN(sample.NSPerOperation)||math.IsInf(sample.NSPerOperation,0)||(mode=="resource"&&sample.PeakHeapBytes==0){
		return sample,fmt.Errorf("invalid Command queue probe output")
	}
	return sample,nil
}
