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
      "name": "clusterIP",
      "type": "string",
      "value": "localhost"
    },
	{
      "name": "keySpace",
      "type": "string",
      "value": "sample"
    },
	{
      "name": "tableName",
      "type": "string",
      "value": "employee"
    },
    {
      "name": "select",
      "type": "string",
      "value": "*"
    },
    {
      "name": "where",
      "type": "string",
      "value": ""
    }	
  ],
  "outputs": [
    {
      "name": "result",
      "type": "any"
    },
    {
      "name": "rowCount",
      "type": "int"
    }
  ]
}

```
## Inputs
| Input       | Description    |
|:------------|:---------------|
| clusterIP   | The cluster to connect to |         
| keySpace    | The keyspace the table to query resides in  |
| tableName   | The table to query |
| select      | The element to return from the query (can be "*") |
| where       | The selection criteria |

## Outputs
| Output       | Description    |
|:------------|:---------------|
| result   | The result of the query |  
| rowCount   | The result of the query |  

