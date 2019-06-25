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
| message    | string | The message to log
| addDetails | bool   | Append contextual execution information to the log message

## Examples
The below example logs a message 'test message':

```json
{
  "id": "log_message",
  "name": "Log Message",
  "activity": {
    "ref": "qingcloud-flow/plugin/activity/log",
    "input": {
      "message": "test message",
      "addDetails": "false"
    }
  }
}
```