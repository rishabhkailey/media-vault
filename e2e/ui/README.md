

Linux

```bash
# file name and source separated by '|'

# sources
# https://file-examples.com/wp-content/storage/2020/03/file_example_WEBM_640_1_4MB.webm
# https://file-examples.com/wp-content/storage/2017/10/file_example_JPG_1MB.jpg
# https://file-examples.com/wp-content/storage/2017/10/file_example_PNG_1MB.png
# https://file-examples.com/wp-content/storage/2020/03/file_example_WEBP_500kB.webp

files=(
  "video-01.webm  | https://file-examples.com/storage/fe8119f4e865f33329898be/2020/03/file_example_WEBM_640_1_4MB.webm"
  "image-01.jpg   | https://file-examples.com/storage/fe8119f4e865f33329898be/2017/10/file_example_JPG_1MB.jpg"
  "image-02.png   | https://file-examples.com/storage/fe8119f4e865f33329898be/2017/10/file_example_PNG_1MB.png"
  "image-03.webp  | https://file-examples.com/storage/fe8119f4e865f33329898be/2020/03/file_example_WEBP_500kB.webp"
)
mkdir -p test-files
for _file in "${files[@]}"; do
  # remove tabs/space
  file=$(echo $_file | tr -d ' ')
  fileName=$(echo $file | cut -d'|' -f1)
  fileSource=$(echo $file | cut -d'|' -f2)
  wget "$fileSource" -O "test-files/$fileName"
done
```

update `cypress/fixtures/test-files.json` with test file names

```bash
yarn install
node_modules/cypress/bin/cypress open
```
