# Simple Setup with KinD

Similar to Docker based setup, we could have Kubernetes to host the Docker image.

For this setup, we use [KinD - Kubernetes-in-Docker](https://kind.sigs.k8s.io/).

## 🛠 Prerequisites

You need the following tools:

- docker
- kubectl
- kind

## 🐾 Steps

### 1. Start local Kubernetes clusters with KinD

```bash
kind create cluster \
    --config https://raw.githubusercontent.com/rytswd/simple-nats-js/main/tools/kind-config/config-4-nodes.yaml
```

<details>

<summary>Details</summary>

To be updated

</details>

---

### 2.1. Prepare for NATS JetStream Cluster Install

```bash
{
    curl -sL -o tmp-simple-nats-js.zip https://github.com/rytswd/simple-nats-js/archive/main.zip
    unzip tmp-simple-nats-js.zip
    cp -r simple-nats-js-main/nats-jetstream-helm .
    rm -rf simple-nats-js-main/ tmp-simple-nats-js.zip
}
```

<details>

<summary>Details</summary>

To be updated

</details>

---

### 2.2 Install NATS JetStream Cluster with Custom Helm Chart

```bash
helm install nats-js-cluster nats-jetstream-helm/
```

<details>

<summary>Details</summary>

To be updated

</details>

---

## 🐾 Reference: Steps with This Repo

### 0. Clone this repository

```bash
$ pwd
/some/path/at

$ git clone https://github.com/rytswd/simple-nats-js

$ cd simple-nats-js
```

From here on, all the steps are assumed to be run from `/some/path/at/simple-nats-jso`.

<details>

<summary>Details</summary>

To be updated

</details>

---

### 1. Start local Kubernetes clusters with KinD

```bash
kind create cluster \
    --config ./tools/kind-config/config-4-nodes.yaml
```

<details>

<summary>Details</summary>

To be updated

</details>

---

### 2. Install NATS JetStream Cluster with Custom Helm Chart

```bash
helm install nats-js-cluster nats-jetstream-helm
```

<details>

<summary>Details</summary>

To be updated

</details>

---
