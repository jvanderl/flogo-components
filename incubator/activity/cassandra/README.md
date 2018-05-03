# Cassandra Query
This activity provides your flogo application the ability to query a Cassandra database.


## Installation

```bash
flogo add activity github.com/jvanderl/flogo-components/incubator/activity/cassandra
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/incubator/activity/cassandra
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "ClusterIP",
      "type": "string",
      "value": "localhost"
    },
	{
      "name": "Keyspace",
      "type": "string",
      "value": "sample"
    },
	{
      "name": "TableName",
      "type": "string",
      "value": "employee"
    },
    {
      "name": "Select",
      "type": "string",
      "value": "*"
    },
    {
      "name": "Where",
      "type": "string",
      "value": ""
    }	
  ],
  "outputs": [
    {
      "name": "result",
      "type": "any"
    }
  ]
}

```
## Inputs
| Input       | Description    |
|:------------|:---------------|
| ClusterIP   | The cluster to connect to |         
| Keyspace    | The keyspace the table to query resides in  |
| TableName   | The table to query |
| Select      | The element to return from the query (can be "*") |
| Where       | The selection criteria |

## Outputs
| Output       | Description    |
|:------------|:---------------|
| result   | The result of the query |  

