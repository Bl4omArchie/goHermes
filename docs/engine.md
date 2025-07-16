# goHermes engine

The engine in goHermes the different components like the database, the log system, sources...
When you first start goHermes, you initiate an engine instance where you can set-up several parameters. 


## Create engine instance

There is the list of the configurable parameters for an engine instance :

| Parameters      |  Description      | Type | Values |
| ------------- | ------------- |  ------------- |  ------------- |
| databaseName | Filename of your database  | string | i.e : "database.db" for sqlite3 |
| wokers | Number of workers for the downloading worker pool | int | 50 in average is enought |
| sources | Every struct that implement the Source interface | Source | hermes.NewEprintSource() |

Learn more about the Source interface in [!source.md]
The engine is also setting-up the logging system which is still mandatory for now.
