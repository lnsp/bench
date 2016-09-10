github-release delete \
	--user lnsp \
	--repo bench \
	--tag nightly

github-release release \
	--user lnsp \
	--repo bench \
	--tag nightly \
	--name "Nightly" \
	--description "Automated build. **Careful, this release is not suited for production!**" \
	-p

for arch in *.tar.gz; do
	github-release upload \
		--user lnsp \
		--repo bench \
		--tag nightly \
		--name $arch \
		--file $arch;
done;

