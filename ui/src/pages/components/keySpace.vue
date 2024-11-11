<template>
  <div class="card" style="margin-top: 1rem;">
    <div class="card-header">
      Key Space
    </div>
    <table class="table">
      <thead>
      <tr>
        <th scope="col">DbName</th>
        <th scope="col">Keys</th>
        <th scope="col">Expires</th>
        <th scope="col">Avg ttl</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(item,key) in keyspace" :key="key">
        <th scope="row">{{item.dbname}}</th>
        <td>{{item.keys}}</td>
        <td>{{item.expires}}</td>
        <td>{{item.avg_ttl}}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
<script setup>
import {defineProps} from "vue";

const props = defineProps({
  keyspace:Array
})

</script>