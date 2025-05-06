# ETL do Projeto ViralGraph

Este módulo é responsável por extrair, transformar e carregar os dados sobre a pandemia de COVID-19 em um banco de dados de grafos Neo4j. 

---

## 🔧 Tecnologias utilizadas

- **Python**
- **Pandas** para manipulação de dados
- **Neo4j** como banco de grafos
- **Driver oficial Neo4j (`neo4j`)** para conexão via Bolt
- **dotenv** para gerenciar variáveis sensíveis

---

## 📁 Estrutura

```
etl/
├── data/                      # Contém os CSVs gerados
├── generate_csv_data.py       # Extrai e transforma os dados em CSVs
├── load_to_neo4j.py           # Carrega os dados no Neo4j com UNWIND
├── README.md
├── requirements.txt
├── .env
```

---

## ⚙️ Configuração

Crie um arquivo `.env` na raiz do projeto com os seguintes parâmetros:

```
NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=sua_senha_aqui
```

---

📦 Instalação das dependências

Execute o seguinte comando para instalar os pacotes necessários à execução da ETL:

```
pip install -r etl/requirements.txt
```

---

## 🚀 Etapas do ETL

### 1. Geração dos CSVs

Execute:
```bash
python etl/generate_graph_data.py
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

Execute:
```bash
python etl/load_to_neo4j.py
```

Este script:
- Cria índices nos atributos principais
- Utiliza **UNWIND** para enviar os dados em lote
- Converte datas para o tipo `date` do Neo4j
- Usa `MERGE` para evitar duplicatas e `SET` para atualizar atributos

---

## 💡 Decisões técnicas

🔸 Fonte de dados escolhida (OWID)

A base de dados da Our World in Data (OWID) foi escolhida por ser amplamente reconhecida, de acesso público, e por disponibilizar dados consistentes e atualizados sobre casos de COVID-19, mortes, vacinação e uso por fabricante. Ela oferece formatos acessíveis (CSV via URL), estrutura padronizada e boa documentação, facilitando a automação da ETL e garantindo confiabilidade nas análises.

🔸 Linguagem escolhida para a ETL

A linguagem Python foi escolhida para a construção da ferramenta de ETL devido à sua ampla maturidade no ecossistema de manipulação de dados, especialmente com bibliotecas como pandas. Apesar de o projeto utilizar Go em outras partes, o Python proporcionou maior agilidade no tratamento e transformação dos dados tabulares.

🔸 Identificadores numéricos (id)

Foi escolhida a utilização de identificadores numéricos inteiros nos nós que exigem unicidade (CovidCase, VaccinationStats, etc.), garantindo maior performance e padronização. Essa decisão também facilita integrações futuras com bancos relacionais.

🔸 Inclusão de ID em nós que não são identificados por ele

Mesmo nos nós Country e Vaccine, onde a identificação primária é feita por iso3 e name, respectivamente, o campo id foi mantido por consistência e para viabilizar eventuais expansões que exijam vínculos relacionais numéricos.

🔸 Modelagem de relacionamentos Country → CovidCase

Embora o modelo sugerido inicialmente relacionasse Country e CovidCase apenas pela data, essa abordagem foi considerada inadequada, uma vez que vários países podem ter registros para a mesma data. Assim, o relacionamento foi implementado com base no identificador exclusivo de cada CovidCase, garantindo integridade e unicidade no grafo.

🔸 Uso de MERGE + SET

Optamos por MERGE (a)-[r:REL]->(b) seguido de SET r.prop = ... para evitar múltiplos relacionamentos com atributos diferentes e permitir atualizações sem duplicações.

🔸 Conversão de datas

As datas foram explicitamente convertidas com date(...) no Cypher, para permitir consultas com filtros de data no formato nativo do Neo4j.

🔸 Performance com UNWIND

Toda a carga de dados foi otimizada com UNWIND, reduzindo drasticamente o número de comandos e aumentando a escalabilidade do processo.

🔸 Eliminação do nó VaccineApproval

Consideramos desnecessária a existência de um nó exclusivo para representar o evento de aprovação de uma vacina, especialmente porque:

- Ele não traria nenhuma propriedade relevante além da data
- Não haveria relacionamento com os países, tornando o nó isolado
- A informação de aprovação poderia ser mais bem representada como atributo no próprio nó Vaccine

Além disso, como a fonte de dados não fornece a data oficial da aprovação regulatória, e sim a data do primeiro uso documentado, usamos esta data como proxy para first_global_use. Já o uso específico por país foi modelado como atributo first_used no relacionamento (:Country)-[:USES]->(:Vaccine).
Essa modelagem simplifica o grafo, evita nós artificiais e mantém a capacidade de responder às perguntas do desafio de forma clara.

🔸 Uso do chunksize ao carregar os dados

A decisão de usar o parâmetro chunksize foi baseada na necessidade de otimizar a performance e evitar sobrecarga de memória no banco de dados Neo4j. Ao processar os dados em lotes controlados, em vez de usar um tamanho indefinido, garantimos que o sistema não sobrecarregue sua memória ao tentar carregar grandes volumes de dados simultaneamente. Essa abordagem também melhora a eficiência ao balancear as operações de leitura e escrita, facilita a escalabilidade ao lidar com grandes volumes de dados e oferece maior controle sobre erros, permitindo a recuperação eficiente sem comprometer todo o processo.

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
- A limpeza de dados alfabéticos nos campos numéricos não foi aplicada, mas é prevista como melhoria futura.

---

Para dúvidas ou contribuições, consulte o `README.md` principal do projeto.
