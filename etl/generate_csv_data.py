import pandas as pd
import os

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
DATA_DIR = os.path.join(BASE_DIR, "data")

# ===================== MAIN DATA =====================
covid_data = "https://covid.ourworldindata.org/data/owid-covid-data.csv"
df = pd.read_csv(covid_data)
df = df[['iso_code', 'location', 'date', 'total_cases', 'total_deaths', 'people_vaccinated']]
df = df[df['iso_code'].str.len() == 3]  # Filter valid ISO country codes

# Assign unique IDs to each CovidCase and VaccinationStats entry (based on country/date rows)
df = df.reset_index(drop=True)
df['covidcase_id'] = range(1, len(df) + 1)
df['vaccstats_id'] = range(1, len(df) + 1)

# ===================== NODE: Country =====================
countries = df[['iso_code', 'location']].drop_duplicates().copy()
countries['id'] = range(1, len(countries) + 1)
countries.rename(columns={
    'iso_code': 'iso3',
    'location': 'name'
}, inplace=True)
countries = countries[['id', 'name', 'iso3']]
countries.to_csv(f"{DATA_DIR}/countries.csv", index=False)
print(f"Saving countries.csv with {len(countries)} rows...")

# ===================== NODE: CovidCase =====================
df_cases = df.dropna(subset=['date']).dropna(subset=['total_cases', 'total_deaths'], how='all').copy()
covid_cases = df_cases[['covidcase_id', 'iso_code', 'date', 'total_cases', 'total_deaths']].copy()
covid_cases.rename(columns={
    'covidcase_id': 'id',
    'iso_code': 'country_iso',
    'total_cases': 'totalCases',
    'total_deaths': 'totalDeaths'
}, inplace=True)
covid_cases.to_csv(f"{DATA_DIR}/covid_cases.csv", index=False)
print(f"Saving covid_cases.csv with {len(covid_cases)} rows...")

# ===================== NODE: VaccinationStats =====================
df_vacc = df.dropna(subset=['iso_code', 'date', 'people_vaccinated'], how='any').copy()
vacc_stats = df_vacc[['vaccstats_id', 'iso_code', 'date', 'people_vaccinated']].copy()
vacc_stats.rename(columns={
    'vaccstats_id': 'id',
    'iso_code': 'country_iso',
    'people_vaccinated': 'totalVaccinated'
}, inplace=True)
vacc_stats['totalVaccinated'] = vacc_stats['totalVaccinated'].astype(int)
vacc_stats.to_csv(f"{DATA_DIR}/vaccination_stats.csv", index=False)
print(f"Saving vaccination_stats.csv with {len(vacc_stats)} rows...")

# ===================== RELATIONSHIP: HAS_CASE =====================
has_case = covid_cases[['country_iso', 'id']].rename(columns={
    'id': 'covidcase_id'
})
has_case.to_csv(f"{DATA_DIR}/has_case.csv", index=False)
print(f"Saving has_case.csv with {len(has_case)} rows...")

# ===================== RELATIONSHIP: VACCINATED_ON =====================
vaccinated_on = vacc_stats[['country_iso', 'id']].rename(columns={
    'id': 'vaccstats_id'
})
vaccinated_on.to_csv(f"{DATA_DIR}/vaccinated_on.csv", index=False)
print(f"Saving vaccinated_on.csv with {len(vaccinated_on)} rows...")

# ===================== VACCINE MANUFACTURER DATA =====================
vac_manuf_url = "https://covid.ourworldindata.org/data/vaccinations/vaccinations-by-manufacturer.csv"
df_vac_by_manuf = pd.read_csv(vac_manuf_url)

vac_br = "https://covid.ourworldindata.org/data/vaccinations/country_data/Brazil.csv"
df_vac_br = pd.read_csv(vac_br)

# Set and unify Brazil data with Global data:
df_vac_br = df_vac_br[['date', 'vaccine']].dropna()
df_vac_br['vaccine'] = df_vac_br['vaccine'].str.split(', ')
df_vac_br = df_vac_br.explode('vaccine')
df_vac_br['location'] = 'Brazil'
df_vac_br = df_vac_br[['location', 'vaccine', 'date']]
df_vac_by_manuf = pd.concat([df_vac_by_manuf, df_vac_br], ignore_index=True)

# ===================== NODE: Vaccine =====================
vaccines = df_vac_by_manuf.groupby('vaccine')['date'].min().reset_index().copy()
vaccines['id'] = range(1, len(vaccines) + 1)
vaccines.rename(columns={'date': 'first_global_use'}, inplace=True)
vaccines.to_csv(f"{DATA_DIR}/vaccines.csv", index=False)
print(f"Saving vaccines.csv with {len(vaccines)} rows...")

# ===================== RELATIONSHIP: USES =====================
# Map country name to ISO code
country_iso_map = countries.set_index('name')['iso3'].to_dict()
uses_raw = df_vac_by_manuf[['location', 'vaccine', 'date']].copy()
uses_raw['country_iso'] = uses_raw['location'].map(country_iso_map)

total_uses = len(uses_raw)
uses_raw_valid = uses_raw.dropna(subset=['country_iso'])
missing_data = total_uses - len(uses_raw_valid)

if missing_data > 0:
    missing_countries = uses_raw[uses_raw['country_iso'].isna()]['location'].unique()
    print("WARNING! Entries were ignored due to country matching failure")
    print(f"Ignored countries: {sorted(missing_countries)}")

# Keep only the first use per country-vaccine pair
uses = uses_raw_valid.groupby(['country_iso', 'vaccine'])['date'].min().reset_index()
uses.rename(columns={'date': 'first_used'}, inplace=True)

uses.to_csv(f"{DATA_DIR}/uses.csv", index=False)
print(f"Saving uses.csv with {len(uses)} rows...")
