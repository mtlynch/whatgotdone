<template>
  <div class="container">
    <div class="row">
      <div class="col-lg-10 col-xl-9 mx-auto">
        <div class="card card-signin flex-row my-5">
          <div class="card-img-left d-none d-md-flex">
            <!-- Background image for card set in CSS! -->
          </div>
          <div class="card-body">
            <h5 class="card-title text-center">Register</h5>
            <form class="form-signin">
              <div class="form-label-group">
                <input
                  type="text"
                  id="inputUserame"
                  class="form-control"
                  placeholder="Username"
                  required
                  autofocus
                />
                <label for="inputUserame">Username</label>
              </div>

              <div class="form-label-group">
                <input
                  type="email"
                  id="inputEmail"
                  class="form-control"
                  placeholder="Email address"
                  required
                />
                <label for="inputEmail">Email address</label>
              </div>

              <hr />

              <div class="form-label-group">
                <input
                  type="password"
                  id="inputPassword"
                  class="form-control"
                  placeholder="Password"
                  required
                />
                <label for="inputPassword">Password</label>
              </div>

              <div class="form-label-group">
                <input
                  type="password"
                  id="inputConfirmPassword"
                  class="form-control"
                  placeholder="Password"
                  required
                />
                <label for="inputConfirmPassword">Confirm password</label>
              </div>

              <button class="btn btn-lg btn-primary btn-block text-uppercase" type="submit">Register</button>
              <router-link class="d-block text-center mt-2 small" to="/login">Sign In</router-link>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "SignUp",
  data() {
    return {
      username: null,
      password: null
    };
  }
};
</script>