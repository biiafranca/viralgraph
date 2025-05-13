# ViralGraph 🦠🌐

ViralGraph é um projeto de análise e visualização de dados relacionados à pandemia de COVID-19, utilizando o banco de dados orientado a grafos Neo4j. O projeto realiza o processamento e carga de dados via ETL Python, e fornece uma API REST em Go para consulta aos dados.

Ele foi desenvolvido como resposta ao desafio técnico proposto em:  
📎 https://github.com/NeowayLabs/jobs/blob/master/graph-analysis/analyst.md

O objetivo é construir uma API REST capaz de responder às seguintes perguntas:

1. **Qual foi o total acumulado de casos e mortes de Covid-19 em um país específico em uma data determinada?**  
2. **Quantas pessoas foram vacinadas com pelo menos uma dose em um determinado país em uma data específica?**  
3. **Quais vacinas foram usadas em um país específico?**  
4. **Em quais datas as vacinas foram autorizadas para uso?**  
5. **Quais países usaram uma vacina específica?**  

## 🔧 Tecnologias

- Go (Golang)
- Neo4j 5.x
- Docker + Docker Compose
- Python (ETL com pandas)
- Makefile
- go-chi (roteador leve para Go)

## 🗂 Estrutura

```
ViralGraph/
├── api/                    # API REST em Go
├── etl/                    # Scripts de ETL em Python
├── docker-compose.yml
├── Makefile
├── .gitignore
├── .env
└── README.md
````

## 🚀 Como rodar o projeto

### Pré-requisitos

- Docker e Docker Compose instalados
- `make` disponível no terminal (Linux/Mac ou via [Chocolatey](https://chocolatey.org/) no Windows)

### Passos

1. Inicie todos os serviços com:

```
make start
```

2. Acesse a API em http://localhost:8080 e o Neo4j Browser em http://localhost:7474

## 🧪 Testes

Para executar os testes automatizados da API:
```
make test
```

## 📦 ETL

Gera arquivos CSV e carrega dados no Neo4j.

Para mais informações, consulte /etl/README.md

## 📖 API

A especificação da API (Swagger/OpenAPI) está disponível em api/docs/swagger.yaml.

Para mais informações, consulte o /api/README.md

## Respostas ao desafio e Exemplos de uso

Seguem os exemplos de uso, para obter respostas para as perguntas solicitadas no desafio. A API também responde algumas perguntas **bônus**, que também serão exemplificadas:

1. **Qual foi o total acumulado de casos e mortes de Covid-19 em um país específico em uma data determinada?**  
   → Rota:`GET /covid-stats/{country}/{date}`
   → Ex:`GET /covid-stats/BRA/2021-08-01`

2. **Quantas pessoas foram vacinadas com pelo menos uma dose em um determinado país em uma data específica?**  
   → Rota:`GET /vaccinations/{country}/{date}`
   → Ex:`GET /vaccinations/BRA/2021-08-01`

3. **Quais vacinas foram usadas em um país específico?**  
   → Rota: `GET /vaccines/used-in/{country}`
   → Ex: `GET /vaccines/used-in/BRA`

4. **Em quais datas as vacinas foram autorizadas para uso?**  
   → Rota: `GET /vaccines/first-use`

5. **Quais países usaram uma vacina específica?**  
   → Rota: `GET /vaccines/{vaccine_id}/used-by`
   → Ex: `GET /vaccines/1/used-by`

6. **Qual foi o primeiro uso de cada vacina em um país específico?**  
   → Incluso na resposta de `/vaccines/used-in/{country}` com campo `first_used`
   → Ex: `GET /vaccines/used-in/BRA`

7. **Qual a quantidade de novos casos e mortes de Covid-19 em um país específico registrados em uma data determinada?**  
   → Rota: `GET /covid-stats/{country}/{date}?only-news=true`
   → Ex: `GET /covid-stats/BRA/2021-08-01?only-news=true`

8. **Qual o total (acumulado ou diário) de casos e mortes de Covid-19 a nível mundial?**  
   → Rota: `GET /covid-stats/{date}` (acumulado)
   → Ex: `GET /covid-stats/2021-08-01` (acumulado)  
   → Rota: `GET /covid-stats/{date}?only-news=true` (novos)
   → Ex: `GET /covid-stats/2021-08-01?only-news=true` (novos)

9. **Qual o total de novas vacinações de Covid-19 em um país específico registrados em uma data determinada?**  
   → Rota: `GET /vaccinations/{country}/{date}?only-news=true`
   → Ex: `GET /vaccinations/BRA/2021-08-01?only-news=true`

10. **Qual o total (acumulado ou diário) de vacinações de Covid-19 a nível mundial?**  
   → `GET /vaccinations/{date}` (acumulado) 
   → `GET /vaccinations/2021-08-01` (acumulado)   
   → `GET /vaccinations/{date}?only-news=true` (novos)
   → `GET /vaccinations/2021-08-01?only-news=true` (novos)

Para testar as respostas às perguntas, utilize a coleção e configuração de ambiente Postman, presentes em /api/docs/
