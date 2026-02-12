.PHONY: deploy-prod deploy-server-send deploy-server-swap deploy-scheduler deploy-worker deploy-tx-indexer deploy-grafana

NS ?= plugin-dca

deploy-prod: deploy-configs deploy-server-send deploy-server-swap deploy-scheduler deploy-worker deploy-tx-indexer deploy-grafana

deploy-configs:
	kubectl -n $(NS) apply -f deploy/prod

deploy-server-send:
	kubectl -n $(NS) apply -k deploy/overlays/send
	kubectl -n $(NS) rollout status deployment/server-send --timeout=300s

deploy-server-swap:
	kubectl -n $(NS) apply -k deploy/overlays/swap
	kubectl -n $(NS) rollout status deployment/server-swap --timeout=300s

deploy-scheduler:
	kubectl -n $(NS) apply -f deploy/01_scheduler.yaml
	kubectl -n $(NS) rollout status deployment/scheduler --timeout=300s

deploy-worker:
	kubectl -n $(NS) apply -f deploy/01_worker.yaml
	kubectl -n $(NS) rollout status deployment/worker --timeout=300s

deploy-tx-indexer:
	kubectl -n $(NS) apply -f deploy/01_tx_indexer.yaml
	kubectl -n $(NS) rollout status deployment/tx-indexer --timeout=300s

deploy-grafana:
	kubectl -n $(NS) apply -f deploy/02_grafana_dashboard.yaml
