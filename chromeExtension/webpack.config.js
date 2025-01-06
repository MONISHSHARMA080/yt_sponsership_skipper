import path from "path";

module.exports = {
  entry: {
    index: "./src/service_worker.js", // Update this to the correct entry file.
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "[name].js",
  },
  resolve: {
    extensions: [".js", ".jsx", ".ts", ".tsx"], // Support both JavaScript and TypeScript files.
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx|ts|tsx)$/,
        exclude: /node_modules/,
        use: "babel-loader", // Use Babel to handle modern JS/TS.
      },
      {
        test: /\.html$/,
        use: [
          {
            loader: "html-loader",
          },
        ],
      },
    ],
  },
  plugins: [],
};
