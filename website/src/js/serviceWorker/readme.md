## Why webpack required
firefox doesn't support es6 imports in service worker

> TODO: this stopped working after upgrading npm packages and changing the type = module in package.json
> for debugging we can directly use es6 file in chrome and attach a debugger


## commands to run
```bash
cd src/js/serviceWorker
npx webpack --mode development --config webpack.config.cjs
```

## todo integrate with vite

