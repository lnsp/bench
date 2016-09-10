github-release delete \
	--user lnsp \
	--repo bench \
	--tag nightly

github-release release \
	--user lnsp \
	--repo bench \
	--tag nightly \
	--name "Nightly" \
	--description "Builds only (source code is not up to date)"
	--pre-release

for arch in *.tar.gz; do
	github-release upload \
		--user lnsp \
		--repo bench \
		--tag nightly \
		--name $arch \
		--file $arch;
done;

