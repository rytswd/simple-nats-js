# Commands


## Create Docker Network

```bash
docker network create tester
```

## Create Volumes

```bash
mkdir /tmp/nats-config
mkdir /tmp/nats-vol-1
mkdir /tmp/nats-vol-2
mkdir /tmp/nats-vol-3
```

## Create Configs for JetStream

```bash
cat << EOF > /tmp/nats-config/cluster-server-1.conf
server_name=n1-c1
listen=4222
http=8222

jetstream {
   store_dir=/nats/storage
   max_file=10Gi
}

cluster {
  name: C1
  port: 6222

  routes: [
    nats-route://my-jetstream-server-1:6222
    nats-route://my-jetstream-server-2:6222
    nats-route://my-jetstream-server-3:6222
  ]
}
EOF
```

```bash
cat << EOF > /tmp/nats-config/cluster-server-2.conf
server_name=n2-c1
listen=4222
http=8222

jetstream {
   store_dir=/nats/storage
   max_file=10Gi
}

cluster {
  name: C1
  port: 6222

  routes: [
    nats-route://my-jetstream-server-1:6222
    nats-route://my-jetstream-server-2:6222
    nats-route://my-jetstream-server-3:6222
  ]
}
EOF
```

```bash
cat << EOF > /tmp/nats-config/cluster-server-3.conf
server_name=n3-c1
listen=4222
http=8222

jetstream {
   store_dir=/nats/storage
   max_file=10Gi
}

cluster {
  name: C1
  port: 6222

  routes: [
    nats-route://my-jetstream-server-1:6222
    nats-route://my-jetstream-server-2:6222
    nats-route://my-jetstream-server-3:6222
  ]
}
EOF
```

## First Docker Conatiner
```bash
docker run \
    -it \
    -p 14222:4222 -p 18222:8222 \
    --rm \
    --name my-jetstream-server-1 \
    --network tester \
    --mount type=bind,source=/tmp/nats-vol-1,dst=/nats/storage \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    nats:2.6.1 -c /home/nats-config/cluster-server-1.conf
```


## Second Docker Container

```bash
docker run \
    -it \
    -p 24222:4222 -p 28222:8222 \
    --rm \
    --name my-jetstream-server-2 \
    --network tester \
    --mount type=bind,source=/tmp/nats-vol-2,dst=/nats/storage \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    nats:2.6.1 -c /home/nats-config/cluster-server-2.conf
```


## Third Docker Container

```bash
docker run \
    -it \
    -p 34222:4222 -p 38222:8222 \
    --rm \
    --name my-jetstream-server-3 \
    --network tester \
    --mount type=bind,source=/tmp/nats-vol-3,dst=/nats/storage \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    nats:2.6.1 -c /home/nats-config/cluster-server-3.conf
```


## Set Up Local NATS CLI

```bash
{
    nats context add tester-cluster -s "nats://127.0.0.1:24222"
    nats context select teser-cluster
}
```

## Check NATS CLI Connection

```bash
nats stream ls
```
