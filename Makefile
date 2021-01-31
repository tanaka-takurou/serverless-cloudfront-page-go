root	:=		$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: clean build deploy

clean:
	rm -rfv bin

build:
	mkdir -p bin
	scripts/create_template.sh
	cp -r templates bin/templates
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main

deploy:
	sam package --output-template-file "${root}"/packaged.yml --s3-bucket "${bucket}"
	sam deploy --stack-name "${stack}" --capabilities CAPABILITY_IAM --template-file "${root}/packaged.yml"
