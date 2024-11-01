import httpx

url = "https://stockanalysis.com/stocks/aapl/financials/__data.json?p=quarterly"
response = httpx.get(url)
raw_data = response.json()

expected_keys = ['datekey', 'fiscalYear', 'fiscalQuarter', 'revenue', 'revenueGrowth', 'cor', 'gp', 'sgna', 'rnd', 'opex', 'opinc', 'interestExpense', 'interestIncome', 'currencyGains', 'otherNonOperating', 'ebtExcl', 'gainInvestments', 'pretax', 'taxexp', 'netinc', 'netinccmn', 'netIncomeGrowth', 'sharesBasic', 'sharesDiluted', 'sharesYoY', 'epsBasic', 'epsdil', 'epsGrowth', 'fcf', 'fcfps', 'dps', 'dividendGrowth', 'grossMargin', 'operatingMargin', 'profitMargin', 'fcfMargin', 'taxrate', 'ebitda', 'depAmorEbitda', 'ebitdaMargin', 'ebit', 'ebitMargin', 'revenueAsReported', 'payoutratio']

nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]
financial_data_index = data_map['financialData']
print(data[financial_data_index])

income_statement_data = {}
income_statement_data_map = data[financial_data_index]
for income_statement_field, income_statement_field_idx in income_statement_data_map.items():
    income_statement_field_values = []
    for field_value_index in data[income_statement_field_idx]:
        income_statement_field_values.append(data[field_value_index])
    
    income_statement_data[income_statement_field] = income_statement_field_values

result = [dict(zip(income_statement_data.keys(), values)) for values in zip(*income_statement_data.values())]
print('-' * 100)
print(result[0].keys())