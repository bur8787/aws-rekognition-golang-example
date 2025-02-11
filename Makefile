.PHONY: all build_lambda deploy clean

LAMBDA_DIR=lambda
LAMBDA_BINARY=$(LAMBDA_DIR)/bootstrap
LAMBDA_ZIP=$(LAMBDA_DIR)/function.zip

all: clean build_lambda deploy

build_lambda:
	cd $(LAMBDA_DIR) && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	cd $(LAMBDA_DIR) && zip function.zip bootstrap

deploy:
	npx cdk bootstrap
	npx cdk deploy

clean:
	rm -f $(LAMBDA_BINARY) $(LAMBDA_ZIP)
