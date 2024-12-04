FROM public.ecr.aws/docker/library/golang:1.23 as build-image

WORKDIR /src

COPY . ./
RUN CGO_ENABLED=0 go build -o gin-lambda-handler
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build-image /src/gin-lambda-handler .
ENTRYPOINT ./gin-lambda-handler
