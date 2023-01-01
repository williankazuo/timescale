# Timescale Challenge

## Dependencies
- Docker Engine >= 19.03.0
- docker-compose >= 1.29.0

## How to run
### docker-compose
This will run using the input `query_params.csv`
```
docker-compose up
```

To run with different parameters modify the `docker-compose.yml`
- To modify the `workers` param just change the number
- To modify the file which will run the benchmarking, put the file in `input` folder and change the path of `filepath`
```
command: "./bench -workers=4 -filepath=/input/fileininput.csv"
```


### Makefile
- Create a `.env` file following the example from `.env.sample`

Run the db via docker-compose
```
docker-compose up database
```
And then run the program
```
make run
```
To run with different parameters modify the `Makefile`
```
go run main.go -workers=2 -filepath="./input/query_params.csv"
```