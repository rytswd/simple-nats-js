# Simple Setup of NATS JetStream

**IMPORTANT NOTE**

This is a personal note of how I understand NATS JetStream offering.

Most of the information is based on my understanding, and if any wording or information does not match with the official documentation, please file an issue or PR.

## Prerequisite

You need:

- Go
- Docker

## Simple Setup with Docker only

The following is to achieve persistent JetStream setup. The setup only considers local testing scenarios, but should be easily converted to more cloud native approach.

This does not any of Go code in this repo. (It will be added as a separate document.)

### 1. Local directory setup

Create directory on local machine so that we can persist the config and data.

```bash
$ {
    mkdir /tmp/nats-config  # For NATS and JetStream configurations
    mkdir /tmp/nats-vol     # For NATS data and objects
}
```

<details>
<summary>Details</summary>

There are 2 different items to be considered for persistent setup.

- Config
- Data

#### Config

The "Config" is further split into 2 phases: (1.) NATS server startup configuration, and (2.) JetStream "Stream" and "Consumer" configurations.

At (1.) NATS server startup, you can provide your config to NATS to change where to store JetStream related data, and other NATS feature setup. This is used at Step#5. The (2.) JetStream "Stream" and "Consumer" configurations are discussed a bit more in later section. Essentially, you can consider these as a one-time setup, and we will be using them for the first time setup.

#### Data

The "Data" is the actual data stored in the NATS server. This refers to the actual messages sent to the NATS server, how many are ack'ed, etc. The NATS server handles JetStream's "Stream" and "Consumer" concepts, and if you choose to create "Stream" with File storage, these can be persisted at the NATS server with files. We are creating a directory `/tmp/nats-vol/` to store all the NATS JetStream data, so we can restart the NATS server without losing data or configuration.

</details>

---

### 2. Create NATS cocnfiguration with JetStream enabled

Create NATS configuration.

```bash
$ cat << EOF > /tmp/nats-config/jetstream.conf
jetstream {
    store_dir: "/data/jsm/"
}
EOF
```

<details>
<summary>Details</summary>

This `/tmp/nats-config/jetstream.conf` is the simplest setup.

It tells the NATS server to use `/data/jsm/` directory to store the JetStream related data. This means that any data / config will be stored under this directory when file storage is used, and also recovers from the files in this directory when the NATS server starts up.

</details>

---

### 3. Create JetStream "Stream" configuration

```bash
$ cat << EOF > /tmp/nats-config/jetstream-stream.json
{
  "name": "AnotherStream",
  "subjects": ["xyz.*"],
  "retention": "limits",
  "max_consumers": -1,
  "max_msgs": -1,
  "max_bytes": -1,
  "max_age": 0,
  "max_msg_size": -1,
  "storage": "file",
  "discard": "old",
  "num_replicas": 1,
  "duplicate_window": 120000000000
}
EOF
```

<details>
<summary>Details</summary>

This is the JetStream "Stream" configuration. We will be creating the "Stream" later using this file.

This file does not need to be persisted. This is saved under `/tmp/nats-config/jetstream-stream.json` just for the ease of the setup. It can be done within Docker image instead if you don't need to hold on to the original configuration file.

Also note that the configuration will be persisted at the NATS server, so it is easy to reccreate the config file from the running NATS server.

---

_TODO: Add reference for each attribute_

</details>

---

### 4. Create JetStream "Consumer" configuration

```bash
$ cat << EOF > /tmp/nats-config/jetstream-consumer.json
{
  "durable_name": "SomeConsumer",
  "deliver_subject": "pull",
  "deliver_policy": "all",
  "ack_policy": "explicit",
  "ack_wait": 2000000000,
  "max_deliver": -1,
  "filter_subject": "xyz.*",
  "replay_policy": "instant"
}
EOF
```

<details>
<summary>Details</summary>

This is the JetStream "Consumer" configuration. We will be creating the "Consumer" later using this file.

This file does not need to be persisted. This is saved under `/tmp/nats-config/jetstream-consumer.json` just for the ease of the setup. It can be done within Docker image instead if you don't need to hold on to the original configuration file.

Also note that the configuration will be persisted at the NATS server, so it is easy to reccreate the config file from the running NATS server.

---

_TODO: Add reference for each attribute_

</details>

---

### 5. Start NATS server with JetStream

Use a separate terminal for this step, as you want to keep this process running.

Also, you will be restarting this later on, so it is better to have it on separate terminal rather than running it on background.

```bash
$ docker run \
    -it \
    -p 4222:4222 \
    --name my-jetstream-server \
    --mount type=bind,source=/tmp/nats-vol,dst=/data/jsm \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    synadia/jsm:nightly server -c /home/nats-config/jetstream.conf \
    ; docker container rm my-jetstream-server
```

<details>
<summary>Details</summary>

Docker command reference:

