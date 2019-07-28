<!--
title: MQTT-BROKER
weight: 4705
-->
# MQTT
This trigger allows you to listen to messages on MQTT.

## Installation

### Flogo CLI
```bash
flogo install github.com/qingcloudhx/edge-contrib/trigger/qingcloud-mqtt
```

## Configuration

### Settings:
| Name          | Type   | Description
| :---          | :---   | :---
| url        | string | The broker URL - ***REQUIRED***
 
 #### *sslConfig* Object: 
 
### Handler Settings
### Output: 

| Name    | Type   | Description
| :---    | :---   | :---
| head | object | The message head
| body | bytes | The message body
    
### Reply:

| Name  | Type   | Description
| :---  | :---   | :---
| data  | object | The data recieved

## Example

```json
{
  "id": "mqtt-trigger",
  "name": "Mqtt Trigger",
  "ref": "github.com/qingcloudhx/edge-contrib/trigger/mqtt",
  "settings": {
      "url" : "tcp://localhost:1883"
  },
  "handlers": {
    "settings": {
    }
  }
}
```