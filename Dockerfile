FROM golang:1.14-alpine

COPY bin/namespace-reaper /namespace-reaper

CMD [ "/namespace-reaper" ]