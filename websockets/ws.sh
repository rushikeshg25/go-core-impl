#!/bin/bash

WS_URL="http://localhost:8080/ws"
SEC_WEBSOCKET_KEY="dGhlIHNhbXBsZSBub25jZQ=="  # Example Base64-encoded key

curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Key: $SEC_WEBSOCKET_KEY" \
  -H "Sec-WebSocket-Version: 13" \
  "$WS_URL"
