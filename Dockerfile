# Build container
FROM msales/go-builder:1.14-base-1.0.0 as builder

# Set token
ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

COPY ./ .
