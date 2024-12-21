const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
  entry: './src/client.js', // Your main JS file
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js',
  },
  mode: 'development', // or 'production' for production builds
  plugins: [
    new HtmlWebpackPlugin({
      template: './src/index.html', // Specify the source HTML file
      filename: 'index.html', // Output file in the `dist` folder
      favicon: './src/assets/favicon.ico', // Add the favicon
      inject: false
    }),
    new CopyWebpackPlugin({
      patterns: [
        { from: './src/assets/favicon.ico', to: 'favicon.ico' }, // Copy favicon
      ],
    }),
  ],
  devServer: {
    static: './dist',
  },
  module: {
    rules:[
      { test: /\.css$/, use: [ 'style-loader', 'css-loader', 'postcss-loader' ] }, 
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: 'babel-loader',
      },
  ]
  },
};

