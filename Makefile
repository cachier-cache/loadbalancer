start:
	go run main.go

# TODO: add this to the logic of routing the request to the server by creating more servers
# for availability
increase-scale:
	# this will increase the scale while running the current servers
	docker-compose up --scale core=5

# useful for when you want to check if this port is the server port or a system port
check-port:
	lsof -i:8086
