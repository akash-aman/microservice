# Kafka Overview

## Core Concepts

### Topics and Events
- On kafka, microservices Can publish events
- microservices Can subscribe to events
- kafka Can store events for later retrieval
- kafka Can replay events later in the event of a failure
- kafka Can transform or process events
- Publisher = Producers
- Consumers = Subscribers
- Events are organized and durabely stored in topics.
- Topic as the name suggest has events of particular topics.
- Topic in kafka are alway multiple publisher & multiple consumers.
- Events are not deleted after consumptions unlike traditional messaging systems, Instead you can define, for how long kafka should retain through a per topic configuration setting, after which all event will be discarded.
- Kafka cluster consist of mulitple broker.
- Broker is nothing but a server. having multiple broker ensures high availability.
- Kafka cluster is managed by kafka zookeeper. 
- Topics are separated into partitions inside kafka clusters.
- Partitions are replicated among the brokers.
- Topic can be read in parallel by many consumers by each partition.
- Producer may also attach a key to a message which directs all message within identical key to that division.
- messages carrying keys are round robin transfered to that partitions, but messages sans keys are progressively added and stored inside partition.
- Using keys you may inforce the processing order of communications in kafka that contain the same key.

###	Cluster Management
- Maintains metadata about the Kafka cluster, including broker information, topic configurations, and partition assignments.
- Keeps track of which brokers are alive and available.
- Zookeeper continuously monitors the Kafka cluster. If a broker goes down, Zookeeper detects it through heartbeats.

###	Leader Election
- Handles the election of partition leaders when a broker goes down or comes back online.
- Ensures high availability by automatically promoting new leaders.
- For partitions previously led by the downed broker, Zookeeper immediately initiates a leader election among the available replicas.
- Kafka clients (producers/consumers) automatically reconnect to the newly elected leader for the affected partitions.

###	Configuration Management
- Stores and manages configuration data for Kafka topics, including replication factor and partition details.
- Allows dynamic updates to topic configurations.

###	State Management
- Maintains the state of Kafka consumers and producers by storing offsets and other transactional information.
- Helps consumers resume from the last processed message after a failure.

###	Fault Tolerance
- Provides fault tolerance by replicating metadata across Zookeeper nodes.
- Ensures data consistency and availability even if one Zookeeper node goes down.

###	Manual or Automated Restart
- Kafka brokers must be restarted manually, or through automation tools like Systemd, Kubernetes, or supervisord.
- Zookeeper does not handle broker restarts or scaling.
