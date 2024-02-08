[![ci](https://github.com/IgorCastilhos/BankApplication/actions/workflows/ci.yml/badge.svg)](https://github.com/IgorCastilhos/BankApplication/actions/workflows/ci.yml)
# Bank Application


## Setup do ambiente de desenvolvimento

### Instale essas ferramentas

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (Linux/Mac)

    ```bash
    brew install golang-migrate
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

### Setup da infraestrutura

- Inicie o contêiner postgres:

    ```bash
    make postgres
    ```

- Crie o banco de dados 'bank':

    ```bash
    make createdb
    ```

- Execute a migração do banco de dados em todas as versões:

    ```bash
    make migrateup
    ```

- Execute a migração do banco de dados para uma versão acima:

    ```bash
    make migrateup1
    ```

- Execute o migratedown para voltar todas as versões:

    ```bash
    make migratedown
    ```

- Execute o migratedown1 para uma versão anterior:

    ```bash
    make migratedown1
    ```


### Como gerar código

- Gere o CRUD SQL com sqlc:

    ```bash
    make sqlc
    ```

- Gere a simulação de banco de dados(mock) com gomock:

    ```bash
    make mock
    ```

### Como executar

- Execute o servidor:

    ```bash
    make server
    ```

- Execute os testes:

    ```bash
    make test
    ```
    
## Design do banco de dados

![image](https://github.com/IgorCastilhos/BankApplication/assets/101683017/761b2b23-b0b7-499d-b66f-49c271e74eb9)

---
## Deadlock
Script usado no TablePlus para identificar um dos Deadlock's
* `SELECT
  a.application_name,
  l.relation::regclass,
  l.transactionid,
  l.mode,
  l.locktype,
  l.GRANTED,
  a.username,
  a.query,
  a.pid
  FROM pg_stat_activity a
  JOIN pg_locks l ON l.pid = a.pid
  WHERE a.application_name = 'psql'
  ORDER BY a.pid;`

Foi possível identificar que a transação 1 estava tentando adquirir um ShareLock pelo transactionID '1167', porém, ela ainda não havia ganho um, pois a transação 2 já possuía um ExclusiveLock, no mesmo transactionID.
![img_1.png](imgs/img_1.png)
Portanto, **a transação 1 deve aguardar a transação 2 terminar antes de continuar**.
![img_2.png](imgs/img_2.png)

Ao tentar atualizar o saldo da conta 1, acontece o deadlock. A transação 2 também precisa aguardar a transação 1 terminar, para receber o resultado da consulta (query).
Resumo: O Deadlock ocorre, pois, ambas as transações concorrentes, **precisam aguardar a outra terminar.**
Para resolver, precisei mover a ordem do UPDATE da transação 2. Agora ambas as transações 1 e 2 **sempre irão atualizar a account1 antes da account2**.
A melhor maneira de prevenir deadlocks é fazer com que a aplicação sempre **adquira locks em uma ordem consistente!**

---
## Validação
### Oneof
Com esse pequeno código `oneof`, o Gin valida o input, **obrigando que seja BRL.**
![img_3.png](imgs/img_3.png)

###### Ao passar uma string vazia ou um câmbio monetário inválido, ele retorna um erro.
![img_4.png](imgs/img_4.png)

### Min
O `min=1` valida que o ID deve ser no mínimo um **inteiro positivo 1**

![img_5.png](imgs/img_5.png)

###### Exemplo de erro no Postman:

![img_6.png](imgs/img_6.png)

---

## Paginação
* No arquivo api/accounts.go, eu limitei a quantidade de contas que podem ser retornadas da seguinte maneira:

![img.png](imgs/img7.png)
* Caso o id enviado seja 0, o **validador** do Gin irá retornar um `failed on the 'required' tag` 

---

## Viper - Carregando configurações a partir de um arquivo e variáveis de ambiente
* Por que um arquivo?
  * Facilita a especificação da configuração padrão para desenvolvimento local e de teste.
* Por que variáveis de ambiente? (Env)
  * Facilitam a sobrescrita das configurações padrões quando é feito o deploy da aplicação.
* Viper pode encontrar, carregar e fazer unmarshal de arquivos de configuração
  * JSON, TOML, YAML, ENV, INI
* Também pode ler valores de variáveis de ambiente ou flags de linha de comando
  * Permite sobrescrever valores padrões
* Se preferir armazenar os arquivos de configuração em um sistema remoto como **etcd ou consul**, o viper pode ler a partir deles também
* Viper observa qualquer modificação nas configurações, podendo notificá-las também
  * Ele lê as mudanças feitas e as salva

---
## Mock - Por que mockar uma base de dados?
* Ajuda a escrever testes independentes facilmente, porque cada teste usará seu próprio banco de dados para armazenar dados, evitando conflitos entre eles.
  * Se tu usar um banco de dados real, todos os testes irão ler e escrever dados no mesmo lugar, aumentando a chance de conflitos, especialmente em grandes projetos.
* Os testes executarão rapidamente, visto que eles não precisarão se comunicar com o banco de dados real e esperar todas as queries serem executadas. Todas as ações são executadas em memória e dentro do mesmo processo.
* 100% de cobertura
  * Com um banco de dados mock, pode-se facilmente configurar e testar **alguns casos extremos**, como _unexpected error_ ou um _connection lost_, que seriam impossíveis de testar usando um banco de dados real.

### Como Mockar?
1. Implemente um banco de dados falso ou "fake", que armazena dados em memória.
   * Um porém: criar estruturas e funções pode demandar muito tempo, por isso utilizei GoMock nesse projeto.
2. GoMock gera e constrói stubs, que retornam dados hardcoded para cada cenário que queremos testar.

---
## Hash Password com o pacote Bcrypt em Go
![image](https://github.com/IgorCastilhos/BankApplication/assets/101683017/89e6d19d-5b32-4325-9091-f3b563b15c94)
**Bcrypt** é uma função de hashing de senha projetada especificamente para proteção de senhas armazenadas. Desenvolvida por Niels Provos e David Mazières para o sistema operacional OpenBSD, ela é baseada no ciframento Blowfish e é amplamente utilizada devido à sua segurança e eficiência. O Bcrypt é projetado para ser lento de propósito, o que é uma característica desejável para funções de hashing de senha, tornando ataques de força bruta muito menos eficazes.

### Características principais do Bcrypt:

- **Custo Adaptável:** Bcrypt permite ajustar o custo (ou fator de trabalho), que é basicamente quantas vezes o hashing é processado. Isso torna o algoritmo escalável; à medida que o hardware fica mais rápido, o custo pode ser aumentado para tornar o hash mais lento.
- **Sal Integrado:** Para proteger contra ataques de dicionário e uso de tabelas rainbow, o Bcrypt automaticamente gera e incorpora um sal (um valor aleatório) às senhas antes de aplicar o hash. Isso garante que hashes de duas instâncias da mesma senha sejam diferentes.
- **Resistente a Ataques de Força Bruta:** Devido à sua natureza inerentemente lenta e ao custo adaptável, o Bcrypt é resistente a ataques de força bruta.

### Como o Bcrypt é usado:

1. **Hashing de Senhas:** Ao armazenar uma senha, você não armazena a senha em si, mas um hash bcrypt dela. Quando um usuário fornece uma senha, você aplica o mesmo processo de hash à senha fornecida e compara com o hash armazenado.
2. **Verificação de Senhas:** Para verificar uma senha, o Bcrypt extrai o sal do hash armazenado (já que o sal é armazenado junto com o hash), aplica o hash à senha fornecida usando esse sal e verifica se o resultado corresponde ao hash armazenado.

### Exemplo em Go:

Para usar o Bcrypt em Go, você pode utilizar o pacote `golang.org/x/crypto/bcrypt`. Aqui está um exemplo simples de como hashar e verificar uma senha com Bcrypt em Go:

```go
package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	password := "mySecretPassword"
	hashedPassword, err := hashPassword(password)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hashed password:", hashedPassword)

	match := checkPasswordHash("mySecretPassword", hashedPassword)
	fmt.Println("Password match:", match) // Deve ser true

	match = checkPasswordHash("wrongPassword", hashedPassword)
	fmt.Println("Password match:", match) // Deve ser false
}
```

Neste exemplo, `bcrypt.GenerateFromPassword` é usado para criar um hash da senha com um sal gerado automaticamente e um custo padrão. Para verificar a senha, `bcrypt.CompareHashAndPassword` compara o hash da senha armazenada com o hash da senha fornecida pelo usuário.


