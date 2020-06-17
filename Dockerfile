FROM golang:1.13-alpine

COPY bin/namespace-reaper /namespace-reaper

CMD [ "/namespace-reaper" ]