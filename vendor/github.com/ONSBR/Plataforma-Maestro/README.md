# Plataforma-Maestro
Maestro é um componente de plataforma para orquestrar o processo de persistência e reprocessamento da plataforma

Fluxo Inicial de persistência:

![fluxo-01](https://www.planttext.com/plantuml/img/TP4zRiCm38LtdOAZSmLdQ8QcG8OMGTAn2slHAG5PSaHPIC_JeGV9nKh9HMkaAA51uC_x7gMD98nf6fmnYPCZU73J9S3ESyVeO6V99-wvGq1ueev4sA8bq7EWCOQImK6R0hnuU4II50CSAMRkI5F6L80nxK6dNmcskJRh_9wYi2JoIbeR3UwXUQO1uetmIBxOg53KKiRhv_KZtAqWlP67HdXO6T1eJnD6Yq0_Z77103yD23qxB1KIhIcNd10qNlMgHhkda-ugm5wDbp61yrJGMPq9nOKxgpKVu9wb2vdUrzL3MM9xASobn8WH5vFnhtT52xgiBdyz_lSGCgDkh1U9uooXfiA0xB_xlEM-tinyBI4fZCPiDA6V_mK0 "Fluxo Inicial")


Fluxo-02: Para verificar se existe um reprocessamento pendente:
![fluxo-02](https://www.planttext.com/plantuml/img/VP0nRiCm34LtdkAFpb2WwEWEoPIfdGfaogBONmWBMJ8asV1zCWJeDNonCCbMeEL4yk7pazoLwdATXY1IjGPY7wOblRo-jJWmgzVEeH2L0pB7d3gMuWR6cZ0ozfOGFU4CpMwzhfU4OyIdOwavuOjvrexM4daOYRGVwmySl0Pt5_urz5r4FHekMimXdRvfC3vrsmtgcH5DqM4Zi6YpuMmu_JEGmGvfegtuInId40p7NlrzpJIAxAoofzm0 "Fluxo para verificar se existe reprocessamento pendente")

Fluxo-03: Aprovar um reprocessamento para início imediato

![Fluxo-03](https://www.planttext.com/plantuml/img/bP2zYW9H38NxF4NAtK8GjXiBjH4i166nsoRCH0pSkRdatfdrTMIBVP1vCUDF8Ig2NVxETyYPvK9MkZO052c1SH6wlOx6NnNEasbFm__mfzWe6djVSyxKSYoAFn5NnBcOuZTRBpNx2E3C0wWskHiE9efqng1YCcbPx97K46ub4FP2E5yl9wvUQthccJWsNh2V-8gfG9KfE5sY-yRQ0OcCRdI6yKfl-1utUKzTQ_HdWvjVlF5t9nxO1-yb5tu1cNx2AHTD03D_mCC-0W00)





### Próximos Passos

- [x] Instalar Maestro na Plataforma
- [x] Conectar o Maestro ao Rabbit
- [x] Receber os eventos de persistência da Plataforma
- [x] Levantar uma API básica
- [x] Implementar um endpoint no domain para retornar a lista de entidades baseada na instancia do processo
- [x] Integrar com o Discovery para capturar as instancias que deverão ser reprocessadas
- [x] Implementar o serviço que permita ao admin e iniciar um reprocessamento
- [x] Implementar o serviço que permita ao admin ignorar um reprocessamento e por consequência ignorar um commit
- [x] Publicar instancias que deverão ser reprocessadas na fila de reprocessamento
- [x] API de check de bloqueio da Plataforma por sistema
- [x] Bloquear as execuções de persistência enquanto um reprocessamento estiver em execução
- [x] Faz o bloqueio de novos eventos no EventManager quando um reprocessamento está em execução
- [x] Incrementa a fila de reprocessamento com novos eventos (Reprocessamento em Cascata)
- [x] Incrementa o documento de reprocessamento com novos eventos (Reprocessamento em Cascata)
- [x] Fazer a persistência no domain quando oportuno

