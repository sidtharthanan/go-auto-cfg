go run main.go test_samples/schema1.yml test_1/config1.go
go test $(go list ./... | grep -v '/vendor/')
