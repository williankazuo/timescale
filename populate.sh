psql -U postgres < input/cpu_usage.sql
psql -U postgres -d homework -c "\COPY cpu_usage FROM input/cpu_usage.csv CSV HEADER"