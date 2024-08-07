# md-blog

## dev
mock postgres-db
```
sudo docker run --name md-blog-db -e POSTGRES_PASSWORD=mdblogdb -e POSTGRES_USER=mdblog -e POSTGRES_DB=mdblog -p 5432:5432 -d postgres
sudo docker exec -it md-blog-db psql -U "mdblog" -d mdblog
```

run api
```
go run .
```
