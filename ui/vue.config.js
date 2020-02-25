const {gitDescribe, gitDescribeSync} = require('git-describe');
process.env.VUE_APP_GIT_HASH = gitDescribeSync().hash;

module.exports = {
  devServer: {
    port: 8081,
    disableHostCheck: true,
  },
  "transpileDependencies": [
    "vuetify"
  ]
};
