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

Para desenvolver para este projeto você deve ter o golang instalado ([veja como instalar o go](https://go.dev/doc/install)).

#### Rodando o ambiente local

Se você nunca desenvolveu neste repositório antes:

- Clone o repositório:

Com HTTPS
```sh
$ git clone https://github.com/phbpx/gobeer.git
```
Com SSH
```sh
$ git@github.com:phbpx/gobeer.git
```

