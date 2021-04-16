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
  chainWebpack: (config) => {
    config.plugin('html').tap((args) => {
      if (args[0].minify) {
        args[0].minify.removeAttributeQuotes = false;
      }
      return args;
    });
  },
};
