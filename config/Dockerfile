# syntax=docker/dockerfile:1
FROM public.ecr.aws/docker/library/golang:1.20.6-alpine3.18 as builder

WORKDIR /app

# Copy go mod and sum files
COPY ./ ./
RUN go mod download

#Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

FROM public.ecr.aws/lambda/go:1 as dev
COPY --from=builder /app/main ./main
CMD [ "main" ]

FROM public.ecr.aws/docker/library/alpine:3.18.2 as prod
COPY --from=builder /app/main ./main
CMD [ "main" ]
