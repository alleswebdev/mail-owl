APPS = scheduler builder email-sender sms-sender
RUN = $(APPS:=-run)
JSON_MODELS = internal/models/scheduler_notice.go

.PHONY: all clean $(APPS) stop build run migrate json

all: build

build: $(APPS)
run: build $(RUN)

$(APPS):
			go build -o $@ cmd/$@/main.go
$(RUN):
			 @./$(@:-run=) > $@.log 2>&1 & \
			  echo pid of $(@:-run=):$$!
clean:
			@-rm $(APPS) *.log 2>/dev/null
stop:
			-killall -s 9 $(APPS)

migrate: migrate-init
			go run cmd/migrations/migrate.go
migrate-init:
			go run cmd/migrations/migrate.go init
json:
			easyjson -all $(JSON_MODELS)