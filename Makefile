all: target

target:
	go install && cp $(GOPATH)/bin/netdisk ./

clean: 
	@rm ./netdisk
