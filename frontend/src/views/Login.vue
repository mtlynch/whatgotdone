<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-9 col-md-7 col-lg-5 mx-auto">
        <div class="card card-signin my-5">
          <div class="card-body">
            <h5 class="card-title text-center">Sign In</h5>
            <form class="form-signin">
              <div class="form-label-group">
                <input
                  type="text"
                  id="username"
                  class="form-control"
                  placeholder="Username"
                  v-model="username"
                  required
                  autofocus
                />
                <label for="username">Username</label>
              </div>

              <div class="form-label-group">
                <input
                  type="password"
                  id="password"
                  class="form-control"
                  placeholder="Password"
                  v-model="password"
                  required
                />
                <label for="password">Password</label>
              </div>

              <button class="btn btn-lg btn-primary btn-block text-uppercase" type="submit">Sign in</button>
              <router-link class="d-block text-center mt-2 small" to="/register">Register</router-link>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import updateLoginState from "../controllers/LoginState.js";

export default {
  name: "Login",
  data() {
    return {
      username: null,
      password: null
    };
  }
};
</script>
