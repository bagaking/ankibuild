# Makefile.bundle
# 获取当前目录下所有 file_bundle_rc 文件列表
FILE_BUNDLE_RCS := $(shell find . -name '*.file_bundle_rc')

# 定义后缀替换规则：将 .file_bundle_rc 结尾的文件替换为 .bundle.txt
BUNDLES := $(FILE_BUNDLE_RCS:.file_bundle_rc=.bundle.txt)

# 默认目标
all: $(BUNDLES)

# 规则来生成 bundle 文件
%.bundle.txt: %.file_bundle_rc
	file_bundle -v -i $< -o $@

# 清理
clean:
	rm -f $(BUNDLES)