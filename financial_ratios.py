import httpx

url = "https://stockanalysis.com/stocks/aapl/financials/ratios/__data.json?p=quarterly"
response = httpx.get(url)
raw_data = response.json()

# expected_keys = ['datekey', 'fiscalYear', 'fiscalQuarter', 'netIncomeCF', 'totalDepAmorCF', 'sbcomp', 'changeAR', 'changeInventory', 'changeAP', 'changeUnearnedRev', 'changeOtherNetOperAssets', 'otheroperating', 'ncfo', 'ocfGrowth', 'capex', 'cashAcquisition', 'salePurchaseIntangibles', 'investInSecurities', 'otherinvesting', 'ncfi', 'debtIssuedShortTerm', 'debtIssuedLongTerm', 'debtIssuedTotal', 'debtRepaidShortTerm', 'debtRepaidLongTerm', 'debtRepaidTotal', 'netDebtIssued', 'commonIssued', 'commonRepurchased', 'commonDividendCF', 'otherfinancing', 'ncff', 'ncf', 'fcf', 'fcfGrowth', 'fcfMargin', 'fcfps', 'leveredFCF', 'unleveredFCF', 'cashInterestPaid', 'cashTaxesPaid', 'changeNetWorkingCapital']
nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]
financial_data_index = data_map['financialData']
print(data[financial_data_index])

ratio_data = {}
ratio_data_map = data[financial_data_index]
for ratio_field, ratio_field_idx in ratio_data_map.items():
    ratio_field_values = []
    for field_value_index in data[ratio_field_idx]:
        ratio_field_values.append(data[field_value_index])
    
    ratio_data[ratio_field] = ratio_field_values

result = [dict(zip(ratio_data.keys(), values)) for values in zip(*ratio_data.values())]
print('-' * 100)
print(result[0].keys())