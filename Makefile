WORKDIR = $(shell pwd)

# 爬虫抓取网站静态页面
.PHONY: crawler
crawler:
	cd $(WORKDIR)/crawler \
	&& go mod tidy  \
	&& go run . -dist ../dist

# 编译 docker 镜像
.PHONY: image
image:
	docker build -t coolshell-nginx .
