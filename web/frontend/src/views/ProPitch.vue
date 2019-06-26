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
    >
      <template slot="basic" slot-scope="data">
        <img src="/images/heavy-check-mark.png" v-if="data.item.basic" class="checkmark">
      </template>
      <template slot="pro" slot-scope="data">
        <img src="/images/heavy-check-mark.png" v-if="data.item.pro" class="checkmark">
      </template>
    </b-table>
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
        { feature: "Unlimited posts", basic: true, pro: true },
        { feature: "Share updates privately", basic: false, pro: true },
        {
          feature: "Automatically send email summaries to your teammates",
          basic: false,
          pro: true
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

.checkmark {
  max-width: 20px;
}
</style>
