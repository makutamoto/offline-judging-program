FROM ubuntu:18.04

WORKDIR /usr/app

RUN apt update && apt install -y build-essential golang python3

COPY ./judger ./
COPY ./languages ./languages

CMD ["./judger"]
