# example-external-simulator
This is an example of the external simulator which is one of the features on [kubernetes-sigs/kube-scheduler-simulator](https://github.com/kubernetes-sigs/kube-scheduler-simulator).

# Deploy to your cluster
## Before you begin
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you create one by using [minikube](https://minikube.sigs.k8s.io/docs/tutorials/multi_node/) or some other methods, and then launch it.
```sh
minikube start --nodes 2 -p multinode-demo
```

Here are three examples. Please choose one you would like to check.
```sh
export workdir=example1
```

## Configure scheduler
To enable/disable default-plugin/your-custom-plugin and set some other setting, you need to use `KubeSchedulerConfiguration`.
Please see `${workdir}/configs/kube-scheduler-config.yaml`

### Set your kubeconfig.yaml path
In order for the scheduler to communicate with the control plane, the absolute path to kubeconfig.yaml must be specified in the KubeSchedulerConfiguration file.
In general, the kubeconfig file for minikube is located at `~/.kube/config`.

So, rewrite this field.
```yaml
clientConnection:
  kubeconfig: <absolute path to kubeconfig of minikube>
```

## Run scheduler



```sh
go run ./${workdir}/main.go --config ./${workdir}/configs/kube-scheduler-config.yaml
```

## Scheduling with our scheduler
Deploy pod to use our scheduler as a working check of the scheduler.
```sh
kubectl apply -f ${workdir}/configs/test-pod1.yaml
kubectl get events
LAST SEEN   TYPE     REASON                    OBJECT                            MESSAGE
...
25s         Normal   Scheduled                 pod/annotation-second-scheduler   Successfully assigned default/annotation-second-scheduler to multinode-demo-m02
24s         Normal   Pulled                    pod/annotation-second-scheduler   Container image "registry.k8s.io/pause:2.0" already present on machine
24s         Normal   Created                   pod/annotation-second-scheduler   Created container pod-with-second-annotation-container
24s         Normal   Started                   pod/annotation-second-scheduler   Started container pod-with-second-annotation-container
```

# References
- [Configure Multiple Schedulers](https://kubernetes.io/docs/tasks/extend-kubernetes/configure-multiple-schedulers/)

