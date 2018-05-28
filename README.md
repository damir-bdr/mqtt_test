# mqtt_test

## Installing Mosquitto MQTT Messaging Broker
* [How to Install and Secure the Mosquitto MQTT Messaging Broker on Ubuntu 16.04](https://www.digitalocean.com/community/tutorials/how-to-install-and-secure-the-mosquitto-mqtt-messaging-broker-on-ubuntu-16-04)

## Build example

    $ make

## Run example
in first terminal start mqttex_server with arguments

    $ ./mqttex_server aaa bbb ccc

in second terminal start mqttex_client with arguments

    $ ./mqttex_client aaa bbb ccc

## Dependencies
This example uses `paho.mqtt.golang` library. To install:

    $ make deps

During the installation of `paho.mqtt.golang` some problems with packages `golang.org/x/net/proxy` and `golang.org/x/net/websocket` are possible. To solve it use commands:

    $ cd $GOPATH/src
    $ git clone https://github.com/golang/net.git golang.org/x/net
