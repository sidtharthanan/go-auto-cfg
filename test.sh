mkdir -p test_1
go run main.go test_samples/schema1.yml test_1 test_1
go test $(go list ./... | grep -v '/vendor/')
