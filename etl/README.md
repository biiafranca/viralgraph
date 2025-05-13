# ETL do Projeto ViralGraph

Este módulo é responsável por extrair, transformar e carregar os dados sobre a pandemia de COVID-19 em um banco de dados de grafos Neo4j. 


## 🔧 Tecnologias utilizadas

- **Python**
- **Pandas** para manipulação de dados
- **Neo4j** como banco de grafos
- **Driver oficial Neo4j (`neo4j`)** para conexão via Bolt

## 📁 Estrutura

```
etl/
├── data/                      # Contém os CSVs gerados
├── generate_csv_data.py       # Extrai e transforma os dados em CSVs
├── load_to_neo4j.py           # Carrega os dados no Neo4j com UNWIND
├── README.md
├── requirements.txt
```

## ⚙️ Configuração

Crie um arquivo `.env` na raiz do projeto com os seguintes parâmetros:

```
NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=sua_senha_aqui
```

📦 Instalação das dependências

Para execução local, execute o seguinte comando para instalar os pacotes necessários à execução da ETL:

```
pip install -r etl/requirements.txt
```

## 🚀 Etapas do ETL

### 1. Geração dos CSVs

Localmente, execute:
```bash
python etl/generate_graph_data.py
```

Para execução via docker-compose:
```
make etl-generate:
```

Este script:
- Baixa os dados atualizados do [Our World in Data](https://ourworldindata.org/covid-vaccinations)
- Filtra dados inválidos
- Gera nós para:
  - Country
  - CovidCase
  - VaccinationStats
  - Vaccine
- Gera relacionamentos:
  - HAS_CASE
  - VACCINATED_ON
  - USES (com atributo `first_used`)

### 2. Carga no Neo4j

Localmente, execute:
```bash
python etl/load_to_neo4j.py
```

Para execução via docker-compose:
```
make etl-load
```

Este script:
- Cria índices nos atributos principais
- Utiliza **UNWIND** para enviar os dados em lote
- Converte datas para o tipo `date` do Neo4j
- Usa `MERGE` para evitar duplicatas e `SET` para atualizar atributos


## 💡 Decisões técnicas

🔸 Linguagem escolhida para a ETL

A linguagem Python foi escolhida para a construção da ferramenta de ETL devido à sua maturidade no ecossistema de manipulação de dados, especialmente com a biblioteca pandas. Apesar de o projeto utilizar Go no desenvolvimento da API, o Python proporcionou maior agilidade no tratamento e transformação dos dados.

🔸 Fonte de dados escolhida (OWID)

A base de dados da Our World in Data (OWID) foi escolhida por ser amplamente reconhecida, de acesso público, e por disponibilizar dados consistentes e atualizados sobre casos de COVID-19, mortes, vacinação e uso por fabricante. Ela oferece formatos acessíveis (CSV via URL), estrutura padronizada e boa documentação, facilitando a automação da ETL e garantindo confiabilidade nas análises.

Foram extraídos os seguintes dados:

- owid-covid-data.csv:
  - 'iso_code': Código iso3 do país
  - 'location': Nome do país, em inglês
  - 'date': Data relativa ao dado
  - 'total_cases': total acumulado de casos na data
  - 'total_deaths': total acumulado de mortes na data
  - 'people_vaccinated': pessoas vacinadas com no mínimo uma dose da vacina, total acumulado
    
- vaccinations-by-manufacturer.csv:
  - 'location': Nome do país, em inglês
  - 'date': Data relativa ao dado 
  - 'vaccine': Nome da vacina

Como o vaccinations-by-manufacturer.csv não contém informações relativas a todos os países, os dados foram enriquecidos com informações específicas do Brasil, para demonstração.

- country_data/Brazil.csv:
  - 'date': Data relativa ao dado 
  - 'vaccine': Nome da vacina

🔸 Identificadores numéricos (id)

Foi escolhida a utilização de identificadores numéricos inteiros, pois, de acordo com as pesquisas realizadas, no banco de dados Neo4J a ordenação, indexação e busca por igualdade são muito rápidas com números do que com strings, byte a byte.

🔸 Inclusão de ID em nós que não são identificados por ele

Mesmo nos nós Country e Vaccine, onde a identificação primária é feita por iso3 e name, respectivamente, o campo id foi mantido por consistência e para viabilizar eventuais expansões que exijam vínculos relacionais numéricos.

🔸 Modelagem de relacionamentos Country → CovidCase

Embora o modelo sugerido inicialmente relacionasse Country e CovidCase apenas pela data, essa abordagem foi considerada inadequada, uma vez que vários países podem ter registros para a mesma data. Assim, o relacionamento considerou também o iso3 do país.

🔸 Uso de MERGE + SET

Foi utilizado o MERGE (a)-[r:REL]->(b) seguido de SET r.prop = ... para permitir atualizações sem duplicações.

🔸 Conversão de datas

As datas foram explicitamente convertidas com date(...) no Cypher, para permitir consultas com filtros de data no formato nativo do Neo4j.

🔸 Performance com UNWIND e chunksize ao carregar os dados

Toda a carga de dados foi otimizada com UNWIND, reduzindo drasticamente o número de comandos e aumentando a escalabilidade do processo. Foi definido um chunksize de 1000 linhas em cada lote carregado, para evitar sobrecarga de memória no banco de dados Neo4j. Ao processar os dados em lotes controlados, em vez de usar um tamanho indefinido, garantimos que o sistema não sobrecarregue sua memória ao tentar carregar grandes volumes de dados simultaneamente 

🔸 Eliminação do nó VaccineApproval

Não foi utilizado um nó exclusivo para representar o evento de aprovação de uma vacina, especialmente porque:

- Ele não traria nenhuma propriedade relevante além da data
- Não haveria relacionamento com os países, tornando o nó isolado
- A informação de aprovação poderia ser mais bem representada como atributo no próprio nó Vaccine

Além disso, como a fonte de dados não fornece a data oficial da aprovação regulatória, e sim a data do primeiro uso documentado, foi utilizada esta data para inferir o dado como first_global_use. Já o uso específico por país foi modelado como atributo first_used no relacionamento (:Country)-[:USES]->(:Vaccine).
Essa modelagem simplifica o grafo, evita nós artificiais e mantém a capacidade de responder às perguntas do desafio de forma clara.

🔸 Melhorias Futuras: Enriquecimento com Outras Fontes

Como aprimoramento futuro, é possível realizar o enriquecimento dos dados com fontes alternativas oficiais, como o Ministério da Saúde do Brasil ou bancos de dados regionais com cobertura mais precisa. Essa melhoria traria maior representatividade e completude à análise global. Contudo, essa etapa foi intencionalmente deixada de fora do escopo original proposto, a fim de manter o foco na implementação da arquitetura da API, modelagem do grafo e demonstração de consultas relevantes sobre os dados já fornecidos.

---

## ✅ Exemplo de query possível após a carga

```cypher
MATCH (c:Country {iso3: "BRA"})-[:HAS_CASE]->(cc:CovidCase)
WHERE cc.date = date("2021-03-01")
RETURN cc.totalCases, cc.totalDeaths
```

---

## 📌 Observações finais

- O script foi testado com mais de 450 mil registros e manteve estabilidade.
- Em caso de dúvidas, consulte também o `README.md` principal do projeto.
