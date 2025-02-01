# Usa a imagem oficial do Golang
FROM golang:1.19

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos do projeto para o contêiner
COPY . .

# Baixa as dependências
RUN go mod init file-storage && go mod tidy

# Compila o código
RUN go build -o main .

# Define o volume onde os arquivos serão armazenados
VOLUME [ "/app/storage" ]

# Define o comando padrão ao rodar o contêiner
CMD [ "./main" ]
