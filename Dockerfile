FROM alpine:latest
ADD ./task .
EXPOSE 8091
ENV ENV=production
CMD [ "./task" ]