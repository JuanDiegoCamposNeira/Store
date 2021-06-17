import json, sys

filePath = '/Users/juancamposneira/Desktop/Reto-Aprendizaje/Data/Encoded/transactions.txt'
jsonFilePath = '/Users/juancamposneira/Desktop/Reto-Aprendizaje/Data/Decoded/transactions.json'

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

        transactionDict = { 'uid': transactionId, 'buyerId': buyerId, 'ip': ip, 'device': device, 'products': products}
        transactionsData.append(transactionDict) 

# Parse data into JSON format
with open(jsonFilePath, 'w', encoding='utf-8') as jsonFile: 
    json.dump(transactionsData, jsonFile, indent = 4, ensure_ascii = False)