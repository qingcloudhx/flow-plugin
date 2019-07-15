<!--
title: Log
weight: 4615
-->

# Log
This activity allows you to write log messages.

## Installation

### Flogo CLI
```bash
flogo install qingcloud-flow/plugin/activity/log
```

## Configuration

### Input:
| Name       | Type   | Description
|:---        | :---   | :---    
| head    | map | The message to log
| body | bytes   | body

## Examples
The below example logs a message 'test message':

```json
{
  "id": "log_message",
  "name": "Log Message",
  "activity": {
    "ref": "qingcloud-flow/plugin/activity/log",
    "input": {
      "head": "test message",
      "body": "false"
    }
  }
}
```