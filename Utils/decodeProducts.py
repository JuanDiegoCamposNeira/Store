import json, sys

csvFilePath =  '/Users/juancamposneira/Desktop/Reto-Aprendizaje/Data/Encoded/Products.txt'
jsonFilePath = '/Users/juancamposneira/Desktop/Reto-Aprendizaje/Data/Decoded/products.json'

# Save data in a dict
csvData = []
with open(csvFilePath, 'r') as csvFile: 
    lines = csvFile.readlines()
    for line in lines: 
        line = line.replace('\n', '') # Remove backspace at the end of the line
        product = line.split('\'')
        dict = { 'type': 'Product', 'uid': '_:' + product[0], 'name': product[1], 'price': product[2]}
        csvData.append(dict)

# Parse data into JSON format
with open(jsonFilePath, 'w', encoding='utf-8') as jsonFile: 
    json.dump(csvData, jsonFile, indent = 4, ensure_ascii = False)