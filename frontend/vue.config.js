module.exports = {
  runtimeCompiler: true,
  devServer: {
    disableHostCheck: true,
  },
  // We have to disable removeAttributeQuotes otherwise webpack will
  // prematurely remove quotes from things that actually need them once go
  // populates template variables like:
  //
  //    content="[[.Title]]"
  //
  chainWebpack: config => {
    config
      .plugin('html').init((Plugin, args) => {
        const newArgs = {
            ...args[0],
        };
        if (newArgs.minify) {
          newArgs.minify.removeAttributeQuotes = false;
        }
        return new Plugin(newArgs);
    });
  }
}