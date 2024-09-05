FROM alpine:latest AS files

WORKDIR /out

RUN echo "nogroup:x:65534:" > /out/group
RUN echo "nobody:*:65534:65534:nobody:/_nonexistent:/bin/false" > /out/passwd

FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=files /out/group /etc/group
COPY --from=files /out/passwd /etc/passwd
COPY overseerr-to-apprise /
USER 65534:65534
ENTRYPOINT [ "/overseerr-to-apprise" ]