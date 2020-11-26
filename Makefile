APP = acg

FLGS = -v
.PHONY: build
build:
	go build $(FLGS) ./cmd/$(APP)

.PHONY: run
run:
	go run $(FLGS) ./cmd/$(APP)

.PHOTY: buildnrace
buildnrace:
	go build $(FLGS) -race ./cmd/$(APP)


DBPORT = 27017
DBDIR = $(shell pwd)/mongodata
.PHONY: rundb
rundb:
	docker run -it -p $(DBPORT):$(DBPORT) -v $(DBDIR):/data/db --name acg_db --rm mongo
	# docker run -it -p $(DBPORT):$(DBPORT) --mount type=volume,src=acg_mongodb,dst=/data/db --name acg_db --rm mongo
	# docker run -it -p $(DBPORT):$(DBPORT) -v $(DBDIR):/data/db --name acg_db --rm mongo

.PHONY: clean
clean:
	rm ./$(APP)

.DEFAULT_GOAL := build