- `-it`: For interactive process
- `-p 4222:4222`: Use local port `4222` on the Docker host, and map to `4222` on Docker container
- `--name my-jetstream-server`: Set a name so that we can use it to link another Docker container later
- `--mount type=bind,source=/tmp/nats-vol,dst=/data/jsm`: Volume mounting for NATS data and objects
- `--mount type=bind,source=/tmp/nats-config,dst=/home/nats-config`: Volume mounting for NATS config
- `synadia/jsm:nightly`: Docker image we are using
- `server`: Docker CMD - this is handled by `synadia/jsm` image with `entrypoint.sh`
- `-c /home/nats-config/jetstream.conf`: Docker ARG - this is handled by `synadia/jsm` image
- `; docker container rm my-jetstream-server`: When stopping the ccontainer, remove the container, so that you can easily restart using the same container name

</details>

---

### 6. Start NATS client

```bash
$ docker run \
    -it \
    --link my-jetstream-server \
    --env NATS_URL=my-jetstream-server:4222 \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    synadia/jsm:latest
```

<details>
<summary>Details</summary>

Docker command reference:

- `-it`: For interactive process
- `--link my-jetstream-server`: Link to running NATS server
- `--env NATS_URL=my-jetstream-server:4222`: Ensure connection is made to the linked container
- `--mount type=bind,source=/tmp/nats-config,dst=/home/nats-config`: Volume mounting for NATS config
- `synadia/jsm:latest`: Docker image we are using - if no argument is provided, it goes to interactive shell by default

</details>

---

### 7. Create JetStream "Stream" and "Consumer

[ From the interactive shell from the Step #6 ]

```bash
$ nats str add --config=/home/nats-config/jetstream-stream.json
```

And then,

```bash
$ nats con add AnotherStream --config=/home/nats-config/jetstream-consumer.json
```

<details>
<summary>Details</summary>

Firstly, create JetStream "Stream":

- `nats`: CLI provided in `synadia/jsm` Docker image
- `str`: Short for `stream`, handles JetStream "Stream" related setup
- `add`: Add a new "Stream"
- `--config=/home/nats-config/jetstream-stream.json`: Use the JSON configuration from Step#3.

Then, create JetStream "Consumer" - this needs to map to an exisitng "Stream".

- `nats`: CLI provided in `synadia/jsm` Docker image
- `con`: Short for `consumer`, handles JetStream "Consumer" related setup
- `add`: Add a new "Consumer"
- `AnotherStream`: "Stream" name "Consumer" should connect to. This is the same name used in Step#3.
- `--config=/home/nats-config/jetstream-consumer.json`: Use the JSON configuration from Step#4.

_NOTE_: For both commands, by omitting `--config` option, you can go into interactive setup mode.

## </details>

---

### 8. Verify setup

[ From the interactive shell from the Step #6 ]

Publish some data

```bash
$ nats pub xyz.test "some random data"
10:29:03 Published 16 bytes to "xyz.test"
```

Check "Stream" has the data and configuration in place

```bash
$ nats str info AnotherStream
```

Check "Consumer"

```bash
$ nats con info AnotherStream SomeConsumer
```

<details>
<summary>Details</summary>

The subject `xyz.test` matches the "Stream" and "Consumer" setup of `xyz.*`.
"Stream" name and "Consumer" name are defined in the JSON above.

</details>

---

### 9. Verify file storage persistence

Using another terminal, verify that file storage holds the JetStream data.

```bash
$ tree /tmp/nats-vol
```

This will result in output like the following

```
/tmp/nats-vol
└── jetstream
    └── $G
        └── streams
            └── AnotherStream
                ├── meta.inf
                ├── meta.sum
                ├── msgs
                │   ├── 1.blk
                │   └── 1.idx
                └── obs
                    └── SomeConsumer
                        ├── meta.inf
                        ├── meta.sum
                        └── o.dat

7 directories, 7 files
```

<details>
<summary>Details</summary>

You can see how there is a new directory `jetstream` is created under the mounted volume.

- `$G` is the account used by JetStream (to be confirmed)
- "Stream" is stored under `streams` directory
- "Consumer" for a given Stream is stored under `obs` directory under `streams` directory

</details>

---

### 10. Stop NATS Server

Using the terminal used for Step#5, simply kill the server with `ctrl-C`.

You can try sending NATS client from Step#6 to confirm the client lost connection to the server.

<details>
<summary>Details</summary>

At this point, NATS client won't be able to establish connection for any commands.

However, the data stored in "Stream", the "Consumer" ack'ed list, etc. are all kept in the file in the local machine at `/tmp/nats-vol`.

</details>

---

### 11. Restart NATS Server

Simply use the same command as Step#5

```bash
$ docker run \
    -it \
    -p 4222:4222 \
    --name my-jetstream-server \
    --mount type=bind,source=/tmp/nats-vol,dst=/data/jsm \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    synadia/jsm:nightly server -c /home/nats-config/jetstream.conf \
    ; docker container rm my-jetstream-server
```

---

### 12. Verify setup

[ From the interactive shell from the Step #6 ]

Check "Stream" has the data and configuration in place

```bash
$ nats str info AnotherStream
```

Check "Consumer"

```bash
$ nats con info AnotherStream SomeConsumer
```

---

### 13. Clean up

The clean-up is straightforward.

- Stop NATS server with `ctrl-C` on Step#11
- Stop NATS client with `ctrl-D` on Step#6
- Remove config and data files with following commands

```bash
$ {
    rm -rf /tmp/nats-config
    rm -rf /tmp/nats-vol
}
```
