# Imagem base com Go
FROM golang:1.21-alpine

# Diretório de trabalho dentro do container
WORKDIR /app

# Copia o código-fonte para o container
COPY . .

# Instala as dependências do projeto
RUN go mod download

# Compila o código Go
RUN go build -o main .

# Expõe a porta que a API irá utilizar
EXPOSE 8080

# Comando para iniciar a API
CMD ["./main"]