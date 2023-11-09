# tinyKVStore
tinyKVStore is a simple KV storage system that supports storing KV pairs of string types. It maintains a skiplist in memory to speed up key retrieval and stores values on disk.  

## Usage
Run the following command to start a tingKVStore database:  
```shell
go run main.go
```
Run the following command to start a tinyKVStore database client:  
```shell
go run db-cli.go
```
Then you can manipulate the database through the database client.  

## Currently supported commands
Using `put [key] [value]` to store a KV pair to the database.  
Using `get [key]` to get the value of the key.
Using `delete [key]` to delete a KV pair.

## TODO
tinyKVStore is under development and future plans are as follows:  
- Implement crash consistency