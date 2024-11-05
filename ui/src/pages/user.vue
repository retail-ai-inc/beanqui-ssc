<template>
  <div class="user">
    <div class="d-flex justify-content-end">
      <button type="button" class="btn btn-primary">Add</button>
    </div>
    <table class="table table-striped table-hover">
      <thead>
      <tr>
        <th scope="col">#</th>
        <th scope="col">Account</th>
        <th scope="col">Active</th>
        <th scope="col">Type</th>
        <th scope="col">Detail</th>
        <th scope="col">Action</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(item, key) in users" :key="key" style="height: 3rem;line-height:3rem">
        <th scope="row">{{parseInt(key)+1}}</th>
        <td>{{item.account}}</td>
        <td>{{item.active}}</td>
        <td>{{item.type}}</td>
        <td>{{item.detail}}</td>
        <td>
          <a href="javascript:;">Edit</a> | <a href="javascript:;">Delete</a>
        </td>
      </tr>
      </tbody>

    </table>
  </div>
</template>
<script setup>
import { reactive,onMounted,toRefs,onUnmounted } from "vue";

let data = reactive({
  users:[]
})

onMounted(async ()=>{
  let res = await userApi.List();
  data.users = res.data;
  console.log(res);
})

function detailDlq(item){

}

function retryItem(item){

}

function deleteItem(item){

}

const {users} = toRefs(data);
</script>

<style scoped>
.user{
  transition: opacity 0.5s ease;
  opacity: 1;
}
</style>