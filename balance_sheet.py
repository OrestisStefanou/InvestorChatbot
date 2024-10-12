import json

# Open and read the JSON file
with open('example_responses/balance_sheet.json', 'r') as file:
    raw_data = json.load(file)

# Print the data
nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]
financial_data_index = data_map['financialData']
# print(financial_data_index)
print(data[financial_data_index])
# print(data[231])
# print(data[919])

balance_sheet_data = {}
balance_sheet_data_map = data[financial_data_index]
for key, value in balance_sheet_data_map.items():
    key_data = []
    for data_index in data[value]:
        key_data.append(data[data_index])
    
    balance_sheet_data[key] = key_data

print('-' * 100)
print(balance_sheet_data)