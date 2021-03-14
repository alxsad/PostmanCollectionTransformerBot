# PostmanCollectionTransformerBot
Telegram bot to transform postman collections from v1 to v2 written in Go

Build:
```commandline
GOOS=linux go build .
```

Run:
```commandline
./PostmanCollectionTransformerBot -token={telegram_token}
```

/lib/systemd/system/PostmanCollectionTransformerBot.service:
```ini
[Unit]
Description=PostmanCollectionTransformerBot
After=network-online.target

[Service]
Type=simple
EnvironmentFile=-/root/PostmanCollectionTransformerBot.env
ExecStart=/root/PostmanCollectionTransformerBot -token $POSTMAN_COLLECTION_TRANSFORMER_BOT_TOKEN
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
