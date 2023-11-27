
# This docker file is for development/local environment
FROM debian:stable-slim

WORKDIR /

COPY SuperEsbAdminWeb /
RUN mkdir /conf
COPY conf /conf

RUN mkdir /views
COPY views /views

RUN mkdir /frontend
COPY frontend /frontend

RUN mkdir /static
COPY static /static

RUN chmod +x /SuperEsbAdminWeb


ENTRYPOINT ["/SuperEsbAdminWeb"]
