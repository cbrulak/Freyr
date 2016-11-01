default: freyr

.PHONY: freyr surtr

pwd := $(shell pwd)

freyr:
	GO15VENDOREXPERIMENT=1 docker run \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		-e "GOPATH=/go/src/github.com/serdmanczyk/freyr/vendor/:/go" \
		-v $(pwd):/go/src/github.com/serdmanczyk/freyr/ \
		-w /go/src/github.com/serdmanczyk/freyr/ \
		golang go build -ldflags "-s" -a -installsuffix cgo -o freyr

surtr:
	GO15VENDOREXPERIMENT=1 docker run \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		-e "GOPATH=/go/src/github.com/serdmanczyk/freyr/vendor/:/go" \
		-v $(pwd):/go/src/github.com/serdmanczyk/freyr/ \
		-w /go/src/github.com/serdmanczyk/freyr/cmd/surtr \
		golang go build -ldflags "-s" -a -installsuffix cgo -o surtr

webpack:
	cd nginx/static && PROD_ENV=1 webpack

rundev:
	docker-compose -f docker-compose.debug.yml -p dev build
	docker-compose -f docker-compose.debug.yml -p dev up --force-recreate

buildstatic:
	docker-compose build

runstatic:
	docker-compose up -d --build

downstatic:
	docker-compose down

integration:
	docker-compose -p integration down
	docker-compose -f docker-compose.integration.yml -p integration build
	docker-compose -f docker-compose.integration.yml -p integration up --force-recreate --abort-on-container-exit

acceptance:
	docker-compose -p acceptance down
	docker-compose -f docker-compose.acceptance.yml -p acceptance build
	docker-compose -f docker-compose.acceptance.yml -p acceptance up --force-recreate --abort-on-container-exit

buildgo:  freyr surtr
test:  integration acceptance
