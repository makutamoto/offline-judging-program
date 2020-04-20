FROM ubuntu:18.04

WORKDIR /usr/app

RUN apt update && apt install -y build-essential golang python3

RUN groupadd -r -g 400 code && useradd -r -u 400 -g 400 code

COPY ./judger ./
COPY ./languages ./languages

ENTRYPOINT ["./judger"]
