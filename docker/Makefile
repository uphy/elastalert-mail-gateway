clear:
	docker-compose exec elasticsearch curl -XDELETE -H 'Content-Type: application/json' localhost:9200/test	

	docker-compose stop elastalert
	docker-compose exec elasticsearch curl -XDELETE -H 'Content-Type: application/json' localhost:9200/elastalert_status
	docker-compose exec elasticsearch curl -XDELETE -H 'Content-Type: application/json' localhost:9200/elastalert_status_status
	docker-compose exec elasticsearch curl -XDELETE -H 'Content-Type: application/json' localhost:9200/elastalert_status_error
	docker-compose exec elasticsearch curl -XDELETE -H 'Content-Type: application/json' localhost:9200/elastalert_status_past
	docker-compose exec elasticsearch curl -XDELETE -H 'Content-Type: application/json' localhost:9200/elastalert_status_silence
	docker-compose start elastalert

insert:
	./insert-data.sh

run:
	docker-compose stop gateway elastalert
	docker-compose up --build -d
	./insert-data.sh
	docker-compose logs -f gateway