import os
import pandas as pd
from neo4j import GraphDatabase
from dotenv import load_dotenv

# Load environment variables from .env
load_dotenv()

NEO4J_URI = os.getenv("NEO4J_URI", "bolt://localhost:7687")
NEO4J_USER = os.getenv("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD")

DATA_DIR = "etl/data"

driver = GraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))

def load_csv_to_neo4j():
    with driver.session() as session:
        # Create indexes
        session.run("CREATE INDEX IF NOT EXISTS FOR (c:Country) ON (c.iso3)")
        session.run("CREATE INDEX IF NOT EXISTS FOR (cc:CovidCase) ON (cc.id)")
        session.run("CREATE INDEX IF NOT EXISTS FOR (vs:VaccinationStats) ON (vs.id)")
        session.run("CREATE INDEX IF NOT EXISTS FOR (v:Vaccine) ON (v.name)")

        # Load Country nodes
        for chunk in pd.read_csv(f"{DATA_DIR}/countries.csv", chunksize=1000):
            countries = chunk.to_dict(orient="records")
            session.run(
                """
                UNWIND $batch AS row
                MERGE (c:Country {iso3: row.iso3})
                SET c.name = row.name, c.id = toInteger(row.id)
                """,
                batch=countries
            )

        # Load CovidCase nodes
        for chunk in pd.read_csv(f"{DATA_DIR}/covid_cases.csv", chunksize=1000):
            covid_cases = chunk.to_dict(orient="records")
            session.run(
                """
                UNWIND $batch AS row
                MERGE (cc:CovidCase {id: toInteger(row.id)})
                SET cc.date = date(row.date),
                    cc.totalCases = CASE WHEN row.totalCases IS NOT NULL THEN toInteger(row.totalCases) ELSE NULL END,
                    cc.totalDeaths = CASE WHEN row.totalDeaths IS NOT NULL THEN toInteger(row.totalDeaths) ELSE NULL END
                """,
                batch=covid_cases
            )

        # Load VaccinationStats nodes
        for chunk in pd.read_csv(f"{DATA_DIR}/vaccination_stats.csv", chunksize=1000):
            vacc_stats = chunk.to_dict(orient="records")
            session.run(
                """
                UNWIND $batch AS row
                MERGE (vs:VaccinationStats {id: toInteger(row.id)})
                SET vs.date = date(row.date),
                    vs.totalVaccinated = toInteger(row.totalVaccinated)
                """,
                batch=vacc_stats
            )

        # Load Vaccine nodes
        for chunk in pd.read_csv(f"{DATA_DIR}/vaccines.csv", chunksize=1000):
            vaccines = chunk.to_dict(orient="records")
            session.run(
                """
                UNWIND $batch AS row
                MERGE (v:Vaccine {name: row.vaccine})
                SET v.id = toInteger(row.id),
                    v.first_global_use = date(row.first_global_use)
                """,
                batch=vaccines
            )

        # Relationships: HAS_CASE
        for chunk in pd.read_csv(f"{DATA_DIR}/has_case.csv", chunksize=1000):
            has_case = chunk.to_dict(orient="records")
            session.run(
                """
                UNWIND $batch AS row
                MATCH (c:Country {iso3: row.country_iso})
                MATCH (cc:CovidCase {id: toInteger(row.covidcase_id)})
                MERGE (c)-[:HAS_CASE]->(cc)
                """,
                batch=has_case
            )

        # Relationships: VACCINATED_ON
        for chunk in pd.read_csv(f"{DATA_DIR}/vaccinated_on.csv", chunksize=1000):
            vaccinated_on = chunk.to_dict(orient="records")
            session.run(
                """
                UNWIND $batch AS row
                MATCH (c:Country {iso3: row.country_iso})
                MATCH (vs:VaccinationStats {id: toInteger(row.vaccstats_id)})
                MERGE (c)-[:VACCINATED_ON]->(vs)
                """,
                batch=vaccinated_on
            )

        # Relationships: USES with attribute
        for chunk in pd.read_csv(f"{DATA_DIR}/uses.csv", chunksize=1000):
            uses = chunk.to_dict(orient="records")
            session.run(
                """
                UNWIND $batch AS row
                MATCH (c:Country {iso3: row.country_iso})
                MATCH (v:Vaccine {name: row.vaccine})
                MERGE (c)-[r:USES]->(v)
                SET r.first_used = date(row.first_used)
                """,
                batch=uses
            )

        print("Data successfully loaded into Neo4j.")

if __name__ == "__main__":
    load_csv_to_neo4j()
    driver.close()
