populate:
	psql -U postgres -h localhost < input/cpu_usage.sql
	psql -U postgres -h localhost -d homework -c "\COPY cpu_usage FROM input/cpu_usage.csv CSV HEADER"

run:
	export $$(grep -v '^#' .env | xargs) && go run main.go