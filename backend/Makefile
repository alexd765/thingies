HANDLER ?= handler
PACKAGE ?= package

all:
	@docker run --rm \
		-v $(GOPATH):/go -v $(PWD):/tmp \
		-e "HANDLER=$(HANDLER)" -e "PACKAGE=$(PACKAGE)" \
		eawsy/aws-lambda-go-shim make _all

clean: _clean

_all: _clean
	@echo -ne "build..."\\r
	@go build -buildmode=plugin -ldflags='-w -s' -o $(HANDLER).so
	@chown $(shell stat -c '%u:%g' .) $(HANDLER).so
	@echo -ne "build, pack"\\r
	@zip -q $(PACKAGE).zip $(HANDLER).so
	@echo -ne "build, pack, inject"\\r
	@cd /; mv /shim $(HANDLER); zip -q -r /tmp/$(PACKAGE).zip $(HANDLER)
	@chown $(shell stat -c '%u:%g' .) $(PACKAGE).zip
	@echo -ne "build, pack, inject, go!"\\n

_clean:
	@rm -rf $(PACKAGE).zip $(HANDLER).so