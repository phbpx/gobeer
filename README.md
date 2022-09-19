# gobeer
[![Test and coverage](https://github.com/phbpx/gobeer/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/phbpx/gobeer/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/phbpx/gobeer/branch/main/graph/badge.svg?token=KV1V27Z55R)](https://codecov.io/gh/phbpx/gobeer)

Projeto de demonstração utilizado na talk sobre estrutura de projetos Go, onde utilizamos princípios de design baseados no __DDD__ (Domain Driven Design) e __Arquitetura Hexagonal__ para o desenvolvimento de um pequeno backend.

## Conteúdo

- [Motivação](#motivação)
- [Projeto](#projeto)
- [Tecnologias](#tecnologias)
- [Organização](#organização)
- [Desenvolvimento](#desenvolvimento)
    - [Clonando o repositório](#clonando-o-repositório)
    - [Makefile](#comandos-makefile)
    - [Setup](#setup)
    - [Executando testes](#executando-testes)
    - [Executando o ambiente local](#executando-o-ambiente-local)
    - [Postman](#postman)
    - [Banco de dados](#banco-de-dados)

## Motivação

Um dos principais motivos para a não adoção da linguagem Go em novos projetos é a dúvida relacionada a como devemos estruturar os projetos. A linguagem não nos dá nenhuma dica a respeito disso, é como se fosse tudo por nossa conta. O que é ótimo na verdade, porque podemos nos divertir com a estrutura e podemos deixá-la do nosso jeito. Mas ao mesmo tempo isso pode ser um problema, especialmente para quem vem de linguagens Orientadas a Objetos como  Java ou C#, e esperam ver classes e objetos. Mas quando se deparam com os idiomas da linguagem não sabem bem por onde começar, ou começam escrevendo código Go com os idiomas que estão acostumados.

## Projeto

Como demontração vamos fazer um serviço para __review de cervejas__, com os seguintes requisitos:
 - O usuário pode criar uma cerveja.
 - O usuário pode listar todas as cervejas.
 - O usuário pode fazer um review de uma cerveja.
 - O usuário pode listar os reviews de uma determinada cerveja.
 - Os dados devem ser armazenados em banco de dados.

Tradução dos requisitos para uma linguagem oblíqua:

- __Context__: Beer tasting (degustação de cervejas)
- __Language__: Beer, Review
- __Entity__: Beer, Review
- __Service__: Beer adder, Beer lister, Review adder, Review lister
- __Events__: Beer added, Review added
- __Repository__: Beer repository, Review Repository

## Tecnologias

Além de __Go__ como linguagem, utilizamos também:
- __PostgreSQL__: Banco de dados.
- __docker/docker-compose__: Para nos ajudar com o ambiente de desenvolvimento.

## Organização

O projeto utiliza uma filosofia de design baseada em __DDD__ (Domain Driven Design) e __Arquitetura Hexagonal__. Tento como principais caracteristicas:

- Consistente
- Fácil de entender e de navegar
- Fácil de mudar (loosely-coupled/acoplamento fraco)
- Fácil de testar
- O mais simples possível, mas não simples demais
- O design reflete exatamente como o software funciona
- A estrutura reflete exatamente o design

## Desenvolvimento

Para desenvolver para este projeto você deve ter em seu computador:
- golang ([veja como instalar](https://go.dev/doc/install)).
- docker ([veja como instalar](https://docs.docker.com/get-docker/)).
- docker-compose ([veja como instalar](https://docs.docker.com/compose/install/)).

#### Rodando o ambiente local

Se você nunca desenvolveu neste repositório antes:

#### Clonando o repositório

`Com HTTPS`
```sh
$ git clone https://github.com/phbpx/gobeer.git
```
`Com SSH`
```sh
$ git clone git@github.com:phbpx/gobeer.git
```

#### Comandos makefile

Para ajudar com o fluxo de trabalho o projeto possui um Makefile com alguns comandos úteis, para conhecê-lo, rode o comando:

```sh
$ make help
```

Output:
```
Usage:
  make <target>

Targets:
  help                 Show help
  setup                Install go tools
  test                 Execute tests
  cover                Execute tests with coverage visualization
  lint                 Execute static check
  dev                  Run local environment
  stop                 Stop local environment
```

#### Setup

Alguns comandos make utilizam ferramentas de linha de comando do Go, e para instalar é necessário rodar o comando:

```sh
make setup
```

#### Executando testes

O projeto possui uma suite de testes de unidade (em use cases) e testes de integração (nas portas/adapters). Para executar, rode o comando:

```sh
$ make test
```

Para executar os testes e vizualisar o relatório de cobertura de código, rode o comando:

```sh
make cover
```

#### Executando o ambiente local

O projeto possui um `docker-compose` com o necessário para rodar localmente a api:

```sh
$ make dev
```

E para desligar o ambiente:

```sh
make stop
```

Endereços:
- gobeer-api: `http://localhost:3000`
  - Adding beer: `POST http://localhost:3000/beers`
  - Listing beers: `GET http://localhost:3000/beers`
  - Adding beer review: `POST http://localhost:3000/beers/:beer_id/reviews`
  - Listing beer reviews: `GET http://localhost:3000/beers/:beer_id/reviews`
  - Helthcheck: `GET http://localhost:3000/debug/health`

#### Postman

Para facilitar a utilização da api, o repositótio possui uma collection postman ([link para o arquivo](https://raw.githubusercontent.com/phbpx/gobeer/main/gobeer-api.postman_collection.json)).

#### Banco de dados

Para acessar o banco de dados com o adminer:

1.  Acesse: [`http://localhost:8080`](http://localhost:8080)
2.  No campo `System`, selecione a opção: `PostgreSQL`
3.  No campo `Server`, preencha com: `db`
4.  No campo `Username`, preencha com: `postgres`
5.  No campo `Password`, preencha com: `postgres`
6.  No campo `Database`, preencha com: `testdb`

Caso queira acessar o banco com um client de sua preferência:

- Host: `localhost:5432`
- Username: `postgres`
- Password: `postgres`
- Database: `testdb`
