// eslint-disable-next-line no-undef
module.exports = {
  entry: "./src/serviceWorker.ts",
  output: {
    filename: "serviceWorker.js",
  },
  resolve: {
    extensions: [".ts", ".js"],
  },
  module: {
    rules: [
      {
        test: /\.ts$/,
        exclude: /node_modules/,
        use: "ts-loader",
      },
    ],
  },
};
