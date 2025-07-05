openapi-generator-cli-version := 4.1.1
project-name := dockerized_mockoon

client: openapi-generator
	java -jar ./clients/openapi-generator-cli.jar generate -g go -i c:/GoCode/src/github.com/diegoavanzini/$(project-name)/data/petstore/petstore.yaml --additional-properties=enumClassPrefix=true --additional-properties=packageName=petstore -o ./clients/petstore
	rm ./clients/petstore/go.mod
	rm ./clients/petstore/go.sum

delete_client:
	rm -rf ./clients/petstore

test:
	go test ./... 

openapi-generator: clean
	curl https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/$(openapi-generator-cli-version)/openapi-generator-cli-$(openapi-generator-cli-version).jar --output ./clients/openapi-generator-cli.jar

clean:
	rm -rf ./clients/petstore
	rm -f ./clients/openapi-generator-cli.jar

build:
	go build