package cassandra

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/gocql/gocql"
)

// THIS IS ADDED
// log is the default package logger which we'll use to log
var log = logger.GetLogger("activity-jl-cassandra")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	// Get the activity data from the context
	clusterIP := context.GetInput("clusterIP").(string)
	keySpace := context.GetInput("keySpace").(string)
	tableName := context.GetInput("tableName").(string)
	selectElements := context.GetInput("select").(string)
	whereClause := context.GetInput("where").(string)

	// Use the log object to log the greeting
	//log.Debugf("The Flogo engine says [%s] to [%s] with table [%s]", clusterIP, keySpace, tableName)
	log.Debugf("Flogo is about to select [%s] from table [%s].[%s] where [%s] on cluster [%s]", selectElements, keySpace, tableName, whereClause, clusterIP)

	// Provide the cassandra cluster instance here.
	cluster := gocql.NewCluster(clusterIP)

	// gocql requires the keyspace to be provided before the session is created.
	// In future there might be provisions to do this later.
	cluster.Keyspace = keySpace

	// cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		log.Debugf("Could not connect to cassandra cluster: ", err)
		context.SetOutput("result", "ERROR_CONNECT")
		context.SetOutput("rowCount", 0)
		return true, err
	}
	log.Debugf("Session Created Sucessfully")

	log.Debugf("Cluster: %v", clusterIP)
	log.Debugf("Keyspace: %v", keySpace)
	log.Debugf("Session Timeout: %v", cluster.Timeout)

	log.Debugf("Next Step is Query Execution")
	queryString := "SELECT " + selectElements + " FROM " + tableName
	if whereClause != "" {
		queryString += " where " + whereClause
	}
	log.Debugf("Query string: [%s]", queryString)

	iter := session.Query(queryString).Iter()
	log.Debugf("number of columns: %v", len(iter.Columns()))

	context.SetOutput("rowCount", iter.NumRows())

	if iter.NumRows() == 0 {
		context.SetOutput("result", "NO_DATA")
		return true, nil
	}
	var result []map[string]interface{}
	for i := 0; i < iter.NumRows(); i++ {
		row := make(map[string]interface{})
		if !iter.MapScan(row) {
			log.Debug("Error Select")
			iter.Close()
		}
		result = append(result, row)
		for _, column := range iter.Columns() {
			log.Debugf("Record [%v], Field [%v], value [%v]", i, column.Name, row[column.Name])
		}
	}

	context.SetOutput("result", result)

	// Signal to the Flogo engine that the activity is completed
	return true, nil
}
