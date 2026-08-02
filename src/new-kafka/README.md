Kafka `Broker` can have many leader partitions, if number of partitions for topic **payments** is more than number of broker instances.

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
