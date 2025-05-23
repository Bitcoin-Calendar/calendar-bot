# CSV Format Guide

This guide explains how to structure your CSV files for use with the Bitcoin Calendar Bot.

## Basic Structure

The CSV files must follow this specific format:

```csv
"date","title","description"
"YYYY-MM-DD","Event Title","Event Description |Optional additional text |Optional URLs"
```

### Required Columns

1. **Date** (Column 1): The date of the event in `YYYY-MM-DD` format
2. **Title** (Column 2): A short title for the event
3. **Description** (Column 3): Detailed description of the event, possibly with links

## Example

Here is an example of a properly formatted CSV file:

```csv
"date","title","description"
"2009-01-03","Genesis Block Mined","Satoshi Nakamoto mines the first Bitcoin block (genesis block), starting the Bitcoin network. |The block contained the text: 'The Times 03/Jan/2009 Chancellor on brink of second bailout for banks' |https://en.bitcoin.it/wiki/Genesis_block"
"2010-05-22","Bitcoin Pizza Day","Laszlo Hanyecz makes the first commercial transaction using Bitcoin, paying 10,000 BTC for two pizzas. |The pizzas were worth about $25, making the Bitcoin-to-USD exchange rate roughly $0.0025 per BTC. |https://bitcointalk.org/index.php?topic=137.0"
"2011-02-09","Bitcoin reaches parity with USD","For the first time in history, 1 BTC = 1 USD on the Mt. Gox exchange. |This was a significant psychological milestone for the cryptocurrency."
```

## Special Formatting

### Line Breaks

To include line breaks in your descriptions:

- Use the pipe character (`|`) to separate paragraphs
- Each pipe character will be converted to a double line break in the final post
- Example: `"First paragraph |Second paragraph |Third paragraph"`

### Date Matching

The bot checks for events matching today's date by comparing the month and day portions:

- The bot extracts the month and day from the date (e.g., `01-03` from `2009-01-03`)
- It compares this with today's date in the same format
- Events are posted only if they match today's month and day

### Best Practices

1.	**Column Formatting**: Enclose values in double quotes to handle special characters properly.
2.	**Quote Usage**: Use quotes only to define column boundaries; replace internal quotes with "" if needed.
3.	**Consistent Formatting**: Maintain uniform formatting across CSV files.
4.	**Verification**: Ensure all URLs are valid and functional.
5.	**Historical Accuracy**: Verify the correctness of all events and dates.
6.	**File Naming**: Follow these naming conventions for your CSV files:
   - English events: `events_en.csv`
   - Russian events: `events_ru.csv`
   - Other languages: `events_xx.csv` (where `xx` is the language code)

## Language-Specific Considerations

### English (`events_en.csv`)

- Use standard English formatting and punctuation
- Provide links to English-language resources

### Russian (`events_ru.csv`)

- Use Russian language and appropriate Cyrillic encoding
- Provide links to Russian-language resources when available

## CSV File Organization

To add or modify events:

1. Open the CSV file in a text editor (VSCode works, while many editors that support CSV will break formatting when exporting your edited version)
2. Add new events or edit existing ones
3. Save the file in CSV format with UTF-8 encoding
4. Ensure the file has the correct name and is in the proper location

## Validation

Before using a CSV file with the bot, validate that:

1. The file uses the correct format with three columns
2. All dates follow the YYYY-MM-DD format
3. All required quotes and separators are present
4. The file is encoded with UTF-8

## Additional Notes

- The first line of the CSV file must contain the column headers
- The bot processes events from the second line onward
- Events with the same date will all be posted, with a 30-minute delay between posts 