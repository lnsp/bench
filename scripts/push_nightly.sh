github-release delete \
	--user lnsp \
	--repo bench \
	--tag nightly

cd $GOPATH/src/github.com/lnsp/bench/
git config --global user.email "builds@travis-ci.com"
git config --global user.name "Travis CI"
git tag --delete nightly
git push --delete origin nightly
git tag nightly -a -m "Automated nightly builds"
git push -q https://lnsp:$GITHUB_TOKEN@github.com/lnsp/bench :nightly

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

