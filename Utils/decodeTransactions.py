import json, sys

if len(sys.argv) != 3: 
    print('Not enought arguments: \npython3 decodeTransactions.py <input_file_path> <output_json_path_file>')
    exit()

filePath = sys.argv[1]
jsonFilePath = sys.argv[2]

# Save data in a dict
transactionsData = []
with open(filePath, 'r') as file: 
    line = file.read() # Read hole file 
    transactions = line.split('#') # Every new transaction begin with '#'
    for transaction in transactions: 
        transaction = transaction.split('\x00')
        
        if len(transaction) == 1: 
            continue

        transactionId = transaction[0]
        buyerId = transaction[1]
        ip = transaction[2]
        device = transaction[3]
        products = transaction[4]
        # Remove parenthesis from products 
        products = products.replace('(', '') 
        products = products.replace(')', '')
        # Split into single products
        products = products.split(',')

        transactionDict = { 'id': transactionId, 'buyerId': buyerId, 'ip': ip, 'device': device, 'products': products}
        print(transactionDict)
        transactionsData.append(transactionDict) 

# Parse data into JSON format
with open(jsonFilePath, 'w', encoding='utf-8') as jsonFile: 
    json.dump(transactionsData, jsonFile, indent = 4, ensure_ascii = False)