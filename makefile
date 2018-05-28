GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get

all: build
build:
	$(GOBUILD) -o mqttex_client client/mqttex_client.go
	$(GOBUILD) -o mqttex_server server/mqttex_server.go
clean:
	rm -f mqttex_client
	rm -f mqttex_server
deps:
	$(GOGET) github.com/eclipse/paho.mqtt.golang
	$(GOGET) github.com/satori/go.uuid
