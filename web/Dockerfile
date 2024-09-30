# Usa a imagem base do Node.js na versão 16 e nomeia esta fase como 'build'
FROM node:16 AS build

# Define variáveis de ambiente que podem ser passadas no momento da construção da imagem
ARG APP_ENVINRONMENT
ARG APP_API_KEY_TOKEN

# Define as variáveis de ambiente com os valores recebidos dos argumentos
ENV APP_ENVINRONMENT=$APP_ENVINRONMENT
ENV APP_API_KEY_TOKEN=$APP_API_KEY_TOKEN

# Define o diretório de trabalho para o aplicativo
WORKDIR /app

# Copia os arquivos de configuração do npm para o diretório de trabalho
COPY package*.json ./

# Executa a instalação das dependências do npm
RUN npm install

# Copia todos os arquivos do projeto para o diretório de trabalho
COPY . .

# Executa o script de build do aplicativo
RUN npm run build

# Inicia uma nova fase a partir da imagem base do Nginx com Alpine
FROM nginx:alpine

# Copia os arquivos gerados na fase de build para o diretório do Nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Define permissões para os arquivos copiados para o Nginx
RUN chmod -R 755 /usr/share/nginx/html

# Copia o arquivo de configuração do Nginx para o diretório apropriado
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copia o script de configuração de ambiente para o diretório de entrada do Docker
COPY env.sh /docker-entrypoint.d/env.sh

# Torna o script copiado executável
RUN chmod +x /docker-entrypoint.d/env.sh

# Expõe a porta 80 para acesso externo
EXPOSE 80

# Comando para iniciar o Nginx em modo não daemon
CMD ["nginx", "-g", "daemon off;"]