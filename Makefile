BEATNAME=cmkbeat
BEATDIR=github.com/jeremyweader/cmkbeat
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile
.PHONY: deps
deps:
	glide up
	
.PHONY: config
config:
	echo "Update config file"
	-rm -f ${BEATNAME}.yml
	cat etc/beat.yml etc/config.yml | sed -e "s/beatname/${BEATNAME}/g" > ${BEATNAME}.yml
	-rm -f ${BEATNAME}.full.yml
	cat etc/beat.yml etc/config.full.yml | sed -e "s/beatname/${BEATNAME}/g" > ${BEATNAME}.full.yml

	# Update doc
	python ${ES_BEATS}/libbeat/scripts/generate_fields_docs.py $(PWD) ${BEATNAME} ${ES_BEATS}

	# Generate index templates
	python ${ES_BEATS}/libbeat/scripts/generate_template.py $(PWD) ${BEATNAME} ${ES_BEATS}
	python ${ES_BEATS}/libbeat/scripts/generate_template.py --es2x $(PWD) ${BEATNAME} ${ES_BEATS}

	# Update docs version
	cp ${ES_BEATS}/libbeat/docs/version.asciidoc docs/version.asciidoc

	# Generate index-pattern
	echo "Generate index pattern"
	-rm -f $(PWD)/etc/kibana/index-pattern/${BEATNAME}.json
	mkdir -p $(PWD)/etc/kibana/index-pattern
	python ${ES_BEATS}/libbeat/scripts/generate_index_pattern.py --index ${BEATNAME}-* --libbeat ${ES_BEATS}/libbeat --beat $(PWD)


.PHONY: install
install:
	mkdir -p /etc/$(BEATNAME)
	mkdir -p /usr/share/$(BEATNAME)/bin
	mkdir -p /var/lib/$(BEATNAME)
	mkdir -p /var/log/$(BEATNAME)
	cp $(BEATNAME) /usr/share/$(BEATNAME)/bin/.
	cp *.yml /etc/$(BEATNAME)/.
	cp *.json /etc/$(BEATNAME)/.
	cp system/$(BEATNAME).service /usr/lib/systemd/system/.
	systemctl enable $(BEATNAME).service

.PHONY: uninstall
uninstall:
	systemctl disable $(BEATNAME).service
	rm /usr/lib/systemd/system/$(BEATNAME).service
	rm -rf /var/log/$(BEATNAME)
	rm -rf /var/lib/$(BEATNAME)
	rm -rf /usr/share/$(BEATNAME)
	rm -rf /etc/$(BEATNAME)
	
.PHONY: all
all:
	make deps
	make cmkbeat
