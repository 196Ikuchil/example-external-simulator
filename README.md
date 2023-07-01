# example-external-simulator
This is an example of the external simulator which is one of the features on [kubernetes-sigs/kube-scheduler-simulator](https://github.com/kubernetes-sigs/kube-scheduler-simulator).

# Deploy to your cluster
## Before you begin
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you create one by using [minikube](https://minikube.sigs.k8s.io/docs/tutorials/multi_node/) or some other methods, and then launch it.
```sh
minikube start --nodes 2 -p multinode-demo
```

## Build
To deploy for the cluster, you need to build your own scheduler to binary.
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/example-external-scheduler main.go
```
Then, build it to docker image.
```sh
docker build -t my-project/example-external-scheduler:1.0 .
```
Upload the built image to any registry. In here, we will push the image directly to minikube.
```sh
minikube image load my-project/example-external-scheduler:1.0 -p=multinode-demo
```

## Configure scheduler
To enable/disable default-plugin/your-custom-plugin and set some other setting, you need to use `KubeSchedulerConfiguration`.
Please see under the [configs](/configs) directory.
- config/configmap-my-scheduler-config.yaml

## Deploy
Deploy the our scheduler and configurations.
```sh
kubectl apply -f configs/configmap-my-scheduler-config.yaml
kubectl apply -f configs/example-external-scheduler.yaml
 ```
You can get a list of pods and check the status.
```sh
kubectl get pods --namespace=kube-system
NAME                                     READY   STATUS    RESTARTS      AGE
...
my-scheduler-7748f5c9fb-s59db            1/1     Running   0             20s
...
```

## Scheduling with our scheduler
Deploy pod to use our scheduler as a working check of the scheduler.
```sh
kubectl apply -f ./configs/test-pod.yaml
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

