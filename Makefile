# Broadly the test here is that _some_ JSON is generated, so the test fails
# if there are no converted JSON files for one of the tests by having the `rm`
# command exit with an error.
.PHONY: test
test: build
	./yamltojson testdata/test.yaml
	rm testdata/*.json
	./yamltojson testdata/folderofyaml
	rm testdata/folderofyaml/*.json
	rm yamltojson

.PHONY: build
build:
	go build .

.PHONY: clean
clean:
	rm yamltojson
