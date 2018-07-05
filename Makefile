build:
	cd ./func/create && dep ensure
	cd ./func/show && dep ensure
	GOOS=linux go build -o bin/create ./func/create
	GOOS=linux go build -o bin/show ./func/show
deploy:
	@make build
	sls deploy

prod-deploy:
	@make build
	sls deploy --stage prod

remove:
	serverless remove -v
