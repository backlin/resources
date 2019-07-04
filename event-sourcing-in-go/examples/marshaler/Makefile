all: item.pb.go

clean:
	rm -f item.pb.go

item.pb.go: item.proto
	protoc --gogofaster_out=. $<
