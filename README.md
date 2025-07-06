# testing with mockoon

## Prerequisites


1. install mockoon and mockoon/cli
2. install Golang
3. install Docker (or equivalent)

Imagine we need to call a “petstore” service to get a list of available pets, but the http APIs are not ready, we only  have the [openAPI interface](https://www.openapis.org/).

Using [Mockoon desktop](https://mockoon.com/download/) I can import the [YAML rapresentation](https://editor.swagger.io/) (the [petstore](https://editor.swagger.io/) example is fine for testing purpose) and export the JSON with the mocked service.

I’ll start with a test… in the “Arrange” section we need to start the service mock using [mockoon cli](https://mockoon.com/cli/).

```go
mockedPetstore, addr, err := NewDockerMockoon().Start("petstore")
if err != nil {
	t.Fatalf("Failed to create DockerMockoon: %v", err)
}
defer mockedPetstore.TearDown()
```


