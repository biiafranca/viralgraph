# API REST - ViralGraph

Esta API fornece acesso a informações sobre estatísticas de COVID-19 e vacinas, consultando dados carregados no banco Neo4j.

## 🔌 Endpoints

### 📊 Estatísticas da COVID-19

- GET `/covid-stats/{country}/{date}`
- GET `/covid-stats/{date}`  
  → Parâmetro opcional: `only-news=true` (retorna apenas casos/mortes com registro no dia solicitado)

### 💉 Vacinação

- GET `/vaccinations/{country}/{date}`
- GET `/vaccinations/{date}`  
  → Suporta `only-news=true` também

### 🌍 Uso de vacinas

- GET `/vaccines`
  → Retorna todas as vacinas cadastradas no banco de dados

- GET `/vaccines/first-use`  
  → Retorna o primeiro uso global de cada vacina

- GET `/vaccines/{vaccine_id}/used-by`  
  → Retorna os países que usaram a vacina
  
- GET `/vaccines/used-in/{country}`  
  → Retorna vacinas aplicadas no país

## 🏗 Estrutura

´´´
api/
├── main.go
├── routes/          # Registro de rotas
├── handlers/        # Implementação dos endpoints
├── neo4j/           # Acesso ao banco
├── utils/           # Funções auxiliares
└── docs/            # Swagger/OpenAPI e Postman
´´´

## 🧪 Testes

Os testes são feitos com `go test`, para executar:
   ```
   make test
   ```

## 💡 Decisões técnicas

## 🧠 Decisões Técnicas

🔸 **Go como linguagem da API**:

A linguagem Go (ou Golang), escolhida para a construção da API do ViralGraph, é rápida, leve e fácil de manter. 
- Vem com recursos nativos importantes, como servidor HTTP, testes e concorrência.
- Entrega respostas rápidas com uso eficiente de CPU e memória, ideal para serviços REST com baixa latência.
- Força um código limpo e padrão com ferramentas como gofmt.

🔸 **go-chi como roteador**: 

A escolha do go-chi foi motivada principalmente pela necessidade de trabalhar com rotas dinâmicas de forma limpa e produtiva.

Antes da adoção do go-chi, o roteamento era feito de forma manual, o que exigia decompor a URL como string e extrair os parâmetros "na mão". Isso era propenso a erro e dificultava a leitura e a manutenção do código, especialmente quando as rotas tinham múltiplos parâmetros.

Com o go-chi, se obteve:

- Suporte nativo a rotas parametrizadas (/{param}), com fácil extração via chi.URLParam(r, "param")
- Um roteador minimalista e leve, sem trazer dependências desnecessárias
- Composição modular de rotas — cada grupo de endpoints pode ser registrado separadamente, favorecendo a organização por domínio (covidstats, vaccinations, etc.)

🔸**Limitação dos testes**:

Por simplicidade inicial, os testes possuem as seguintes limitações:

- A cobertura de testes atual inclui apenas a verificação de status de resposta de alguns cenários positivos e negativos. Como melhoria futura, recomenda-se uma maior cobertura, inclusive de verificação da estrutra das respostas.
- Os testes estão consultando o banco de dados real, não sendo propriamente testes unitários. O uso de mock é recomendado neste caso, e está incluso como melhoria futura.

🔸**Documentação .yaml estática**

Também por simplicidade inicial, a documentação .yaml modelo OpenAPI foi gerada de forma estática. Porém, para aplicação futura, é possível implementar bibliotecas como a Swaggo, capaz de gerar a documentação automaticamente, a atualizando facilmente no decorrer do desenvolvimento da API.

🔸**Organização modular da API**: 

Rotas e handlers foram separados por domínio (`covidstats`, `vaccinations`, `vaccines`).
