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
		-d '{"Records":[{"eventVersion":"2.2","eventSource":"aws:s3","awsRegion":"us-east-1","eventTime":"2025-06-19T04:51:34Z","eventName":"s3:ObjectCreated:Put","userIdentity":{"PrincipalId":"lfp-9dbf396e770b3361c426"},"requestParameters":{"sourceIPAddress":"10.244.3.134"},"responseElements":{"x-amz-request-id":"","x-amz-id-2":""},"s3":{"s3SchemaVersion":"1.0","configurationId":"webhook-global","bucket":{"name":"mybucket","ownerIdentity":{"PrincipalId":"lfp-9dbf396e770b3361c426"},"arn":"arn:aws:s3:::/mybucket/main.go"},"object":{"key":"main.go","size":2367,"eTag":"bfac34da409e8c51dc9346319e7f660b","versionId":null,"sequencer":"B"}},"glacierEventData":{"restoreEventData":{"lifecycleRestorationExpiryTime":"","lifecycleRestoreStorageClass":""}}}]}' | jq .