# quick start
## create app

    serverless create -u https://github.com/gaku3601/serverless-go/ -p your_app_name

## deploy

    make deploy

## deploy production

    make prod-deploy

# curl

    [show]
    curl -X <url>/<id>
    [index]
    curl -X GET <URL>?start=1\&end=10
