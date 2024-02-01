[![ci-test](https://github.com/IgorCastilhos/BankApplication/actions/workflows/ci.yml/badge.svg)](https://github.com/IgorCastilhos/BankApplication/actions/workflows/ci.yml)
### Bank Application

* Script usado no TablePlus para identificar um dos Deadlock's
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

* Foi possível identificar que a transação 1 estava tentando adquirir um ShareLock pelo transactionID '1167', porém, ela ainda não havia ganho um, pois a transação 2 já possuía um ExclusiveLock, no mesmo transactionID.
![img_1.png](img_1.png)
* Portanto, a transação 1 deve aguardar a transação 2 terminar antes de continuar.
![img_2.png](img_2.png)
* Ao tentar atualizar o saldo da conta 1, acontece o deadlock. A transação 2 também precisa aguardar a transação 1 terminar, para receber o resultado da consulta (query).
* Resumo: O Deadlock ocorre, pois, ambas as transações concorrentes, **precisam aguardar a outra terminar.**
* Para resolver, precisei mover a ordem do UPDATE da transação 2. Agora ambas as transações 1 e 2 **sempre irão atualizar a account1 antes da account2**.
* A melhor maneira de prevenir deadlocks é fazer com que a aplicação sempre **adquira locks em uma ordem consistente!**
