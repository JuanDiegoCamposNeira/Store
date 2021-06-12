import json, sys

if len(sys.argv) != 3: 
    print('Not enought arguments: \npython3 decodeCSV.py <input_csv_file_path> <output_json_path_file>')
    exit()

csvFilePath = sys.argv[1]
jsonFilePath = sys.argv[2]

# Save data in a dict
csvData = []
with open(csvFilePath, 'r') as csvFile: 
    lines = csvFile.readlines()
    for line in lines: 
        line = line.replace('\n', '') # Remove backspace at the end of the line
        product = line.split('\'')
        dict = { 'id': product[0], 'name': product[1], 'price': product[2]}
        csvData.append(dict)

# Parse data into JSON format
with open(jsonFilePath, 'w', encoding='utf-8') as jsonFile: 
    json.dump(csvData, jsonFile, indent = 4, ensure_ascii = False)