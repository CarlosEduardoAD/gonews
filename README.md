## GoNews

# GoNews é uma aplicação web desenvolvida em Go que permite aos usuários acessar e visualizar notícias de diversas fontes de forma eficiente e organizada. A aplicação utiliza Docker para facilitar a implantação e garantir um ambiente consistente.

-----

# Funcionalidades

- Agregação de notícias de múltiplas fontes.

- Interface web intuitiva para navegação e leitura de notícias.

- Desenvolvida em Go, garantindo desempenho e eficiência.

- Containerização com Docker para facilitar a implantação.

# Estrutura do Projeto

/
|-- internal/             # Pasta com o projeto e seus códigos
|-- main.go               # Arquivo principal da aplicação em Go
|-- go.mod                # Gerenciamento de dependências do Go
|-- go.sum                # Verificação de integridade das dependências
|-- Dockerfile.prod       # Configuração para construção do container Docker
|-- docker-compose.prod.yaml  # Orquestração de containers Docker
|-- .gitignore            # Especifica arquivos a serem ignorados pelo Git
|-- README.Docker.md      # Instruções relacionadas ao uso do Docker

# Tecnologias Utilizadas

- Go: Linguagem de programação principal utilizada no desenvolvimento da aplicação.

- Docker: Utilizado para containerização, facilitando a implantação e garantindo consistência entre ambientes.

Desenvolvido por Carlos Eduardo.
