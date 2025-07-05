# testing with mockoon

## prerequisiti


1. install mockoon and mockoon/cli
2. install Golang
3. install Docker


Imagine we have to call a petstore service to get a list of available pets but the http APIs are not ready, we have only the openapi interface.


Using mockoon desktop I can import the yaml rapresentation (the example [petstore](https://editor.swagger.io/) is ok for test purpose) and export the json with the mocked service.


I’ll start with a test… in the arrange section we have to start the service mock using mockoon cli.


```go
mockedPetstore, addr, err := NewDockerMockoon().Start("petstore")
if err != nil {
	t.Fatalf("Failed to create DockerMockoon: %v", err)
}
defer mockedPetstore.TearDown()
```


