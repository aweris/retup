FROM alpine:3.12

WORKDIR /app

COPY retup /bin/retup

LABEL vendor="aweris" \
      name="retup" \
      description="A tool that allows creating a distribution directory more flexible way for the mono repos" \
      maintainer="Ali Akca <ali@akca.io>"

ENTRYPOINT ["/bin/retup"]