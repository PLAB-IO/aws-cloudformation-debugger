ifndef profile
$(error You must specify profile parameter)
endif

ifndef bucket
$(error You must specify bucket parameter)
endif

deploy-demo:
	aws --profile $(profile) \
	  cloudformation package \
		--template-file demo/main.yml \
		--s3-bucket $(bucket) \
		--output-template-file .cf-main-output.yml

	aws --profile $(profile) \
	  cloudformation deploy \
		--template-file .cf-main-output.yml \
		--capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
		--stack-name cfdbg-demo