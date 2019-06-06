<template>
  <div class="pro">
    <div class="jumbotron">
      <h1 class="display-3">Upgrade to Pro</h1>
      <p class="lead">Unlock additional features to collaborate with your teammates.</p>
      <hr class="my-4">
      <p>Pricing starts at just $9.99/month per user.</p>
      <p class="lead float-right">
        <b-button class="btn btn-primary btn-lg" variant="outline-primary" to="/upgrade">Upgrade Now</b-button>
      </p>
      <p>&nbsp;</p>
    </div>
    <b-table
      class="features"
      head-variant="light"
      striped
      hover
      outlined
      :items="items"
      :fields="fields"
    ></b-table>
  </div>
</template>

<script>
export default {
  name: "ProUpgrade",
  data() {
    return {
      fields: [
        {
          key: "feature",
          label: ""
        },
        {
          key: "basic"
        },
        {
          key: "pro"
        }
      ],
      items: [
        { feature: "Unlimited posts", basic: "✔️", pro: "✔️" },
        { feature: "Share updates privately", basic: "️", pro: "✔️" },
        {
          feature: "Automatically send email summaries to your teammates",
          basic: "️",
          pro: "✔️"
        }
      ]
    };
  }
};
</script>

<style scoped>
.pro {
  text-align: left;
}

.features {
  margin-top: 100px;
}
</style>
