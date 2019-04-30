# flogo-components
Collection of custom built flogo components

## Components

### Activities
* [amqp](activity/amqp): Publishes a message on AMQP queue or topic
* [blewrite](activity/blewrite): Write data to a Bluetooth BLE Device
* [checkiban](activity/checkiban): Validates an International Bank Account Number
* [combine](activity/combine): Combine separate parts into a single string
* [eftl](activity/eftl): Send message to eFTL
* [encoder](activity/encoder): Encode en decode strings using Base64, Base32
* [filter](activity/filter): Filter out unwanted data
* [getjson](activity/getjson): Retrieve specific elements from a JSON structure
* [jsontodata](activity/jsontodata): Convert JSON string to data object
* [kafka](activity/kafka): Send message to Kafka
* [matchresponse](activity/matchresponse): response to input text based on seach data
* [mqtt](activity/mqtt): Publish MQTT Message
* [redis](activity/redis): Interact with Redis keyspace
* [replace](activity/replace): Find and replace characters in up to 8 strings
* [slack](activity/slack): Send message to Slack
* [splitjson](activity/splitjson): Splits JSON structure into separate name-value pairs
* [splitpath](activity/splitpath): Splits a path into separate parts
* [statechange](activity/statechange): Detects state change for up to eight inputs
* [systeminfo](activity/systeminfo): Retrieve System Information
* [tcmpub](activity/tcmpub): Publish message on TIBCO Cloud Messaging
* [throttle](activity/throttle): Throttle data based on interval
* [udp](activity/udp): Send message over UDP
* [wsmessage](activity/wsmessage): Send message over WebSockets

### Triggers
* [amqp](trigger/amqp):Receive messages from AMQP
* [blemaster](trigger/blemaster):Receive BLE Data (Master)
* [eftl](trigger/eftl): Receive messages from TIBCO eFTL
* [gpio](trigger/gpio): Use GPIO pin to start flows
* [kafka](trigger/kafka): Receive messages from Kafka
* [mqtt2](trigger/mqtt2): Receive messages from MQTT
* [slack](trigger/slack): Receive messages from Slack
* [tcmsub](trigger/tcmsub): Receive messages from TIBCO Cloud Messaging
* [timer2](trigger/timer2): Start flow from timer
* [wsserver](trigger/wsserver): Receive messages on built-in WebSocket server


## Related Information
These components are built for TIBCO flogo.
Please check (http://www.flogo.io/) for more information.
On Github: (https://github.com/TIBCOSoftware/flogo)
