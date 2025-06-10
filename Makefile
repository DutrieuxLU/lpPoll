push: 
	git add . 
	git commit -m "additions"
	git push
test:
	cd tests
	./endpoints.sh
