build:
	go build 
push: 
	git add . 
	git commit -m "additions"
	git push
test:
	./tests/endpoints.sh
