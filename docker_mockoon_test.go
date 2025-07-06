package dockerized_mockoon

import (
	"context"
	"testing"

	"github.com/antihax/optional"

	"github.com/diegoavanzini/dockerized_mockoon/clients/petstore"
)

func TestPetStoreCall_WhenMockoonSetted_ShouldReturnAResponse(t *testing.T) {
	// ARRANGE
	serviceMockFactory, err := NewDockerMockoon()
	if err != nil {
		t.Fatalf("Failed to create DockerMockoon: %v", err)
	}
	serviceMockFactory.AddService(Service{
		Name:        "petstore",
		ExposedPort: 3005,
	})
	mockedPetstore, addr, err := serviceMockFactory.Start("petstore")
	if err != nil {
		t.Fatalf("Failed to create DockerMockoon: %v", err)
	}
	defer mockedPetstore.TearDown()
	petstoreConfiguration := petstore.NewConfiguration()
	petstoreConfiguration.BasePath = *addr
	client := petstore.NewAPIClient(petstoreConfiguration)

	// ACT
	pet, _, err := client.PetApi.FindPetsByStatus(context.Background(), &petstore.FindPetsByStatusOpts{
		Status: optional.NewString("available"),
	})
	if err != nil {
		t.Fatalf("Failed to FindPetsByStatus: %v", err)
	} // ASSERT
	if len(pet) == 0 {
		t.Fatal("Expected at least one pet in response")
	}
	if pet[0].Name == "" {
		t.Fatal("Expected pet to have a name")
	}
}

func TestAddService_WhenServiceAlreadySetted_ShouldReturnError(t *testing.T) {
	// ARRANGE
	serviceMockFactory, err := NewDockerMockoon()
	if err != nil {
		t.Fatalf("Failed to create DockerMockoon: %v", err)
	}
	err = serviceMockFactory.AddService(Service{
		Name:        "petstore",
		ExposedPort: 3005,
	})
	if err != nil {
		t.Fatalf("Failed to add service: %v", err)
	}
	err = serviceMockFactory.AddService(Service{
		Name:        "petstore",
		ExposedPort: 3005,
	})
	if err == nil {
		t.Fatal("Expected error when adding service with the same name")
	}
	if err.Error() != "Service with name 'petstore' already exists" {
		t.Fatalf("Expected error 'Service with name 'petstore' already exists', got '%s'", err.Error())
	}
}
