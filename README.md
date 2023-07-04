# example-external-simulator
This is an example of the external simulator which is one of the features on [kubernetes-sigs/kube-scheduler-simulator](https://github.com/kubernetes-sigs/kube-scheduler-simulator).

# Deploy to your cluster
## Before you begin
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you create one by using [minikube](https://minikube.sigs.k8s.io/docs/tutorials/multi_node/) or some other methods, and then launch it.
```sh
minikube start --nodes 2 -p multinode-demo
```


```sh
export workdir=example1
```

## Build
To deploy for the cluster, you need to build your own scheduler to binary.
```sh
GOOS=linux GOARCH=amd64 go build -a -o bin/example-external-scheduler ${workdir}/main.go
```


Then, build it to docker image.
```sh
docker build -t my-project/example-external-scheduler:1.0 .
```
Upload the built image to any registry. In here, we will push the image directly to minikube.
```sh
minikube image rm my-project/example-external-scheduler:1.0 -p=multinode-demo
minikube image load my-project/example-external-scheduler:1.0 -p=multinode-demo
```

## Configuration
NOTE: I followed [this page](https://kubernetes.io/docs/tasks/extend-kubernetes/configure-multiple-schedulers/) to pass the yaml files via ConfigMap. This minikube is just an example, so you use it in any way that fits your cluster.

### Configure scheduler(Optional)
To enable/disable default-plugin/your-custom-plugin and set some other setting, you need to use `KubeSchedulerConfiguration`.
Please see `${workdir}/configs/configmap-my-scheduler-config.yaml`

### Set your kubeconfig.yaml (Step3 only)
`external scheduler` feature requires us to pass the kubeconfig.yaml of your cluster. This is because the updating process of the pod annotation occurs.
This time, we cheat a little and pass kubeconfig.yaml instead.

WARN: You DO NOT this means in production environment. 

Connect to minikube via ssh.
```sh
minikube ssh -p=multinode-demo
sudo su
```
Copy contents of `/etc/kubernetes/admin.conf` to `my-kubeconfig.yaml` on `example3/configs/configmap-my-scheduler-config.yaml`.

```sh
cat /etc/kubernetes/admin.conf
cat /etc/kubernetes/scheduler.conf
```
Then, overwrite the `server` field with reference to scheduler.conf's one.
```admin.conf
- server: https://control-plane.minikube.internal:8443
+ server: https://192.168.49.2:8443
```

## Deploy
Deploy the our scheduler and configurations.
```sh
kubectl apply -f ${workdir}/configs/configmap-my-scheduler-config.yaml
kubectl apply -f example-external-scheduler.yaml
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

