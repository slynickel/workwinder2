
# nucular_gio should be used for mac os. It's supposed use this by default
# but for some reason doesnt?
run:
	go run -tags nucular_gio cmd/nucular/main.go cmd/nucular/settings.go

run-def:
	go run  cmd/nucular/main.go cmd/nucular/settings.go

run-def2:
	go run cmd/nucular/main.go cmd/nucular/settings.go

chart:
	go run cmd/chart/main.go

j:
	jq . test.json