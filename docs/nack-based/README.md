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
$ mkdir /tmp/nack-js

$ {
    kind create cluster --config ./tools/kind-config/config-2-nodes.yaml --name kind-nats
}
```

<details>
<summary>Details</summary>

To be updated

</details>

---

### 2. STEP TO BE UPDATED

```bash
# Creates cluster of NATS Servers that are not JetStream enabled
$ kubectl apply -f https://raw.githubusercontent.com/nats-io/k8s/master/nats-server/simple-nats.yml

# Creates NATS Server with JetStream enabled as a leafnode connection
$ kubectl apply -f https://raw.githubusercontent.com/nats-io/k8s/master/nats-server/nats-js-leaf.yml


$ kubectl apply -f https://raw.githubusercontent.com/nats-io/nack/main/deploy/crds.yml
$ kubectl apply -f https://raw.githubusercontent.com/nats-io/nack/main/deploy/rbac.yml
$ kubectl apply -f https://raw.githubusercontent.com/nats-io/nack/main/deploy/deployment.yml
```

<details>
<summary>Details</summary>

To be updated

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
