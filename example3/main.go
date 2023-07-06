package main

import (
	"fmt"
	"os"

	"k8s.io/component-base/cli"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"

	"sigs.k8s.io/kube-scheduler-simulator/simulator/pkg/externalscheduler"
	"github.com/196Ikuchil/example-external-simulator/plugins/communicating"
	"github.com/196Ikuchil/example-external-simulator/plugins/nodenumber"

)

func main() {
	command, cancelFn, err  := externalscheduler.NewSchedulerCommand(
		externalscheduler.WithPlugin(communicating.Name, communicating.New),
		externalscheduler.WithPlugin(nodenumber.Name, nodenumber.New),
	)
	if err != nil {
		klog.Info(fmt.Sprintf("failed to build the scheduler command: %+v", err))
		os.Exit(1)
	}

	logs.InitLogs()
	defer logs.FlushLogs()

	code := cli.Run(command)

	cancelFn()
	os.Exit(code)
}
