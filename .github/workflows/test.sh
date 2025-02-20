#!/bin/bash

# Simulating the existing content
EXISTING_CONTENT="<table class=\"relative-table\" style=\"width: 21.25%;\"><colgroup><col style=\"width: 26.1538%;\" /><col style=\"width: 24.8718%;\" /><col style=\"width: 48.9744%;\" /></colgroup><tbody><tr><th>Release Version</th><th>Kueue</th><th>Codeflare-sdk</th></tr><tr><td>2.17</td><td>v0.8.1</td><td>v0.25.0</td></tr></tbody></table>"
echo "Existing Content: $EXISTING_CONTENT"

# Define new row
NEW_RELEASE="2.16"
KUEUE_VERSION="0.8"
CODEFLARE_VERSION="0.24.0"

NEW_ROW="<tr class=\"\"><td>$NEW_RELEASE</td><td>$CODEFLARE_VERSION</td><td>$KUEUE_VERSION</td></tr>"

echo "Creating new row: $NEW_ROW"

# Add new row before </tbody>
UPDATED_CONTENT=$(echo "$EXISTING_CONTENT" | sed "s|</tbody>|${NEW_ROW}</tbody>|")

# # Ensure proper escaping
# UPDATED_CONTENT=$(jq -n --arg content "$UPDATED_CONTENT" '$content')
# Debugging step
echo "UPDATED_CONTENT after row insertion: $UPDATED_CONTENT"

# Ensure proper JSON formatting
UPDATED_CONTENT=$(jq -Rs . <<< "$UPDATED_CONTENT")

echo "UPDATED_CONTENT is: $UPDATED_CONTENT"
