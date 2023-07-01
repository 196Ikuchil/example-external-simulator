package example1

import (
	"os"

	"k8s.io/component-base/cli"
	"k8s.io/component-base/logs"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"github.com/196Ikuchil/example-external-simulator/plugins/communicating"
	"github.com/196Ikuchil/example-external-simulator/plugins/nodenumber"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(communicating.Name, communicating.New),
		app.WithPlugin(nodenumber.Name, nodenumber.New),
	)

	logs.InitLogs()
	defer logs.FlushLogs()

	code := cli.Run(command)
	os.Exit(code)
}
