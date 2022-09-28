build-bot:
	mkdir -p bot/src
	cp -r ../source/bot/* bot/src

	docker build -t shortcut-telegram-integration-bot:latest bot
	rm -rf bot/src

build-webhook:
	mkdir -p webhook/src
	cp -r ../source/webhook/* webhook/src

	docker build -t shortcut-telegram-integration-webhook:latest webhook
	rm -rf webhook/src

build-db:
	mkdir -p db/src
	cp -r ../source/db/* db/src

	docker build -t shortcut-telegram-integration-db:latest db
	rm -rf db/src