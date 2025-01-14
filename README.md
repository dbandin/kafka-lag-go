# Kafka Lag Go

Kafka Lag Calculator is a lightweight, stateless application designed to calculate Kafka consumer group lag in both offsets and seconds. It efficiently processes Kafka consumer group data using a pipeline pattern implemented with Go’s goroutines and channels, ensuring high performance and scalability. The application employs consistent hashing to distribute work evenly among multiple nodes, making it suitable for deployment in distributed environments.

The method for estimating lag in seconds used by Kafka Lag Calculator is inspired by the implementation found in the [Kafka Lag Exporter](https://github.com/seglo/kafka-lag-exporter) project, which is currently in maintenance mode. Kafka Lag Calculator uses both interpolation and extrapolation techniques to approximate the lag in seconds for each partition within a Kafka consumer group.

This strategy provides a reasonable estimate of the time delay between message production and consumption, even in cases where the offset data is sparse or not evenly distributed. However, it is important to note that this is an approximation and not an exact measure, making it useful for gaining insights into lag trends rather than precise calculations.

## Features

- Offset Lag Calculation: Accurately computes the difference between the latest produced offset and the committed offset for each partition, topic, and consumer group.
- Time Lag Calculation: Calculates the lag in seconds, offering insight into the time delay between message production and consumption.
- Max Lag Calculation: Determines the maximum lag in offsets and seconds at both the consumer group and topic levels, providing a clear view of the most delayed parts of your system.
- Stateless Service: Designed to be stateless, allowing for easy scaling and distribution of work across multiple instances or nodes.
- Consistent Hashing: Uses consistent hashing to split the workload among nodes, ensuring an even distribution of data processing tasks.
- Docker & Kubernetes Ready: Can be deployed as a Docker container, making it easy to run in Kubernetes or other container orchestration systems, as well as standalone.
- Pipeline Usage in Redis: The application utilizes Redis pipelines to reduce round trip times and improve overall performance when interacting with Redis. By batching multiple commands together and sending them in a single network request, the application minimizes the overhead of multiple round trips.


# Architecture Overview

This project is designed to efficiently manage data flow and processing tasks using Go's concurrency model. For a detailed explanation of the architecture, please refer to the [Architecture Documentation](docs/Architecture.md).

## Performance

In performance tests against other open-source solutions, Kafka Lag Calculator has demonstrated significant improvements, processing Kafka consumer group lag up to 18 times faster. These improvements are achieved through efficient use of Go’s concurrency model, consistent hashing for load distribution, and the lightweight, stateless architecture of the application.

### Prerequisites

- Docker (optional, it can run in standalone as well)
- A running Redis instance (for now the only supported storage)
- A running Kafka cluster

## Platform Compatibility

These images are available on Docker Hub and can be pulled and run on systems with the corresponding platforms:

- Linux/amd64
- Linux/arm64
- Linux/arm/v7

If you require Kafka Lag Calculator to run on other architectures or platforms, you can easily compile the application from the source code to suit your needs. The flexibility of the Go programming language ensures that Kafka Lag Calculator can be built and run virtually anywhere Go is supported.


## Installation

## Deployment Options

To start monitoring Kafka consumer group lag with Kafka Lag Monitor, you have the following deployment options:

1. Standalone Mode:
  - Run the service directly on your machine by compiling the source code and executing the binary.

2. Containerized Deployment:
  - Kafka Lag Monitor is fully containerized, making it easy to deploy in containerized environments like Kubernetes (K8s).
  - Build the Docker container yourself or pull the pre-built image from Docker Hub.


### Step 1: Clone the Repository

Begin by cloning the Kafka Lag Monitor repository:

```bash
git clone https://github.com/sciclon2/kafka-lag-go.git
cd kafka-lag-go
```

### Step 2: Build the Docker Image

Use the following command to build the Docker image for Kafka Lag Monitor:

```bash
docker build -t kafka-lag-go:latest .
```

### Step 3: Prepare the Configuration File

Kafka Lag Monitor requires a YAML configuration file. Create a file named `config.yaml` and customize it based on your environment. Below is an example configuration:

```yaml
prometheus:
  metrics_port: 9090
  labels:
    env: production
    service: kafka-lag-calculator
    
kafka:
  brokers:
    - "broker1:9092"
    - "broker2:9092"
  client_request_timeout: "30s"
  metadata_fetch_timeout: "5s"
  consumer_groups:
    whitelist: ".*"
    blacklist: "test.*"
  ssl:
    enabled: true
    client_certificate_file: "/path/to/cert.pem"
    client_key_file: "/path/to/key.pem"
    insecure_skip_verify: true
  sasl:
    enabled: true
    mechanism: "SCRAM-SHA-512"
    user: "kafkaUser"
    password: "kafkaPassword"

storage:
  type: "redis"
  redis:
    address: "redis-server"
    port: 6379
    client_request_timeout: "60s"
    client_idle_timeout: "5m"
    retention_ttl_seconds: 7200

app:
  cluster_name: "kafka-cluster"
  iteration_interval: "30s"
  num_workers: 10
  log_level: "info"
  health_check_port: 8080
  health_check_path: "/healthz"
```

### Step 4: Run the Docker Container

After building the image and preparing the configuration file, run the Kafka Lag Monitor Docker container:

```bash
docker run --rm -v /path/to/config.yaml:/app/config.yaml kafka-lag-monitor:latest --config-file /app/config.yaml
```

Replace `/path/to/config.yaml` with the actual path to your configuration file.


### Downlaod the image 
Kafka Lag Calculator is available as a Docker image, making it easy to deploy in containerized environments like Kubernetes.
You can download the Docker image from Docker Hub using the following command:
```
docker pull sciclon2/kafka-lag-go
```


## Configuration

The Kafka Lag Monitor requires a YAML configuration file to customize its behavior. Below is a description of the available configuration options:

- `prometheus.metrics_port`: The port on which Prometheus metrics will be exposed.
- `prometheus.labels`: Additional labels to be attached to the exported metrics.
- `kafka.brokers`: The list of Kafka broker addresses.
- `kafka.client_request_timeout`: The timeout for Kafka client requests.
- `kafka.metadata_fetch_timeout`: The timeout for fetching Kafka metadata.
- `kafka.consumer_groups.whitelist`: A regular expression to whitelist consumer groups.
- `kafka.consumer_groups.blacklist`: A regular expression to blacklist consumer groups.
- `kafka.ssl.enabled`: Whether SSL/TLS is enabled for Kafka connections.
- `kafka.ssl.client_certificate_file`: The path to the client certificate file for SSL/TLS.
- `kafka.ssl.client_key_file`: The path to the client key file for SSL/TLS.
- `kafka.ssl.insecure_skip_verify`: Whether to skip SSL/TLS verification.
- `kafka.sasl.enabled`: Whether SASL authentication is enabled for Kafka connections.
- `kafka.sasl.mechanism`: The SASL mechanism to use.
- `kafka.sasl.user`: The username for SASL authentication.
- `kafka.sasl.password`: The password for SASL authentication.
- `storage.type`: The type of storage backend to use.
- `storage.redis.address`: The address of the Redis server for Redis storage.
- `storage.redis.port`: The port of the Redis server for Redis storage.
- `storage.redis.client_request_timeout`: The timeout for Redis client requests.
- `storage.redis.client_idle_timeout`: The idle timeout for Redis clients.
- `storage.redis.retention_ttl_seconds`: The time-to-live (TTL) for Redis keys.
- `app.cluster_name`: The name of the Kafka cluster.
- `app.iteration_interval`: The interval at which the lag monitor iterates over consumer groups.
- `app.num_workers`: The number of worker goroutines to use.
- `app.log_level`: The log level for the lag monitor.
- `app.health_check_port`: The port on which the health check endpoint will be exposed.
- `app.health_check_path`: The path of the health check endpoint.

Please refer to the `config.go` file for more details on each configuration option.



## Health Check

The health check feature in Kafka Lag Monitor monitors the accessibility of both Kafka and Redis. It ensures that the application can successfully connect to and interact with these components. By periodically checking the health of Kafka and Redis, the health check feature provides insights into the overall availability and reliability of the system.

The health check endpoint is exposed on the specified port (`app.health_check_port`) and path (`app.health_check_path`) in the configuration file. By accessing this endpoint, you can obtain information about the health status of Kafka Lag Monitor.


## Prometheus Metrics

This application exposes a comprehensive set of Prometheus metrics that provide insights into the lag experienced by Kafka consumer groups at both the group and topic levels. These metrics help monitor the health and performance of your Kafka consumers.

### Metrics Overview

- `kafka_consumer_group_lag_in_offsets (group, topic, partition)`: The lag in offsets for a specific partition within a Kafka topic for a consumer group. (Type: Gauge)
- `kafka_consumer_group_lag_in_seconds (group, topic, partition)`: The lag in seconds for a specific partition within a Kafka topic for a consumer group. (Type: Gauge)
- `kafka_consumer_group_max_lag_in_offsets (group)`: The maximum lag in offsets across all topics and partitions within a Kafka consumer group. (Type: Gauge)
- `kafka_consumer_group_max_lag_in_seconds (group)`: The maximum lag in seconds across all topics and partitions within a Kafka consumer group. (Type: Gauge)
- `kafka_consumer_group_topic_max_lag_in_offsets (group, topic)`: The maximum lag in offsets for a specific Kafka topic within a consumer group. (Type: Gauge)
- `kafka_consumer_group_topic_max_lag_in_seconds (group, topic)`: The maximum lag in seconds for a specific Kafka topic within a consumer group. (Type: Gauge)
- `kafka_consumer_group_sum_lag_in_offsets (group)`: The sum of lag in offsets across all topics and partitions within a Kafka consumer group. (Type: Gauge)
- `kafka_consumer_group_sum_lag_in_seconds (group)`: The sum of lag in seconds across all topics and partitions within a Kafka consumer group. (Type: Gauge)
- `kafka_consumer_group_topic_sum_lag_in_offsets (group, topic)`: The sum of lag in offsets for a specific Kafka topic within a consumer group. (Type: Gauge)
- `kafka_consumer_group_topic_sum_lag_in_seconds (group, topic)`: The sum of lag in seconds for a specific Kafka topic within a consumer group. (Type: Gauge)
- `kafka_total_groups_checked`: The total number of Kafka consumer groups checked in each iteration. (Type: Gauge)
- `kafka_iteration_time_seconds`: The time taken to complete an iteration of checking all Kafka consumer groups. (Type: Gauge)


Once the container is running, Kafka Lag Monitor will expose Prometheus metrics on the port specified in the configuration file (`metrics_port`). You can access the metrics at:

```
http://<docker-host-ip>:<metrics_port>/metrics
```

## Next Steps
Please check issues section.
For more details on usage and advanced configuration, refer to the full documentation (coming soon).

## License

This project is licensed under the Apache License 2.0. You may obtain a copy of the License at:

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.