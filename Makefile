run:
	export $$(grep -v '^#' .env | xargs) && go run main.go -workers=2 -filepath="./input/query_params.csv"