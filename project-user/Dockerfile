FROM alpine
WORKDIR /Initial
COPY ./bin/project-user .
COPY config/config-docker.yaml .
RUN mkdir config && mv config-docker.yaml config/config.yaml
EXPOSE 9501 9001
ENTRYPOINT ["./project-user"]