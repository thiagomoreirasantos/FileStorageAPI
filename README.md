# 📂 File Storage API
Esta é uma API REST em **Golang** para **upload, listagem, download e exclusão de arquivos**.  
Ela utiliza **Docker**, **Swagger** para documentação e **Gorilla Mux** como roteador HTTP.

---

## 🚀 **Funcionalidades**
✅ Upload de arquivos (`POST /upload`)  
✅ Listar arquivos armazenados (`GET /files`)  
✅ Baixar um arquivo (`GET /files/{filename}`)  
✅ Excluir um arquivo (`DELETE /files/{filename}`)  
✅ Documentação interativa com **Swagger** (`/swagger/index.html`)  

---

## 🛠 **Tecnologias Utilizadas**
- **Golang** (API)
- **Gorilla Mux** (Roteamento)
- **Swagger** (Documentação)
- **Docker** (Containerização)
- **Volumes Docker** (Persistência de arquivos)

---

## 📥 **Instalação e Execução**
### 1️⃣ **Clonar o repositório**
```bash
git clone https://github.com/seu-usuario/file-storage-api.git
cd file-storage-api
```

### 2️⃣ **Instalar as dependências**
```bash
go mod tidy
```

### 3️⃣ **Gerar a documentação Swagger**
```bash
swag init
```

### 4️⃣ **Rodar a API**
```bash
go run main.go
```

###  **A API estará disponível em:**
```bash
http://localhost:8080
http://localhost:8080/swagger/index.html
```

### 🐳 **Rodando com Docker**

### **Se preferir rodar a API dentro de um contêiner Docker:**
### 1️⃣ **Construir a imagem Docker**
```bash
docker build -t file-storage-api .
```

### 2️⃣ **Rodar o contêiner**
```bash
docker run -p 8080:8080 -v storage_data:/app/storage file-storage-api
```

### ***Agora, a API pode ser acessada em http://localhost:8080.***