# osrs_market_analysis

A tool used to track OSRS market data

Docker:
$ docker network create osrsflip
$ docker build -t angular-app .
$ cd srv
$ docker build -t go-server .

    $ sudo docker run --name osrs_db -p 5432:5432 -e POSTGRES_PASSWORD=1234 --network osrsflip -d postgres
    $ sudo docker run -d --name go-server -p 8080:8080 --network osrsflip go-server
    $ sudo docker run -d --name angular-app -p 4200:4200 --network osrsflip angular-app

Migrating to AWS:

1. Tag all docker images:
   $ docker tag angular-app pipthedev/angular-app:latest
   $ docker tag go-server pipthedev/go-server:latest
   $ docker tag postgres pipthedev/postgres:latest

2. Push all docker images:
   $ docker push pipthedev/angular-app:latest
   $ docker push pipthedev/go-server:latest
   $ docker push pipthedev/postgres:latest

3. Sign into ec2 instance with ssh

4. Install docker: https://docs.docker.com/engine/install/ubuntu/

5. Pull docker images:
   $ docker pull pipthedev/angular-app:latest
   $ docker pull pipthedev/go-server:latest
   $ docker pull pipthedev/postgres:latest

6. Run all docker images:
   $ sudo docker run --name osrs_db -p 5432:5432 -e POSTGRES_PASSWORD=1234 --network osrsflip -d pipthedev/postgres
   $ sudo docker run -d --name go-server -p 8080:8080 --network osrsflip pipthedev/go-server
   $ sudo docker run -d --name angular-app -p 80:4200 --network osrsflip angular-app:latest
