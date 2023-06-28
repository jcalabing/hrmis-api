build:
	go build -o ./cmd/hrmis-api ./cmd

run:
	CompileDaemon -command="./hrmis-api" 

clean:
	rm -f ./cmd/hrmis-api