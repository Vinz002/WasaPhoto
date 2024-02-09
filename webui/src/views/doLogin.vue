// Path: webui/src/views/doLogin.vue
<script setup>
</script>
<script>
export default {
    data:
     function(){
            return {
                username: null,
                token: null,
            }
    },
    methods: {
        async dologin() {
            this.errormsg = null;
            try {
                let response = await this.$axios.post("/session", {
                    username: this.username,
                });

                if (response.status != 201) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                console.log(response)
                this.token = response.data.Id;
                localStorage.setItem("Identifier",this.token)
                localStorage.setItem("ProfileId",this.token)
                console.log("Login Success: ", this.token)
                this.$router.push({ path: "/user/stream/" + this.token });

            } catch (e) {
                this.errormsg = e.toString();
                console.log(e)
            }
        }
    }
}

</script>

<template>
  <div class="login-container">
    <h1 class="login-title">Benvenuto!</h1>
    <form @submit.prevent="dologin" class="login-form">
      <div class="form-group">
        <label for="username" class="visually-hidden">Username</label>
        <div class="input-group">
          <input
            type="text"
            id="username"
            class="form-control"
            placeholder="Inserisci il tuo username"
            v-model="username"
            required
          />
          <button type="submit" class="btn btn-primary">Login</button>
        </div>
      </div>
    </form>
    <div v-if="errormsg" class="alert alert-danger mt-3" role="alert">
      {{ errormsg }}
    </div>
  </div>
</template>

<style>
.login-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #f8f9fa;
}

.login-title {
  text-align: center;
  color: #007bff;
  font-size: 24px;
  margin-bottom: 20px;
}

.login-form {
  margin-top: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.input-group {
  display: flex;
}

.form-control {
  flex: 1;
  padding: 12px;
  font-size: 16px;
  border: 1px solid #ced4da;
  border-radius: 4px 0 0 4px;
}

.btn-primary {
  background-color: #007bff;
  border: 1px solid #007bff;
  border-radius: 0 4px 4px 0;
  color: #fff;
  cursor: pointer;
}

.btn-primary:hover {
  background-color: #0056b3;
  border: 1px solid #0056b3;
}

.alert {
  margin-top: 15px;
}
</style>