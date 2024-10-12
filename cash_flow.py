import httpx

url = "https://stockanalysis.com/stocks/aapl/financials/cash-flow-statement/__data.json?p=quarterly"
response = httpx.get(url)
raw_data = response.json()

expected_keys = ['datekey', 'fiscalYear', 'fiscalQuarter', 'netIncomeCF', 'totalDepAmorCF', 'sbcomp', 'changeAR', 'changeInventory', 'changeAP', 'changeUnearnedRev', 'changeOtherNetOperAssets', 'otheroperating', 'ncfo', 'ocfGrowth', 'capex', 'cashAcquisition', 'salePurchaseIntangibles', 'investInSecurities', 'otherinvesting', 'ncfi', 'debtIssuedShortTerm', 'debtIssuedLongTerm', 'debtIssuedTotal', 'debtRepaidShortTerm', 'debtRepaidLongTerm', 'debtRepaidTotal', 'netDebtIssued', 'commonIssued', 'commonRepurchased', 'commonDividendCF', 'otherfinancing', 'ncff', 'ncf', 'fcf', 'fcfGrowth', 'fcfMargin', 'fcfps', 'leveredFCF', 'unleveredFCF', 'cashInterestPaid', 'cashTaxesPaid', 'changeNetWorkingCapital']
nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]
financial_data_index = data_map['financialData']
print(data[financial_data_index])

cash_flow_data = {}
cash_flow_data_map = data[financial_data_index]
for cash_flow_field, cash_flow_field_idx in cash_flow_data_map.items():
    cash_flow_field_values = []
    for field_value_index in data[cash_flow_field_idx]:
        cash_flow_field_values.append(data[field_value_index])
    
    cash_flow_data[cash_flow_field] = cash_flow_field_values

result = [dict(zip(cash_flow_data.keys(), values)) for values in zip(*cash_flow_data.values())]
print('-' * 100)
print(result[0].keys())