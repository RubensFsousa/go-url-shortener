# Nome do binário
BINARY_NAME := go-url-shortener

# Comandos
BUILD_CMD := go build -o $(BINARY_NAME).exe
RUN_CMD := ./$(BINARY_NAME).exe

# Alvo padrão: roda o CompileDaemon
.PHONY: run
run:
	CompileDaemon -build="$(BUILD_CMD)" -command="$(RUN_CMD)"

# Compilar manualmente sem usar o CompileDaemon
.PHONY: build
build:
	$(BUILD_CMD)

# Executar manualmente o binário
.PHONY: start
start:
	$(RUN_CMD)

# Limpar arquivos gerados (.exe e temporários)
.PHONY: clean
clean:
	rm -f $(BINARY_NAME).exe $(BINARY_NAME).exe~

# Gera a documentação Swagger
.PHONY: swag
swag:
	swag init --output ./docs --parseDependency --parseInternal