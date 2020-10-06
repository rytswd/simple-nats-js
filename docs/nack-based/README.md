# Simple Setup with [NACK](https://github.com/nats-io/nack)

This setup uses [NACK (NATS Controllers for Kubernetes)](https://github.com/nats-io/nack).

NACK works as a K8s controller.

For this setup, we use [KinD - Kubernetes-in-Docker](https://kind.sigs.k8s.io/) to test locally.

## Prerequisites

You need the following tools:

- docker
- kubectl
- kind
- helm

## Steps

### 1. Start local Kubernetes cluster with KinD

```bash
$ {
    mkdir /tmp/nack-js
    kind create cluster --config ./tools/kind-config/config-2-nodes.yaml --name kind-nats
}
```

<details>
<summary>Details</summary>

You can skip this step if you are using existing Kubernetes cluster. This step is only to demonstrate how to get started with absolutely no cluster runnig right now.

Also, while at this step, this creates a temporary directory to store K8s definition files which we will be creating later on.

</details>

---

### 2. STEP TO BE UPDATED

```bash
$ {
    # 1.
    # Deploy StatefulSet of NATS Servers that are not JetStream enabled
    kubectl apply -f https://raw.githubusercontent.com/nats-io/k8s/master/nats-server/simple-nats.yml

    # 2.
    # Deploy StatefulSet of JetStream enabled NATS Server as a leafnode connection
    kubectl apply -f https://raw.githubusercontent.com/nats-io/k8s/master/nats-server/nats-js-leaf.yml

    # 3.
    # Deploy JetStream specific CRDs and RBAC setup
    kubectl apply \
      -f https://raw.githubusercontent.com/nats-io/nack/main/deploy/crds.yml \
      -f https://raw.githubusercontent.com/nats-io/nack/main/deploy/rbac.yml

    # 4.
    # Deploy JetStream-Controller
    kubectl apply -f https://raw.githubusercontent.com/nats-io/nack/main/deploy/deployment.yml

    # 5. (Not required)
    # Deploy NATS management utility for debugging
    kubectl apply -f https://nats-io.github.io/k8s/tools/nats-box.yml
}
```

<details>
<summary>Details</summary>

1. NATS Server StatefulSet:  
   `curl -sL https://raw.githubusercontent.com/nats-io/k8s/master/nats-server/simple-nats.yml | less`  
   This creates a cluster of NATS Servers, which are not JetStream enabled. _TODO: Check why this is needed._
1. JetStream enabled NATS Server:  
   `curl -sL https://raw.githubusercontent.com/nats-io/k8s/master/nats-server/nats-js-leaf.yml | less`  
   This creates a StatefulSet with 1 replica of JetStream enabled NATS Server. As of writing (Oct 2020), the JetStream clustering is not supported, and thus the replica count is set to 1.
1. JetStream Controller:  
   `curl -sL https://raw.githubusercontent.com/nats-io/nack/main/deploy/crds.yml | less`  
   `curl -sL https://raw.githubusercontent.com/nats-io/nack/main/deploy/rbac.yml | less`  
   `curl -sL https://raw.githubusercontent.com/nats-io/nack/main/deploy/deployment.yml | less`
   CRD (Custom Resource Definition) and RBAC (Role Based Access Control) are the definitions users will be interacting with.

</details>

---

### 3. STEP TO BE UPDATED

```bash
$ {
    cat << EOF > /tmp/nack-js/stream.conf
---
apiVersion: jetstream.nats.io/v1beta1
kind: Stream
metadata:
  name: mystream
spec:
  name: mystream
  subjects: ["orders.*"]
  storage: memory
  maxAge: 1h
EOF

    kubectl apply --context kind-kind-nats -f /tmp/nack-js/stream.conf
}
```

<details>
<summary>Details</summary>

To be updated

</details>

---

### 4. STEP TO BE UPDATED

```bash
$ {
    cat << EOF > /tmp/nack-js/push-consumer.conf
---
apiVersion: jetstream.nats.io/v1beta1
kind: Consumer
metadata:
  name: my-push-consumer
spec:
  streamName: mystream
  durableName: my-push-consumer
  deliverSubject: my-push-consumer.orders
  deliverPolicy: last
  ackPolicy: none
  replayPolicy: instant
EOF

    cat << EOF > /tmp/nack-js/pull-consumer.conf
---
apiVersion: jetstream.nats.io/v1beta1
kind: Consumer
metadata:
  name: my-pull-consumer
spec:
  streamName: mystream
  durableName: my-pull-consumer
  deliverPolicy: all
  filterSubject: orders.received
  maxDeliver: 20
  ackPolicy: explicit
EOF

    kubectl apply --context kind-kind-nats \
        -f /tmp/nack-js/push-consumer.conf \
        -f /tmp/nack-js/pull-consumer.conf
}
```

<details>
<summary>Details</summary>

To be updated

</details>

---

### 5. STEP TO BE UPDATED

For testing, get into the `nats-box` Pod.

```bash
$ kubectl exec -it nats-box -- /bin/sh -l
```

The following commands set the running context.

```bash
$ {
    nats context save jetstream -s nats://nats:4222
    nats context select jetstream
}
```

With the above context set, you can now start publishing message to JetStream enabled NATS Server.

```bash
$ nats pub orders.received "order 1"
$ nats pub orders.received "order 2"
$ nats pub orders.other "other order ABCDEF"
```

<details>
<summary>Details</summary>

To be updated

</details>

---
