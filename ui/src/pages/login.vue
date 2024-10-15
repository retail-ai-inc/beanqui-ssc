<template>
  <div class="container text-center" style="background: #f8f9fa">
    <div class="row align-items-start" style="height: 100vh;">
      <div class="col left-col">
        {{title}}
      </div>
      <div class="col right-col" >
        <div class="bq-box">
          <div style="width: 100%">
            <input class="form-control" type="text" placeholder="Username" name="userName" autocomplete="off" aria-label="default input example" v-model="user.username">
            <input class="form-control" type="password" placeholder="Password" name="password" autocomplete="off" aria-label="default input example" style="margin-top: 0.9375rem" v-model="user.password">
          </div>

          <button type="button" class="btn btn-primary" style="margin-top: 0.625rem" @click="onSubmit">Login</button>
          <div id="errorMsg" style="color: red;margin-top:0.625rem;">{{msg}}</div>
        </div>

`
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive,toRefs,onMounted,onUnmounted } from "vue";
import { useRouter } from 'vueRouter';

const data = reactive({
  user:{"username":"","password":""},
  msg:"",
  title:config.title,
})
const useRe = useRouter();

async function onSubmit(event){

  if (data.user.username == "" || data.user.password == ""){
    console.log("can not empty");
    return;
  }
  //,{headers:{"Content-Type":"multipart/form-data"}}
  try{
    let res = await loginApi.Login(data.user.username,data.user.password);
    sessionStorage.setItem("token",res.data.token);
    useRe.push("/admin/home");
  }catch(err){
    if (err.response.status === 401){
      data.msg = err.response.data.msg;
    }
  }
}
const {user,msg,title} = toRefs(data);
</script>
<style scoped>
.left-col{
  background: #7364dd;height: 100%;display: flex;justify-content: center;align-items: center;font-size: 1.5rem;font-weight: bold;color: #fff;
}
.right-col{
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100vh;
}
.bq-box{
  display: flex;width: 70%;
  background: #fff;
  padding: 1.56rem;
  border:0.0625rem solid #ced4da;
  border-radius: 0.3125rem;
  box-shadow:0.1rem 0.2rem 0.2rem #ccc;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  margin-left: 1.875rem;
}
</style>
