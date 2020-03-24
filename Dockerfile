FROM golang:1.14-alpine

COPY namespace-reaper /namespace-reaper

ENTRYPOINT [ "/namespace-reaper" ]