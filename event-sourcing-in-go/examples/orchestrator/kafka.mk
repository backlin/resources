topic-cmd := kafka-topics --zookeeper localhost:2181
groups-cmd := kafka-consumer-groups --bootstrap-server localhost:9092

.PHONY: create-topics delete-topics describe-topics reset-offsets describe-groups



create-topics:
	$(topic-cmd) --create --partitions 2 --replication-factor 1 --topic worker
	# Tracker topic is automatically created with default settings on produce

delete-topics:
	$(topic-cmd) --delete --topic tracker
	$(topic-cmd) --delete --topic worker

describe-topics:
	$(topic-cmd) --describe



reset-offsets:
	$(groups-cmd) --reset-offsets --topic tracker --group tracker
	$(groups-cmd) --reset-offsets --topic worker  --group worker

describe-groups:
	$(groups-cmd) --describe --group tracker
	$(groups-cmd) --describe --group worker

