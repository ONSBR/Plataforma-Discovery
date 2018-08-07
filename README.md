### Discovery

Discovery é o componente de plataforma responsável por localizar instâncias que devem ser reprocessadas caso a persistência de uma instância em execução seja realizada.

### Build

O build da aplicação é feito através de um arquivo Makefile, para buildar a aplicação execute o seguinte comando:

```sh
$ make
```

Após executar o make será criada uma pasta dist e o executável da aplicação discovery.

### Deploy

O processo de deploy do discovery na plataforma é feito através do installer, os componentes em Go são compilados e comitados dentro do installer então para atualizar a versão do discovery para atualizar a versão do discovery na plataforma utilize o seguinte comando:

```sh
$ mv dist/discovery ~/installed_plataforma/Plataforma-Installer/Dockerfiles
$ plataforma --upgrade discovery
```

### API

Retornar a lista de instâncias que devem ser reprocessadas
```http
GET /v1.0.0/discovery/entities?systemID=ec498841-59e5-47fd-8075-136d79155705&instanceID=f6f72706-2beb-4b7c-bf0f-bb571427f1bd HTTP/1.1
Host: localhost:8090
```

### Organização do código

1. actions
    * São as principais ações do serviço, por exemplo, recuperar as entidades que serão persistidas por uma instância, listar todos os cenários que estão disponíveis dentro da plataforma e etc;
2. api
    * É a declaração da API do discovery;
3. db
    * É o pacote que faz a conexão com o postgres para executar os filtros dos mapas;
4. helpers
    * Pacote de funções utilitárias
5. models
    * Define o modelo de domínio usado pelo discovery
6. vendor
    * É um pacote do Go onde ficam todas as bibliotecas de terceiros, os arquivos deste pacote jamais devem ser alterados diretamente;