# Makefile
all: setup hooks

# requires `nvm use --lts` or `nvm use node`
.PHONY: setup
setup: 
	npm install -g @commitlint/config-conventional @commitlint/cli  


.PHONY: hooks
hooks:
	@git config --local core.hooksPath .githooks/

.PHONY: test-webhook
test-webhook:
	@curl -sX POST http://localhost:8080/webhook \
		-H "Content-Type: application/json" \
		-d '{"eventVersion":"2.0","eventSource":"aws:s3","awsRegion":"us-east-1","eventTime":"2024-01-01T00:00:00.000Z","eventName":"s3:ObjectCreated:Put","userIdentity":{"principalId":"test"},"requestParameters":{"sourceIPAddress":"127.0.0.1"},"responseElements":{"x-amz-request-id":"test","x-amz-id-2":"test"},"s3":{"s3SchemaVersion":"1.0","configurationId":"test","bucket":{"name":"test-bucket","ownerIdentity":{"principalId":"test"},"arn":"arn:aws:s3:::test-bucket"},"object":{"key":"test-file.txt","size":123,"eTag":"test","versionId":"test"}}}' | jq .