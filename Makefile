docker: compile
	docker build -t dasmith/appreg .

compile:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o appreg .
