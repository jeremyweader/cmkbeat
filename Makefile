BEATNAME=cmkbeat
BEATDIR=github.com/jeremyweader/cmkbeat
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile


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
