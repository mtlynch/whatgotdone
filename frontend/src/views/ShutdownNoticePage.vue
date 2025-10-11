<template>
  <div class="shutdown-notice">
    <h1>Shutdown Notice</h1>

    <p>What Got Done is unfortunately shutting down at the end of 2025.</p>

    <p>
      I, Michael Lynch, founder of What Got Done am happy to help anyone who
      wants to export their data and migrate it to another service.
    </p>

    <h2>Why?</h2>

    <h3>Low activity</h3>

    <p>
      There are not many users on What Got Done, so the site is a bit overkill
      in terms of time and cost for the number of users it serves.
    </p>

    <h3>Changing views of open-source</h3>

    <p>
      I originally created What Got Done with hopes of it being a commercial
      product in 2019. When
      <a href="https://mtlynch.io/retrospectives/2019/08/">it flopped</a>, I
      decided to run What Got Done as an open-source project and a free service.
    </p>

    <p>
      Since 2019, my views on open-source have changed. I no longer think it's
      sustainable for me to run a service for free. All of my open-source
      projects these days require users to self-host, and that's been a more
      sustainable strategy for my personal bandwidth.
    </p>

    <h2>What can users do?</h2>

    <p>
      The easiest option (and what Michael is doing) is to port your data to a
      static site generator, such as Hugo, Jekyll, or Gatsby. You are welcome to
      fork
      <a href="https://github.com/mtlynch/weeks.mtlynch.io-old"
        >Michael's Hugo implementation</a
      >
      (<a href="https://weeks.mtlynch.io/">live site</a>) and run the
      <a
        href="https://github.com/mtlynch/weeks.mtlynch.io-old/blob/master/dev-scripts/import-from-whatgotdone"
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
      You'd also have to write custom code to import from your JSON or Markdown
      export file back into What Got Done's SQLite database.
    </p>

    <h2>Shutdown timeline</h2>
    <ul>
      <li>
        <b>May 2025</b>: Signups close. No new users may sign up to What Got
        Done.
      </li>
      <li>
        <b>July 2025</b>:
        <router-link to="/export">"Export as Markdown"</router-link> feature
        added. You can export all of your public posts as a ZIP of Markdown
        files.
      </li>
      <li>
        <b>July 2025</b>: Add support for a
        <router-link to="/forwarding-address">forwarding address</router-link>
        feature that allows you to redirect your What Got Done posts to a new
        URL.
      </li>
      <li>
        <b>October 2025</b>: Users will no longer be able to post new updates.
      </li>
      <li>
        <b>January 2026</b>: Replace the site with a shutdown notice. Keep
        redirects in place on a best-effort basis.
      </li>
    </ul>

    <h2>Questions</h2>

    <p>
      If you have any questions about the shutdown, please email
      michael@whatgotdone.com.
    </p>
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
  max-width: 70ch;
}

h2 {
  margin-top: 3rem;
}

h3 {
  margin-top: 2rem;
}
</style>
