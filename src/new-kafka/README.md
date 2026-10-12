Kafka `Broker` can have many partitions on the same node(it tries to spread all leaders across different nodes, but if leader partitions count is more than number of nodes - it can place several leaders on the same node), if number of partitions for topic **payments** is more than number of broker instances.

`LEADER` partition is carried as replica - replication factor 3 means that(1 leader, 2 replicas).
`min.insync.replicas = 3` means that leader and 2 replicas should confirm they had written messages on disk no later than `replica.lag.time.max.ms` since last new message was sent by producer.

### Example:

3 instances of Kafka running:

- 192.168.0.1:9092
- 192.168.0.2:9092
- 192.168.0.3:9092

If you create a topic with `Parttion: 5 and ReplicationFactor: 1` - here is example of approximate partition distribution across 3 Kafka nodes(brokers):

- **Broker1** (192.168.0.1:9092)
  partition-0 (Leader)
  partition-3 (Leader)
- **Broker2** (192.168.0.2:9092)
  partition-1 (Leader)
  partition-4 (Leader)
- **Broker3** (192.168.0.3:9092)
  partition-2 (Leader)

If you create a topic with `Parttion: 5 and ReplicationFactor: 2` - here is example of approximate leader partitions and replicas distribution across 3 Kafka nodes(brokers):

- **Broker1** (192.168.0.1:9092)
  partition-0 (Leader)
  partition-3 (Leader)
  partition-2 (Replica)
  partition-4 (Replica)
- **Broker2** (192.168.0.2:9092)
  partition-1 (Leader)
  partition-4 (Leader)
  partition-3 (Replica)
- **Broker3** (192.168.0.3:9092)
  partition-2 (Leader)
  partition-0 (Replica)
  partition-1 (Replica)

About **acknowledge** factor:

- acks=0: producer doesn't wait of confirmation from any partition
- acks=1: producer waits that message is written in leader-partition
- acks=all: producer waits that message is written in leader and all in-sync replicas
