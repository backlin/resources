all: event.pb.go bin/cli bin/tracker bin/worker

clean:
	rm -f event.pb.go
	rm -rf bin/*

event.pb.go: event.proto orchestrator.go
	protoc --gogofaster_out=. $<

bin/cli: cli/main.go | event.pb.go
	go build -o bin/cli $^

bin/tracker: tracker/main.go tracker/tracker.go tracker/state.go | event.pb.go
	go build -o bin/tracker $^

bin/worker: worker/main.go worker/worker.go | event.pb.go
	go build -o bin/worker $^
