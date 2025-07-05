package dockerized_mockoon

import (
	"context"
	"testing"

	"github.com/antihax/optional"

	"github.com/diegoavanzini/dockerized_mockoon/clients/petstore"
)

func TestPetStoreCall_WhenMockoonSetted_ShouldReturnAResponse(t *testing.T) {
	// ARRANGE
	mockedPetstore, addr, err := NewDockerMockoon().Start("petstore")
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
