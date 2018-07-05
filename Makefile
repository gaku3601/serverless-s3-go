build:
	cd ./func/create && dep ensure
	GOOS=linux go build -o bin/create ./func/create
deploy:
	@make build
	sls deploy

prod-deploy:
	@make build
	sls deploy --stage prod

remove:
	serverless remove -v
