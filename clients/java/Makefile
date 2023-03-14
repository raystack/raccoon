.PHONY: all

all: generate-proto build

# generates the raccoon protos from the https://github.com/goto/proton using the buf
generate-proto:
	rm -rf .temp
	mkdir -p .temp
	curl -o .temp/proton.tar.gz -L http://api.github.com/repos/goto/proton/tarball/main; tar xvf .temp/proton.tar.gz -C .temp/ --strip-components 1
	buf generate --path=.temp/goto/raccoon

clean:
	rm -rf .temp

build:
	 ./gradlew build