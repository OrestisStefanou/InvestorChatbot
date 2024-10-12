import httpx


url = "https://stockanalysis.com/stocks/nvda/financials/balance-sheet/__data.json?p=quarterly"
response = httpx.get(url)
raw_data = response.json()

expected_keys = ['datekey', 'fiscalYear', 'fiscalQuarter', 'cashneq', 'investmentsc', 'totalcash', 'cashGrowth', 'accountsReceivable', 'otherReceivables', 'receivables', 'inventory', 'restrictedCash', 'othercurrent', 'assetsc', 'netPPE', 'investmentsnc', 'goodwill', 'otherIntangibles', 'othernoncurrent', 'assets', 'accountsPayable', 'accruedExpenses', 'debtc', 'currentPortDebt', 'currentCapLeases', 'currentIncomeTaxesPayable', 'currentUnearnedRevenue', 'otherCurrentLiabilities', 'currentLiabilities', 'debtnc', 'capitalLeases', 'longTermUnearnedRevenue', 'longTermDeferredTaxLiabilities', 'otherliabilitiesnoncurrent', 'liabilities', 'commonStock', 'retearn', 'otherEquity', 'equity', 'liabilitiesequity', 'sharesOutFilingDate', 'sharesOutTotalCommon', 'bvps', 'tangibleBookValue', 'tangibleBookValuePerShare', 'debt', 'netcash', 'netCashGrowth', 'netcashpershare', 'workingcapital', 'land', 'machinery', 'leaseholdImprovements', 'tradingAssetSecurities']

nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]
financial_data_index = data_map['financialData']
print(data[financial_data_index])

balance_sheet_data = {}
balance_sheet_data_map = data[financial_data_index]
for balance_sheet_field, balance_sheet_field_idx in balance_sheet_data_map.items():
    balance_sheet_field_values = []
    for field_value_index in data[balance_sheet_field_idx]:
        balance_sheet_field_values.append(data[field_value_index])
    
    balance_sheet_data[balance_sheet_field] = balance_sheet_field_values

result = [dict(zip(balance_sheet_data.keys(), values)) for values in zip(*balance_sheet_data.values())]
print('-' * 100)
print(result[0].keys())