### Build

```
go build -buildmode=plugin plugins/unencrypted_storage/unencrypted_storage.go
go build -buildmode=plugin plugins/local_bus/local_bus.go
go build -buildmode=plugin plugins/unix_bus/unix_bus.go
go build -buildmode=plugin plugins/my_service/my_service.go
go build -buildmode=plugin plugins/my_plugin/my_plugin.go
go run test.go
```

### Multi-process Bus
```
go run test.go
go run remote.go
```