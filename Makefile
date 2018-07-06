build:
	cd ./func/create && dep ensure
	cd ./func/show && dep ensure
	cd ./func/index && dep ensure
	GOOS=linux go build -o bin/create ./func/create
	GOOS=linux go build -o bin/show ./func/show
	GOOS=linux go build -o bin/index ./func/index
deploy:
	@make build
	sls deploy

prod-deploy:
	@make build
	sls deploy --stage prod

remove:
	serverless remove -v
