.PHONY: build deploy local datainit generate-init-data

build:
	sam build

deploy: build
	$(eval SAMENV := $(shell cat .env | tr '\n' ' '))
	sam deploy --parameter-overrides $(SAMENV)

local: build
	sam local start-api --env-vars env.json

# DynamoDB初期データのロード
datainit: generate-init-data
	aws --region ap-northeast-1 --profile serverlessUserAWS dynamodb batch-write-item --request-items file://initial_data.json

# DynamoDBの初期データ定義ファイルの生成
# 予め定義したデータにここの処理でテーブル名を上書きしている
generate-init-data:
	$(eval SAMPLE := $(shell aws dynamodb list-tables --region ap-northeast-1 --profile serverlessUserAWS | jq -r '.TableNames[] | select(startswith("t2-bot-stack-DynamoDBTable"))'))
	sed -i .bak "s/t2-bot-stack-DynamoDBTable-.\{12\}/$(SAMPLE)/" initial_data.json
