import httpx

url = 'https://stockanalysis.com/stocks/aapl/forecast/__data.json'
response = httpx.get(url)
raw_data = response.json()

nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]

print(data_map)
print(data[187])    # targets data
print(data[1])  # estimates data

estimates_data_index = data_map['estimates']
estimates_table_data_index = data[estimates_data_index]['table']
quarterly_estimates_data_index = data[estimates_table_data_index]['quarterly']
quarterly_estimates_data_map = data[quarterly_estimates_data_index]
quarterly_estimates_data = {}
for estimation_field, estimation_field_idx in quarterly_estimates_data_map.items():
    if estimation_field == 'lastDate':
        continue
    print(estimation_field, estimation_field_idx)
    estimation_field_values = []
    for field_value_index in data[estimation_field_idx]:
        if data[field_value_index] == '[PRO]':
            continue

        estimation_field_values.append(data[field_value_index])
    
    quarterly_estimates_data[estimation_field] = estimation_field_values

estimations_doc = [dict(zip(quarterly_estimates_data.keys(), values)) for values in zip(*quarterly_estimates_data.values())]
print('-' * 100)
print(estimations_doc)