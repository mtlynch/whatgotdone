<template>
  <div class="shutdown-notice">
    <h1>Shutdown Notice</h1>

    <p>What Got Done is unfortunately shutting down at the end of 2025.</p>

    <p>
      I, Michael Lynch, founder of What Got Done will help anyone who wants to
      export their data and migrate it to another service.
    </p>

    <h2>Why?</h2>

    <h3>Low activity</h3>

    <p>There are not many users on What Got Done</p>

    <h2>What can users do?</h2>

    <p>
      The easiest option (and what Michael is doing) is to port your data to a
      static site generator, such as Hugo, Jekyll, or Gatsby. You are welcome to
      fork
      <a href="https://github.com/mtlynch/weeks.mtlynch.io"
        >Michael's Hugo implementation</a
      >
      (<a href="https://weeks.mtlynch.io/">live site</a>) and run the
      <a
        href="https://github.com/mtlynch/weeks.mtlynch.io/blob/0d2b39c75d0fba2015fe9e2afebd59aa2cff6acc/dev-scripts/import-from-whatgotdone"
        ><code>import-from-whatgotdone</code></a
      >
      script to import your data from What Got Done.
    </p>

    <p>
      Alternatively, if you'd like to continue using What Got Done, you can
      <a href="https://github.com/mtlynch/whatgotdone">self-host it</a>, though
      you'd have to write a bit of custom code. You'd have to replace
      authentication, as whatgotdone.com depends on UserKit, which is no longer
      accepting new customers, so you'd need to replace UserKit in the app.
    </p>

    <h2>Timeline</h2>
    <ul>
      <li>
        <b>May 2025</b>: Signups close. No new users may sign up to What Got
        Done.
      </li>
      <li>
        <b>July 5, 2025</b>: "Export as Markdown" feature
        <a href="https://github.com/mtlynch/whatgotdone/pull/963">added</a>. You
        can export all of your public posts as a ZIP of Markdown files.
      </li>
      <li>
        <b>August 2025</b>: Add a redirect feature that allows you to redirect
        your pages to a new URL.
      </li>
      <li><b>October 2025</b>: New updates will be closed.</li>
      <li>
        <b>January 2026</b>: Add a redirect feature that allows you to redirect
        your pages to a new URL.
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'ShutdownNoticePage',
};
</script>

<style scoped>
* {
  text-align: left;
}

h2 {
  margin-top: 2rem;
}
</style>
