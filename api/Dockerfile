# Etapa de build
FROM golang:1.23.1-alpine AS builder

# Instala as dependências necessárias, incluindo o WebP
RUN apk update && apk add --no-cache libwebp-dev gcc g++ make

# Defina o diretório de trabalho
WORKDIR /app

# Copie os arquivos de dependências go.mod e go.sum
COPY ./src/go.mod ./src/go.sum ./

# Baixe as dependências
RUN go mod download

# Copie o código fonte
COPY ./src .

# Compile o binário, nomeando-o como "api"
RUN go build -o api .

# Etapa de produção
FROM alpine:3.18

# Defina o diretório de trabalho
WORKDIR /app

# Copie apenas o binário gerado do estágio anterior
COPY --from=builder /app/api /app/api

# Expõe a porta 8080 para acesso ao serviço
EXPOSE 8080

# Comando padrão para executar o binário
CMD ["/app/api"]
