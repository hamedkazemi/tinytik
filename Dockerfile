# First stage: build the executable.
# It is important that these ARG's are defined after the FROM statement
FROM golang:1.20-alpine AS builder

# Add unprivileged user/group
RUN mkdir /user && \
	echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
	echo 'nobody:x:65534:' > /user/group

RUN echo http://repository.fit.cvut.cz/mirrors/alpine/v3.8/main > /etc/apk/repositories; \
    echo http://repository.fit.cvut.cz/mirrors/alpine/v3.8/community >> /etc/apk/repositories

# Create the user and group files that will be used in the running
RUN apk add --no-cache ca-certificates git

# container to run the process as an unprivileged user.
# Create a netrc file using the credentials specified using --build-arg
#RUN git config --global url."https://hamed.kazemi:XGmoxAMNZyCS_xnXVo1P@git.pinsvc.net".insteadof https://git.pinsvc.net
# Set the working directory outside $GOPATH to enable the support for modules.
# Working directory outside $GOPATH
WORKDIR /src
# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./go.work ./
# Import the code from the context.
COPY ./config.toml ./

COPY ./ ./
# Build the executable to `/app`. Mark the build as statically linked.
RUN GO111MODULE=on CGO_ENABLED=0 go build \
	-installsuffix 'static' \
	-o /app \
	.
    # Final stage: the running container.
FROM scratch AS final
# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled executable from the first stage.
COPY --from=builder /app /
# Perform any further action as an unprivileged user.
COPY config.toml /


# Open ports (if needed)
EXPOSE 8080
#EXPOSE 80
#EXPOSE 443

# Will run as unprivileged user/group
USER nobody:nobody

# Entry point for the built application
ENTRYPOINT ["/app"]