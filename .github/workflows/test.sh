#!/bin/bash

set -e  # Exit on error

# Set the values to be updated
RHOAI_RELEASE_VERSION="2.17"
KUEUE_VERSION="12"
CODEFLARE_SDK_VERSION="34"

# Simulated existing Confluence page content
EXISTING_CONTENT='<table class="relative-table wrapped" style="width: 21.25%;"><colgroup><col style="width: 26.1538%;" /><col style="width: 24.8718%;" /><col style="width: 48.9744%;" /></colgroup><tbody><tr><th>Release Version</th><th>Kueue</th><th>Codeflare-sdk</th></tr><tr><td>2.16</td><td>v0.8.0</td><td>0.24.0</td></tr><tr class=""><td>2.18</td><td>v0.10.0</td><td>0.25.1</td></tr><tr class=""><td>2.17</td><td>v0.8.1</td><td>0.25.0</td></tr></tbody></table>'

echo "Existing Content:"
echo "$EXISTING_CONTENT"

# Convert newlines to a placeholder to handle multi-line processing
PLACEHOLDER="__NL__"
MODIFIED_CONTENT=$(echo "$EXISTING_CONTENT" | tr '\n' "$PLACEHOLDER")

# Check if the release version already exists in the table
if echo "$MODIFIED_CONTENT" | grep -q "<tr[^>]*><td>$RHOAI_RELEASE_VERSION</td>"; then
    # Update the existing row with new values
    UPDATED_CONTENT=$(echo "$MODIFIED_CONTENT" | sed -E "s|(<tr[^>]*><td>$RHOAI_RELEASE_VERSION</td><td>)[^<]+(</td><td>)[^<]+(</td></tr>)|\1$KUEUE_VERSION\2$CODEFLARE_SDK_VERSION\3|")
else
    # If the row doesn't exist, insert a new row before the closing </tbody>
    UPDATED_ROW="<tr class=\"\"><td>$RHOAI_RELEASE_VERSION</td><td>$KUEUE_VERSION</td><td>$CODEFLARE_SDK_VERSION</td></tr>"
    UPDATED_CONTENT=$(echo "$MODIFIED_CONTENT" | sed "s|</tbody>|$UPDATED_ROW</tbody>|")
fi

# Restore newlines
UPDATED_CONTENT=$(echo "$UPDATED_CONTENT" | tr "$PLACEHOLDER" '\n')

# Print the updated content
echo "Updated Content:"
echo "$UPDATED_CONTENT"
