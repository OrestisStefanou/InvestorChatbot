import httpx

url = 'https://stockanalysis.com/stocks/nvda/forecast/__data.json'
response = httpx.get(url)
raw_data = response.json()

nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]

# ESTIMATES SCRAPING
estimates_data_index = data_map['estimates']
estimates_table_data_index = data[estimates_data_index]['table']
quarterly_estimates_data_index = data[estimates_table_data_index]['quarterly']
quarterly_estimates_data_map = data[quarterly_estimates_data_index]
quarterly_estimates_data = {}
for estimation_field, estimation_field_idx in quarterly_estimates_data_map.items():
    if estimation_field == 'lastDate':
        continue
    estimation_field_values = []
    for field_value_index in data[estimation_field_idx]:
        if data[field_value_index] == '[PRO]':
            continue

        estimation_field_values.append(data[field_value_index])
    
    quarterly_estimates_data[estimation_field] = estimation_field_values

estimations_doc = [dict(zip(quarterly_estimates_data.keys(), values)) for values in zip(*quarterly_estimates_data.values())]
print(estimations_doc)
print('-------------------------------')

# TARGET PRICE SCRAPING
target_data_index = data_map['targets']
target_data_map = data[target_data_index]
target_price_keys = ['average', 'high', 'low', 'median']
target_price_doc = {}
for target_key in target_price_keys:
    target_value_index = target_data_map[target_key]
    target_price_doc[target_key] = data[target_value_index]

print(target_price_doc)