# Open Cert
Projeto que eu comecei a desenvolver como software real (não de portfólio) de uma plataforma para registrar e validar conhecimentos profissionais.

O projeto não chegou a implementar a plataforma completa. O trabalho ficou concentrado no backend e autenticação. Nessa época eu acabei iniciando o projeto implementando sistema de autenticação. 

## Implementação
O backend foi desenvolvido em Go e organizado em camadas, separando transporte HTTP, serviços de domínio, entidades, repositórios e infraestrutura. A composição das dependências é feita na inicialização do serviço.

Os fluxos implementados incluem cadastro e gerenciamento de contas, autenticação com tokens JWT de acesso e renovação, autenticação em duas etapas com TOTP e operações relacionadas a avaliadores e experiências profissionais. O serviço também integra persistência compatível com PostgreSQL/CockroachDB, cache com Memcached e mensageria com RabbitMQ.
